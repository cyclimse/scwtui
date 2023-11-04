package demo

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/cyclimse/scaleway-dangling/internal/resource"
	"github.com/cyclimse/scaleway-dangling/internal/testhelpers"
	"github.com/scaleway/scaleway-sdk-go/namegenerator"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func NewDiscovery() *Discovery {
	allLocalities := make([]resource.Locality, 0, len(scw.AllRegions)+len(scw.AllZones)+1) // +1 for global locality
	for _, region := range scw.AllRegions {
		allLocalities = append(allLocalities, resource.Region(region))
	}
	for _, zone := range scw.AllZones {
		allLocalities = append(allLocalities, resource.Zone(zone))
	}
	allLocalities = append(allLocalities, resource.Global)

	numFakeProjects := gofakeit.IntRange(1, 10)
	fakeProjects := make([]resource.Resource, 0, numFakeProjects)

	for i := 0; i < numFakeProjects; i++ {
		projectID := gofakeit.UUID()
		description := gofakeit.Sentence(10)

		fakeProject := &testhelpers.MockResource{
			MetadataValue: resource.Metadata{
				Name:        namegenerator.GetRandomName(),
				ID:          projectID,
				ProjectID:   projectID,
				Description: &description,
				Tags:        nil,
				Type:        resource.TypeProject,
			},
		}
		fakeProjects = append(fakeProjects, fakeProject)
	}

	return &Discovery{
		allLocalities: allLocalities,
		fakeProjects:  fakeProjects,
	}
}

type Discovery struct {
	allLocalities []resource.Locality
	fakeProjects  []resource.Resource
}

func (d *Discovery) Discover(_ context.Context, ch chan resource.Resource) error {
	ch <- d.resource()
	return nil
}

func (d *Discovery) resource() resource.Resource {
	chosenProject := d.fakeProjects[gofakeit.IntRange(0, len(d.fakeProjects)-1)]

	numTags := gofakeit.IntRange(0, 10)
	tags := make([]string, 0, numTags)
	for i := 0; i < numTags; i++ {
		tags = append(tags, gofakeit.Word())
	}

	rType := resource.Type(gofakeit.IntRange(0, int(resource.NumberOfResourceTypes)-1))
	locality := d.allLocalities[gofakeit.IntRange(0, len(d.allLocalities)-1)]

	return &testhelpers.MockResource{
		MetadataValue: resource.Metadata{
			Name:      namegenerator.GetRandomName(),
			ID:        gofakeit.UUID(),
			ProjectID: chosenProject.Metadata().ID,
			Tags:      tags,
			Type:      rType,
			Locality:  locality,
		},
	}
}

func (d *Discovery) Projects() []resource.Resource {
	return d.fakeProjects
}
