package ui

// Focus is the current focus of the application.
// It is used to determine which view should be rendered in the root component.
// It's also used to provide contextual help.

type Focused int

const (
	TableFocused Focused = iota
	SearchFocused
	ConfirmFocused
	JournalFocused
	NumViews // The number of views in the app
)

// nolint:gochecknoglobals
var (
	ViewsSwitchableByTab = []Focused{
		TableFocused,
		SearchFocused,
	}
)
