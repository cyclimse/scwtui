package ui

import (
	"github.com/cyclimse/scwtui/internal/resource"
)

func ApplyIDsFilter(previous []resource.Resource, filter resource.SetOfIDs) []resource.Resource {
	resources := make([]resource.Resource, 0, len(filter))

	for _, r := range previous {
		_, ok := filter[r.Metadata().ID]
		if ok {
			resources = append(resources, r)
		}
	}

	return resources
}
