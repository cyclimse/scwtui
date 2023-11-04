package scaleway

import (
	"fmt"

	"github.com/cyclimse/scaleway-dangling/internal/resource"
)

func statusPtr[T fmt.Stringer](v T) *resource.Status {
	s := resource.Status(v.String())
	return &s
}
