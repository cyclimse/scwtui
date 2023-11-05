package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cyclimse/scaleway-dangling/internal/discovery"
	"github.com/cyclimse/scaleway-dangling/internal/discovery/demo"
	"github.com/cyclimse/scaleway-dangling/internal/discovery/scaleway"
	"github.com/cyclimse/scaleway-dangling/internal/observability/cockpit"
	"github.com/cyclimse/scaleway-dangling/internal/resource"
	"github.com/cyclimse/scaleway-dangling/internal/search/bleve"
	"github.com/cyclimse/scaleway-dangling/internal/store/sqlite"
	"github.com/cyclimse/scaleway-dangling/internal/ui"
	"github.com/cyclimse/scaleway-dangling/internal/ui/keys"
	"github.com/cyclimse/scaleway-dangling/internal/ui/scenes"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/sync/errgroup"
)

type TuiCmd struct {
	Demo bool `default:"false" help:"Run the TUI in demo mode"`
}

func (cmd *TuiCmd) Run(cmdCtx *CmdContext) error {
	logger, err := initLogger()
	if err != nil {
		return err
	}

	p, err := loadScalewayProfile(cmdCtx.Profile)
	if err != nil {
		return err
	}

	store, err := sqlite.NewStore(context.Background(), "scaleway-dangling.db")
	if err != nil {
		return err
	}
	defer store.Close()

	var discoverer discovery.ResourceDiscoverer
	var client *scw.Client
	var projects []resource.Resource

	if cmd.Demo {
		demoDiscovery := demo.NewDiscovery()
		copy(projects, demoDiscovery.Projects())
		discoverer = demoDiscovery
	} else {
		client, err = scw.NewClient(scw.WithUserAgent("scaleway-dangling"), scw.WithProfile(p))
		if err != nil {
			return err
		}

		projects, err = scaleway.ListProjects(context.Background(), client)
		if err != nil {
			return err
		}

		discoverer = scaleway.NewResourceDiscoverer(client, projects, &scaleway.ResourceDiscovererConfig{
			NumWorkers: 10,
		})
	}

	projectIDsToNames := make(map[string]string, len(projects))
	for _, p := range projects {
		metadata := p.Metadata()
		projectIDsToNames[metadata.ID] = metadata.Name
	}

	search, err := bleve.NewSearch(projectIDsToNames)
	if err != nil {
		return err
	}

	for _, r := range projects {
		if err := store.Store(context.Background(), r); err != nil {
			logger.Error("failed to store resource", slog.Any("resource", r))
			return err
		}
		if err := search.Index(r); err != nil {
			logger.Error("failed to index resource", slog.Any("resource", r))
			return err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, runCtx := errgroup.WithContext(ctx)

	ch := make(chan resource.Resource, 10000)

	g.Go(func() error {
		err := discoverer.Discover(runCtx, ch)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			slog.With("err", err).Error("failed to discover resources")
		}
		return err
	})
	g.Go(func() error {
		for {
			select {
			case <-runCtx.Done():
				return runCtx.Err()
			case r, ok := <-ch:
				if !ok {
					return nil
				}
				if err := store.Store(runCtx, r); err != nil {
					logger.Error("failed to store resource", slog.Any("resource", r))
					return err
				}
				if err := search.Index(r); err != nil {
					logger.Error("failed to index resource", slog.Any("resource", r))
					return err
				}
			}
		}
	})

	profileName := cmdCtx.Profile
	if profileName == "" {
		profileName = "default"
	}

	appState := ui.ApplicationState{
		Logger:            logger,
		Store:             store,
		Search:            search,
		ScwClient:         client,
		ScwProfileName:    profileName,
		Monitor:           cockpit.NewCockpit(logger, client),
		Keys:              keys.DefaultKeyMap(),
		ProjectIDsToNames: projectIDsToNames,
	}
	m := scenes.Root(appState)

	g.Go(func() error {
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			return err
		}
		cancel()
		return nil
	})

	err = g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
	}

	return nil
}

func initLogger() (*slog.Logger, error) {
	w, err := os.OpenFile("scaleway-dangling.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	loggerHandler := slog.NewJSONHandler(w, nil)
	return slog.New(loggerHandler), nil
}

func loadScalewayProfile(profileName string) (*scw.Profile, error) {
	cfg, err := scw.LoadConfig()
	if err != nil {
		return nil, err
	}

	p, err := cfg.GetActiveProfile()
	if err != nil {
		return nil, err
	}

	if profileName != "" {
		p, err = cfg.GetProfile(profileName)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}
