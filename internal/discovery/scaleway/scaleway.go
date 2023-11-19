package scaleway

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/sync/errgroup"
)

// ErrShouldRetry is returned when the request should be retried.
var ErrShouldRetry = errors.New("should retry")

func NewResourceDiscoverer(logger *slog.Logger, client *scw.Client, projects []resource.Resource, config *ResourceDiscovererConfig) *ResourceDiscover {
	return &ResourceDiscover{
		logger:    logger,
		client:    client,
		config:    config,
		projects:  projects,
		requested: make(chan requestResources, 100),
		regions:   discoveryRegions(logger, client),
		zones:     discoveryZones(logger, client),
	}
}

type ResourceDiscover struct {
	logger *slog.Logger

	client *scw.Client
	config *ResourceDiscovererConfig

	projects  []resource.Resource
	requested chan requestResources

	regions []scw.Region
	zones   []scw.Zone
}

type ResourceDiscovererConfig struct {
	NumWorkers int
	MaxRetries int
}

func (d *ResourceDiscover) Discover(ctx context.Context, ch chan resource.Resource) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(d.config.NumWorkers)

	// Discover resources
	for i := 0; i < d.config.NumWorkers; i++ {
		g.Go(func() error {
			return d.runWorker(ctx, ch)
		})
	}

	// Request resources
	d.requested <- requestResources{
		Get: d.discoverCockpits,
	}
	d.requested <- requestResources{
		Get: d.discoverIAMApplications,
	}
	for _, region := range d.regions {
		region := region // !important
		d.discoverInRegion(region, d.discoverRegistryNamespacesInRegion)
		d.discoverInRegion(region, d.discoverContainersInRegion)
		d.discoverInRegion(region, d.discoverFunctionsInRegion)
		d.discoverInRegion(region, d.discoverRdbInstancesInRegion)
		d.discoverInRegion(region, d.discoverKapsuleClustersInRegion)
		d.discoverInRegion(region, d.discoverJobsInRegion)
	}
	for _, zone := range d.zones {
		zone := zone // !important
		d.discoverInZone(zone, d.discoverInstancesInZone)
	}

	return g.Wait()
}

func (d *ResourceDiscover) discoverInRegion(region scw.Region, discoverFunc func(ctx context.Context, region scw.Region) ([]resource.Resource, error)) {
	d.requested <- requestResources{
		Get: func(ctx context.Context) ([]resource.Resource, error) {
			return discoverFunc(ctx, region)
		},
	}
}

func (d *ResourceDiscover) discoverInZone(zone scw.Zone, discoverFunc func(ctx context.Context, zone scw.Zone) ([]resource.Resource, error)) {
	d.requested <- requestResources{
		Get: func(ctx context.Context) ([]resource.Resource, error) {
			return discoverFunc(ctx, zone)
		},
	}
}

type requestResources struct {
	Get          func(ctx context.Context) ([]resource.Resource, error)
	CurrentRetry int
}

func (d *ResourceDiscover) runWorker(ctx context.Context, ch chan resource.Resource) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case req := <-d.requested:
			resources, err := req.Get(ctx)
			if err != nil {
				if errors.Is(err, ErrShouldRetry) {
					// Retry later
					req.CurrentRetry++
					if req.CurrentRetry < d.config.MaxRetries {
						d.requested <- req
						continue
					}
					return err
				}
				d.logger.With("err", err).Error("discover: failed to get resources")
			}

			d.logger.Info("discover: got resources", slog.Int("num_resources", len(resources)))
			for _, r := range resources {
				ch <- r
			}
		}
	}
}

func handleRequestError(err error) error {
	if err == nil {
		return nil
	}

	// nolint:errorlint // will not be wrapped
	if _, ok := err.(scw.SdkError); !ok {
		return err
	}

	var resErr *scw.ResourceNotFoundError
	if errors.As(err, &resErr) {
		return nil
	}

	var respErr *scw.ResponseError
	if errors.As(err, &respErr) {
		if respErr.StatusCode == http.StatusTooManyRequests {
			return ErrShouldRetry
		}
	}

	return err
}

// discoveryRegions returns the list of regions to discover resources in.
func discoveryRegions(logger *slog.Logger, client *scw.Client) []scw.Region {
	region, ok := client.GetDefaultRegion()
	if !ok {
		return scw.AllRegions
	}

	// check if region is in the list of available regions
	for _, r := range scw.AllRegions {
		if r == region {
			return scw.AllRegions
		}
	}

	logger.Warn("discover: configured default region is unknown", slog.String("region", string(region)))

	return []scw.Region{region}
}

// discoveryZones returns the list of zones to discover resources in.
func discoveryZones(logger *slog.Logger, client *scw.Client) []scw.Zone {
	zone, ok := client.GetDefaultZone()
	if !ok {
		return scw.AllZones
	}

	// check if zone is in the list of available zones
	for _, z := range scw.AllZones {
		if z == zone {
			return scw.AllZones
		}
	}

	logger.Warn("discover: configured default zone is unknown", slog.String("zone", string(zone)))

	return []scw.Zone{zone}
}
