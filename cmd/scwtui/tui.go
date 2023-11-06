package main

import (
	"context"
	"errors"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cyclimse/scwtui/internal/discovery"
	"github.com/cyclimse/scwtui/internal/discovery/demo"
	"github.com/cyclimse/scwtui/internal/discovery/scaleway"
	"github.com/cyclimse/scwtui/internal/observability/cockpit"
	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/search/bleve"
	"github.com/cyclimse/scwtui/internal/store/sqlite"
	"github.com/cyclimse/scwtui/internal/ui"
	"github.com/cyclimse/scwtui/internal/ui/scenes"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/sync/errgroup"
)

type TuiCmd struct {
	Demo bool `default:"false" help:"Run the TUI in demo mode"`
}

//nolint:funlen,gocognit // this is a command
func (cmd *TuiCmd) Run(cmdCtx *CmdContext) error {
	logger := cmdCtx.Logger

	p, err := loadScalewayProfile(cmdCtx.Profile)
	if err != nil {
		return err
	}

	store, err := sqlite.NewStore(context.Background(), "scwtui.db")
	if err != nil {
		return err
	}
	defer store.Close()

	var discoverer discovery.ResourceDiscoverer
	var client *scw.Client
	var projects []resource.Resource

	if cmd.Demo {
		demoDiscovery := demo.NewDiscovery()
		projects = demoDiscovery.Projects()

		logger.Info("running in demo mode", slog.Any("projects", projects))

		discoverer = demoDiscovery
		store.SetUnmarshaller(&sqlite.DemoResourceUnmarshal{})
	} else {
		client, err = scw.NewClient(scw.WithUserAgent("scwtui"), scw.WithProfile(p))
		if err != nil {
			return err
		}

		projects, err = scaleway.ListProjects(context.Background(), client)
		if err != nil {
			return err
		}

		discoverer = scaleway.NewResourceDiscoverer(logger, client, projects, &scaleway.ResourceDiscovererConfig{
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

	profileName := cmdCtx.Profile
	if profileName == "" {
		profileName = "default"
	}

	appState := ui.ApplicationState{
		Logger: logger,

		Store:   store,
		Search:  search,
		Monitor: cockpit.NewCockpit(logger, client),

		ScwClient:         client,
		ScwProfileName:    profileName,
		ProjectIDsToNames: projectIDsToNames,

		Keys: ui.DefaultKeyMap(),

		Styles:                 ui.DefaultStyles(),
		SyntaxHighlighterTheme: cmdCtx.Config.Tui.Theme,
	}
	m := scenes.Root(appState)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := start(ctx, logger, discoverer, store, search)

	g.Go(func() error {
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			return err
		}
		cancel() // cancel the context after the ui has exited
		return nil
	})

	if err = g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	return nil
}

const (
	channelBufferSize = 100
)

func start(ctx context.Context, logger *slog.Logger, discoverer discovery.ResourceDiscoverer, store resource.Storer, search resource.Searcher) *errgroup.Group {
	g, runCtx := errgroup.WithContext(ctx)

	ch := make(chan resource.Resource, channelBufferSize)

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

	return g
}

func loadScalewayProfile(profileName string) (*scw.Profile, error) {
	cfg, err := scw.LoadConfig()
	if err != nil {
		return nil, err
	}

	var p *scw.Profile

	if profileName != "" {
		// if the profile is overridden via the command line,
		// we load it directly
		p, err = cfg.GetProfile(profileName)
		if err != nil {
			return nil, err
		}
	} else {
		// otherwise we load the active profile
		p, err = cfg.GetActiveProfile()
		if err != nil {
			return nil, err
		}
	}

	// finally, merge it with the environment variables overrides
	p = scw.MergeProfiles(p, scw.LoadEnvProfile())

	return p, nil
}
