package discovery

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
)

type ResourceDiscoverer interface {
	// Discover discovers resources and sends them to the given channel.
	Discover(ctx context.Context, r chan resource.Resource) error
}
