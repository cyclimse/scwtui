package resource

import (
	"context"
	"time"
)

type Log struct {
	Timestamp time.Time
	Line      string
}

type Monitorer interface {
	// Logs returns the logs of a resource.
	Logs(ctx context.Context, r Resource) ([]Log, error)
}
