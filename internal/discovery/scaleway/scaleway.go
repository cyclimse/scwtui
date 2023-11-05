package scaleway

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/sync/errgroup"

	"github.com/cyclimse/scwtui/internal/resource"
)

const (
	// MaxRetries is the maximum number of retries for a request.
	MaxRetries = 3
)

var (
	// ErrShouldRetry is returned when the request should be retried.
	ErrShouldRetry = errors.New("should retry")
)

func NewResourceDiscoverer(client *scw.Client, projects []resource.Resource, config *ResourceDiscovererConfig) *ResourceDiscover {
	return &ResourceDiscover{
		client:    client,
		config:    config,
		requested: make(chan requestResources, 1000),
		projects:  projects,
	}
}

type ResourceDiscover struct {
	client *scw.Client
	config *ResourceDiscovererConfig

	projects  []resource.Resource
	requested chan requestResources
}

type ResourceDiscovererConfig struct {
	NumWorkers int
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
	for _, region := range scw.AllRegions {
		region := region // !important
		d.discoverInRegion(region, d.discoverRegistryNamespacesInRegion)
		d.discoverInRegion(region, d.discoverContainersInRegion)
		d.discoverInRegion(region, d.discoverFunctionsInRegion)
		d.discoverInRegion(region, d.discoverRdbInstancesInRegion)
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
					if req.CurrentRetry < MaxRetries {
						d.requested <- req
						continue
					}
					return err
				}
				slog.With("err", err).Error("failed to get resources")
			}

			// slog.With("resource", resources).Info("got resource")
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

	if _, ok := err.(scw.SdkError); !ok {
		slog.With("err", err).Error("got unexpected error")
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
