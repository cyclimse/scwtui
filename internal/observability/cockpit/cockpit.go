package cockpit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/grafana/loki/pkg/loghttp"
	sdk "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	tokenName = "scaleway-dangling"

	// queryTemplateWithName is the template used to query logs from Loki.
	queryTemplateWithName = `{resource_name="%s", resource_type="%s"} |~ "^{.*}$" | json | line_format "{{.message}}"`

	// queryTemplateWithID is the template used to query logs from Loki.
	queryTemplateWithID = `{resource_id="%s", resource_type="%s"} |~ "^{.*}$" | json | line_format "{{.message}}"`
)

var (
	// ErrCockpitNotActivated is returned when the cockpit is not activated for a project.
	ErrCockpitNotActivated = fmt.Errorf("cockpit: cockpit not activated")
)

func NewCockpit(logger *slog.Logger, scwClient *scw.Client) *Cockpit {
	return &Cockpit{
		logger:                logger,
		scwClient:             scwClient,
		httpClient:            &http.Client{},
		tokenPerProject:       make(map[string]sdk.Token),
		lokiAddressPerProject: make(map[string]string),
	}
}

type Cockpit struct {
	logger                *slog.Logger
	scwClient             *scw.Client
	httpClient            *http.Client
	tokenPerProject       map[string]sdk.Token
	lokiAddressPerProject map[string]string
}

func (c *Cockpit) Logs(ctx context.Context, r resource.Resource) ([]resource.Log, error) {
	metadata := r.Metadata()
	cockpitMetadata := r.CockpitMetadata()

	if !cockpitMetadata.CanViewLogs {
		return nil, nil
	}

	address, err := c.lokiAddressForProject(metadata.ProjectID)
	if err != nil {
		return nil, err
	}

	token, err := c.tokenForProject(metadata.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("cockpit: unable to get token for project %s: %w", metadata.ProjectID, err)
	}

	query := buildQuery(r)
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	limit := 5000

	queryURL := buildQueryURL(address, query, start, end, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("cockpit: unable to create request: %w", err)
	}

	req.Header.Set("X-Token", *token.SecretKey)
	req.Header.Set("X-Datasource", "product")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cockpit: unable to execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cockpit: unexpected status code: %d", resp.StatusCode)
	}

	return c.parseLogs(resp)
}

func (c *Cockpit) parseLogs(resp *http.Response) ([]resource.Log, error) {
	var r loghttp.QueryResponse
	err := json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("cockpit: unable to decode response: %w", err)
	}

	switch r.Data.ResultType {
	case loghttp.ResultTypeStream:
		streams := r.Data.Result.(loghttp.Streams)
		logs := make([]resource.Log, 0, len(streams))

		for _, stream := range streams {
			for _, entry := range stream.Entries {
				logs = append(logs, resource.Log{
					Timestamp: entry.Timestamp,
					Line:      entry.Line,
				})
			}
		}

		return logs, nil
	default:
		return nil, fmt.Errorf("cockpit: unexpected result type: %s", r.Data.ResultType)
	}
}

func buildQuery(r resource.Resource) string {
	cockpitMetadata := r.CockpitMetadata()

	if cockpitMetadata.ResourceID != "" {
		return fmt.Sprintf(queryTemplateWithID, cockpitMetadata.ResourceID, cockpitMetadata.ResourceType)
	}
	return fmt.Sprintf(queryTemplateWithName, cockpitMetadata.ResourceName, cockpitMetadata.ResourceType)
}

func buildQueryURL(address, query string, start, end time.Time, limit int) string {
	queryURL := address + "/loki/api/v1/query_range?query=" + url.QueryEscape(query)
	queryURL += "&limit=" + strconv.Itoa(limit)
	queryURL += "&start=" + url.QueryEscape(start.UTC().Format(time.RFC3339))
	queryURL += "&end=" + url.QueryEscape(end.UTC().Format(time.RFC3339))
	return queryURL
}

func (c *Cockpit) lokiAddressForProject(projectID string) (string, error) {
	if address, ok := c.lokiAddressPerProject[projectID]; ok {
		return address, nil
	}

	api := sdk.NewAPI(c.scwClient)

	cockpit, err := api.GetCockpit(&sdk.GetCockpitRequest{
		ProjectID: projectID,
	})
	if err != nil {
		var resourceNotFoundError *scw.ResourceNotFoundError
		if errors.As(err, &resourceNotFoundError) {
			return "", ErrCockpitNotActivated
		}
		return "", fmt.Errorf("cockpit: unable to get cockpit for project %s: %w", projectID, err)
	}

	c.lokiAddressPerProject[projectID] = cockpit.Endpoints.LogsURL
	return cockpit.Endpoints.LogsURL, nil
}

func (c *Cockpit) tokenForProject(projectID string) (*sdk.Token, error) {
	if token, ok := c.tokenPerProject[projectID]; ok {
		return &token, nil
	}

	// we attempt to remove the previously generated token with the same name
	// this avoids having a lot of dangling tokens
	err := c.deleteTokenWithExistingName(projectID)
	if err != nil {
		c.logger.Warn("cockpit: unable to delete token with existing name",
			slog.String("project_id", projectID),
			slog.String("error", err.Error()))
	}

	api := sdk.NewAPI(c.scwClient)

	// we have to recreate a token as the secret key is not returned by the API
	// it's only available when creating the token
	token, err := api.CreateToken(&sdk.CreateTokenRequest{
		ProjectID: projectID,
		Name:      tokenName,
		Scopes: &sdk.TokenScopes{
			QueryLogs: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("cockpit: unable to create token for project %s: %w", projectID, err)
	}

	if token == nil || token.SecretKey == nil {
		return nil, fmt.Errorf("cockpit: unable to get secret key for token")
	}

	c.tokenPerProject[projectID] = *token
	return token, nil
}

func (c *Cockpit) deleteTokenWithExistingName(projectID string) error {
	api := sdk.NewAPI(c.scwClient)

	req, err := api.ListTokens(&sdk.ListTokensRequest{
		ProjectID: projectID,
	}, scw.WithAllPages())
	if err != nil {
		return fmt.Errorf("cockpit: unable to list tokens for project %s: %w", projectID, err)
	}

	for _, token := range req.Tokens {
		if token == nil {
			continue
		}

		if token.Name == tokenName {
			err := api.DeleteToken(&sdk.DeleteTokenRequest{
				TokenID: token.ID,
			})
			if err != nil {
				return fmt.Errorf("cockpit: unable to delete token for project %s: %w", projectID, err)
			}
		}
	}

	return nil
}
