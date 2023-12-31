package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cyclimse/scwtui/internal/discovery"
	demo_discovery "github.com/cyclimse/scwtui/internal/discovery/demo"
	"github.com/cyclimse/scwtui/internal/discovery/scaleway"
	"github.com/cyclimse/scwtui/internal/observability/cockpit"
	demo_monitor "github.com/cyclimse/scwtui/internal/observability/demo"
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
func (cmd *TuiCmd) Run(rs *RootState) error {
	logger := rs.Logger

	p, err := loadScalewayProfile(rs.Profile)
	if err != nil {
		return err
	}

	store, err := sqlite.NewStore(context.Background(), "scwtui.db")
	if err != nil {
		return err
	}
	defer store.Close()

	var discoverer discovery.ResourceDiscoverer
	var monitor resource.Monitorer

	var client *scw.Client
	var projects []resource.Resource

	if cmd.Demo {
		store.SetUnmarshaller(&sqlite.DemoResourceUnmarshal{})

		projects = demo_discovery.ListProjects()
		discoverer = demo_discovery.NewDiscovery(projects)
		monitor = demo_monitor.NewDemo()
	} else {
		client, err = scw.NewClient(
			scw.WithUserAgent(fmt.Sprintf("scwtui/%s", rs.Version)),
			scw.WithProfile(p),
		)
		if err != nil {
			return err
		}

		projects, err = scaleway.ListProjects(context.Background(), client)
		if err != nil {
			return err
		}

		discoverer = scaleway.NewResourceDiscoverer(logger, client, projects, &scaleway.ResourceDiscovererConfig{
			NumWorkers: 10,
			MaxRetries: 3,
		})

		monitor = cockpit.NewCockpit(logger, client)
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
			logger.Error("tui: failed to store resource", slog.Any("resource", r), slog.String("err", err.Error()))
			return err
		}
		if err := search.Index(r); err != nil {
			logger.Error("tui: failed to index resource", slog.Any("resource", r), slog.String("err", err.Error()))
			return err
		}
	}

	profileName := rs.Profile
	if profileName == "" {
		profileName = "default"
	}

	appState := ui.ApplicationState{
		Logger: logger,

		Store:   store,
		Search:  search,
		Monitor: monitor,

		ScwClient:         client,
		ScwProfileName:    profileName,
		ProjectIDsToNames: projectIDsToNames,

		Keys: ui.DefaultKeyMap(),

		Styles:                 ui.DefaultStyles(),
		SyntaxHighlighterTheme: rs.Config.Tui.Theme,
	}
	m := scenes.Root(appState)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := start(ctx, logger, discoverer, resource.NewIndex(store, search))

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

func start(ctx context.Context, logger *slog.Logger, discoverer discovery.ResourceDiscoverer, index resource.Indexer) *errgroup.Group {
	g, runCtx := errgroup.WithContext(ctx)

	ch := make(chan resource.Resource, channelBufferSize)

	g.Go(func() error {
		err := discoverer.Discover(runCtx, ch)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			logger.Error("tui: failed to discover resources", slog.String("err", err.Error()))
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
				if err := index.Index(runCtx, r); err != nil {
					logger.Error("tui: failed to index resource", slog.Any("resource", r), slog.String("err", err.Error()))
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
