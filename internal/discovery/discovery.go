package discovery

import (
	"context"

	"github.com/cyclimse/scaleway-dangling/internal/resource"
)

type ResourceDiscoverer interface {
	// Discover discovers resources and sends them to the given channel.
	Discover(context.Context, chan resource.Resource) error
}
