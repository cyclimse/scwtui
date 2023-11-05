package ui

import (
	"log/slog"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// ApplicationState is the state passed to the UI.
// It contains the store and the Scaleway client.
type ApplicationState struct {
	Logger *slog.Logger

	Store   resource.Storer
	Search  resource.Searcher
	Monitor resource.Monitorer

	ScwClient         *scw.Client
	ScwProfileName    string
	ProjectIDsToNames map[string]string

	Keys KeyMap

	Styles Styles

	// The theme to use for syntax highlighting.
	SyntaxHighlighterTheme string
}
