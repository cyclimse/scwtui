package scenes

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scaleway-dangling/internal/resource"
	"github.com/cyclimse/scaleway-dangling/internal/ui"
	"github.com/cyclimse/scaleway-dangling/internal/ui/confirm"
	"github.com/cyclimse/scaleway-dangling/internal/ui/header"
	"github.com/cyclimse/scaleway-dangling/internal/ui/journal"
	"github.com/cyclimse/scaleway-dangling/internal/ui/search"
	"github.com/cyclimse/scaleway-dangling/internal/ui/table"
)

type Focus int

const (
	tableFocus Focus = iota
	searchFocus
	confirmFocus
	journalFocus
	numViews // The number of views in the app
)

var viewsSwitchableByTab = []Focus{
	tableFocus,
	searchFocus,
}

func Root(state ui.ApplicationState) tea.Model {
	m := Model{
		state:       state,
		header:      header.Header(state),
		search:      search.Search(state),
		table:       table.Table(state.ProjectIDsToNames),
		focusedView: tableFocus,
	}
	m.setFocused(m.focusedView)
	return &m
}

func refreshEvery(state ui.ApplicationState, d time.Duration, skip bool) tea.Cmd {
	return tea.Every(d, func(t time.Time) tea.Msg {
		// If we're skipping the refresh, return an empty slice of resources.
		if skip {
			return []resource.Resource{}
		}

		ctx, cancel := context.WithDeadline(context.Background(), t.Add(d))
		defer cancel()

		resources, err := state.Store.ListAllResources(ctx)
		if err != nil {
			state.Logger.Error("failed to list resources", slog.String("error", err.Error()))
			return []resource.Resource{}
		}

		return resources
	})
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		refreshEvery(m.state, time.Second, false),
		m.header.Init(),
		m.search.Init(),
		m.table.Init(),
		m.confirm.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.header, _ = m.header.Update(msg)
		m.table, _ = m.table.Update(msg)
		return m, nil
	case []resource.Resource:
		// If the search input is dirty, don't update the table.
		if !m.search.Dirty() {
			m.table.UpdateResources(msg)
		}
		return m, refreshEvery(m.state, time.Second, m.search.Dirty())
	case search.SearchResultsMsg:
		m.table.UpdateResources(msg.Resources)
		return m, nil
	case Focus:
		cmd = m.setFocused(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.Keys.Quit):
			if m.focusedView == confirmFocus || m.focusedView == journalFocus {
				m.setFocused(tableFocus)
				return m, nil
			}
			return m, tea.Quit
		// After the user presses enter in the search input, focus the table.
		// This allow quick navigation with the keyboard.
		case msg.String() == "enter" && m.focusedView == searchFocus:
			m.setFocused(tableFocus)
			return m, nil
		case msg.Type == tea.KeyTab:
			cmd = m.focusNext()
			return m, cmd
		}

		// Absorb all keys when those components are focused.
		switch m.focusedView {
		case searchFocus:
			m.search, cmd = m.search.Update(msg)
			if !m.search.Dirty() {
				m.setFocused(tableFocus)
				// we ignore the cmd because we don't want to refresh the table
				return m, nil
			}
			return m, cmd
		case confirmFocus:
			m.confirm, cmd = m.confirm.Update(msg)
			if m.confirm.Deleted() {
				m.setFocused(tableFocus)
			}
			return m, cmd
		case journalFocus:
			m.journal, cmd = m.journal.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, m.state.Keys.Help):
			m.header.ToggleHelp()
			return m, nil
		case key.Matches(msg, m.state.Keys.Delete):
			if m.focusedView == tableFocus {
				m.setFocused(confirmFocus)
				return m, nil
			}
		case key.Matches(msg, m.state.Keys.Logs):
			if m.focusedView == tableFocus {
				canViewLogs := m.table.SelectedResource().CockpitMetadata().CanViewLogs
				if canViewLogs {
					cmd = m.setFocused(journalFocus)
					return m, cmd
				}
			}
		case key.Matches(msg, m.state.Keys.Search):
			cmd = m.setFocused(searchFocus)
			return m, cmd
		case key.Matches(msg, m.state.Keys.ToggleAltView):
			m.table.ToggleAltView()
			return m, nil
		}

		m.setFocused(tableFocus)
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	// Used to pass the blink command to the search component.
	switch m.focusedView {
	case searchFocus:
		m.search, cmd = m.search.Update(msg)
	case confirmFocus:
		m.confirm, cmd = m.confirm.Update(msg)
		if m.confirm.Deleted() {
			cmd = tea.Batch(
				cmd,
				// Wait a second before focusing the table again.
				// This allows the user to see the notification.
				tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
					return tableFocus
				}),
			)
		}
	case journalFocus:
		m.journal, cmd = m.journal.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(m.header.View())
	b.WriteRune('\n')
	switch m.focusedView {
	case searchFocus, tableFocus:
		b.WriteString(m.search.View())
		b.WriteString("\n")
		b.WriteString(m.table.View())
	case confirmFocus:
		b.WriteString("\n\n")
		b.WriteString(lipgloss.PlaceHorizontal(m.table.Width(), lipgloss.Center, m.confirm.View()))
	case journalFocus:
		b.WriteString(m.journal.View())
	}
	return b.String()
}

func (m *Model) focusNext() tea.Cmd {
	return m.setFocused(viewsSwitchableByTab[(int(m.focusedView)+1)%len(viewsSwitchableByTab)])
}

func (m *Model) setFocused(focused Focus) tea.Cmd {
	m.focusedView = focused
	switch m.focusedView {
	case tableFocus:
		m.table.Focus()
	case searchFocus:
		m.table.Blur()
		return m.search.Focus()
	case confirmFocus:
		m.table.Blur()
		m.confirm = confirm.Confirm(m.state, m.table.SelectedResource(), m.table.Width(), m.table.Height())
	case journalFocus:
		m.table.Blur()
		// add exta height to account for table header and border
		const extraHeight = 2
		m.journal = journal.Journal(m.state, m.table.SelectedResource(), m.table.Width()-2, m.table.Height()+extraHeight)
		return m.journal.Init()
	}

	return nil
}

type Model struct {
	state       ui.ApplicationState
	header      header.Model
	search      search.Model
	table       table.Model
	confirm     confirm.Model
	journal     journal.Model
	focusedView Focus
}
