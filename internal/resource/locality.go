package resource

import (
	"fmt"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	// Global is the global locality.
	Global = GlobalLocality("global")
)

type Locality interface {
	IsRegion() bool
	IsZone() bool
	fmt.Stringer
}

type Region scw.Region

func (r Region) IsRegion() bool {
	return true
}

func (r Region) IsZone() bool {
	return false
}

func (r Region) String() string {
	return string(r)
}

type Zone scw.Zone

func (z Zone) IsRegion() bool {
	return false
}

func (z Zone) IsZone() bool {
	return true
}

func (z Zone) String() string {
	return string(z)
}

type GlobalLocality string

func (g GlobalLocality) IsRegion() bool {
	return false
}

func (g GlobalLocality) IsZone() bool {
	return false
}

func (g GlobalLocality) String() string {
	return "global"
}
