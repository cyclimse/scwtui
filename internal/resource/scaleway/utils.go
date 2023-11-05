package scaleway

import (
	"fmt"

	"github.com/cyclimse/scwtui/internal/resource"
)

func statusPtr[T fmt.Stringer](v T) *resource.Status {
	s := resource.Status(v.String())
	return &s
}
