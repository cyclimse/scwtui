package bleve

import (
	"context"

	"github.com/blevesearch/bleve/v2"
	"github.com/cyclimse/scwtui/internal/resource"
)

func NewSearch(projectIDsToNames map[string]string) (*Search, error) {
	mapping := bleve.NewIndexMapping()
	mapping.DefaultAnalyzer = "standard"

	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, err
	}

	return &Search{
		index:             index,
		projectIDsToNames: projectIDsToNames,
	}, nil
}

type Search struct {
	index             bleve.Index
	projectIDsToNames map[string]string
}

func (s *Search) Search(ctx context.Context, query string) ([]string, error) {
	q := bleve.NewQueryStringQuery(query)
	search := bleve.NewSearchRequest(q)
	searchResults, err := s.index.SearchInContext(ctx, search)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, hit := range searchResults.Hits {
		ids = append(ids, hit.ID)
	}

	return ids, nil
}

type indexResource struct {
	// Needs to be embedded for the field selector to work
	resource.Resource `json:""`
	resource.Metadata `json:""`
	Type              string `json:"type"`
	Project           string `json:"project"`
	Zone              string `json:"zone"`
	Region            string `json:"region"`
}

func (s *Search) Index(r resource.Resource) error {
	metadata := r.Metadata()

	region := ""
	zone := ""
	if metadata.Locality.IsRegion() {
		region = metadata.Locality.String()
	}
	if metadata.Locality.IsZone() {
		zone = metadata.Locality.String()
	}

	return s.index.Index(metadata.ID, indexResource{
		Resource: r,
		Metadata: metadata,
		Type:     metadata.Type.String(),
		Project:  s.projectIDsToNames[metadata.ProjectID],
		Zone:     zone,
		Region:   region,
	})
}
