package scenes

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/ui"
	"github.com/cyclimse/scwtui/internal/ui/confirm"
	"github.com/cyclimse/scwtui/internal/ui/describe"
	"github.com/cyclimse/scwtui/internal/ui/header"
	"github.com/cyclimse/scwtui/internal/ui/journal"
	"github.com/cyclimse/scwtui/internal/ui/search"
	"github.com/cyclimse/scwtui/internal/ui/table"
)

func Root(state ui.ApplicationState) tea.Model {
	m := Model{
		state:   state,
		focused: ui.TableFocused,

		header: header.Header(ui.TableFocused, state),
		search: search.Search(state),
		table:  table.Table(state),
	}
	m.setFocused(m.focused)
	return &m
}

func refreshEvery(state ui.ApplicationState, d time.Duration, skip bool) tea.Cmd {
	return tea.Every(d, func(t time.Time) tea.Msg {
		// if we're skipping the refresh, return an empty slice of resources.
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
	)
}

//nolint:funlen,gocognit // will address this later.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.updateWindowsResize(msg), nil
	case []resource.Resource:
		// if the search input is dirty, don't update the table.
		if !m.search.Dirty() {
			m.table.UpdateResources(msg)
		}
		return m, refreshEvery(m.state, time.Second, m.search.Dirty())
	case search.ResultsMsg:
		m.table.UpdateResources(msg.Resources)
		return m, nil
	case ui.Focused:
		// this is used to set the focus asynchroniously
		// e.g. after a resource is deleted, we want to let the user see the
		// confirmation message before focusing the table again.
		cmd = m.setFocused(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.state.Keys.Quit):
			if m.focused != ui.TableFocused {
				m.setFocused(ui.TableFocused)
				return m, nil
			}
			return m, tea.Quit
		case msg.Type == tea.KeyEnter:
			// allows natural navigation after a search.
			if m.focused == ui.SearchFocused {
				m.setFocused(ui.TableFocused)
				return m, nil
			}
		case msg.Type == tea.KeyTab:
			cmd = m.focusNext()
			return m, cmd
		}

		// Absorb all keys when those components are focused.
		switch m.focused {
		case ui.SearchFocused:
			m.search, cmd = m.search.Update(msg)
			if !m.search.Dirty() {
				m.setFocused(ui.TableFocused)
				// we ignore the cmd because we don't want to refresh the table
				return m, nil
			}
			return m, cmd
		case ui.DescribeFocused:
			m.describe, cmd = m.describe.Update(msg)
			return m, cmd
		case ui.ConfirmFocused:
			m.confirm, cmd = m.confirm.Update(msg)
			return m, cmd
		case ui.JournalFocused:
			m.journal, cmd = m.journal.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, m.state.Keys.Search):
			cmd = m.setFocused(ui.SearchFocused)
			return m, cmd
		case key.Matches(msg, m.state.Keys.Describe):
			if m.focused == ui.TableFocused {
				cmd = m.setFocused(ui.DescribeFocused)
				return m, cmd
			}
		case key.Matches(msg, m.state.Keys.Delete):
			if m.focused == ui.TableFocused {
				cmd = m.setFocused(ui.ConfirmFocused)
				return m, cmd
			}
		case key.Matches(msg, m.state.Keys.Logs):
			if m.focused == ui.TableFocused {
				canViewLogs := m.table.SelectedResource().CockpitMetadata().CanViewLogs
				if canViewLogs {
					cmd = m.setFocused(ui.JournalFocused)
					return m, cmd
				}
			}
		}

		m.setFocused(ui.TableFocused)
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	// Used to pass the blink command to the search component.
	switch m.focused {
	case ui.SearchFocused:
		m.search, cmd = m.search.Update(msg)
	case ui.DescribeFocused:
		m.describe, cmd = m.describe.Update(msg)
	case ui.ConfirmFocused:
		m.confirm, cmd = m.confirm.Update(msg)
		if m.confirm.Deleted() {
			cmd = tea.Batch(
				cmd,
				// wait a second before focusing the table again.
				// this allows the user to see the notification.
				tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
					return ui.TableFocused
				}),
			)
		}
	case ui.JournalFocused:
		m.journal, cmd = m.journal.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(m.header.View())
	b.WriteRune('\n')
	switch m.focused {
	case ui.TableFocused, ui.SearchFocused:
		b.WriteString(m.search.View())
		b.WriteString("\n")
		b.WriteString(m.table.View())
	case ui.DescribeFocused:
		b.WriteString(m.describe.View())
	case ui.ConfirmFocused:
		b.WriteString("\n\n")
		b.WriteString(lipgloss.PlaceHorizontal(m.table.Width(), lipgloss.Center, m.confirm.View()))
	case ui.JournalFocused:
		b.WriteString(m.journal.View())
	}
	return b.String()
}

func (m *Model) focusNext() tea.Cmd {
	return m.setFocused(ui.ViewsSwitchableByTab[(int(m.focused)+1)%len(ui.ViewsSwitchableByTab)])
}

const (
	fullViewExtraHeight   = 2
	fullViewExtraPaddding = 2
)

func (m *Model) setFocused(focused ui.Focused) tea.Cmd {
	m.focused = focused
	m.header.SetFocused(focused)

	switch m.focused {
	case ui.TableFocused:
		m.table.Focus()
	case ui.SearchFocused:
		m.table.Blur()
		return m.search.Focus()
	case ui.DescribeFocused:
		m.table.Blur()
		m.describe = describe.Describe(m.state, m.table.SelectedResource(), m.table.Width()-fullViewExtraPaddding, m.table.Height()+fullViewExtraHeight)
		return m.describe.Init()
	case ui.ConfirmFocused:
		m.table.Blur()
		m.confirm = confirm.Confirm(m.state, m.table.SelectedResource(), m.table.Width(), m.table.Height())
		return m.confirm.Init()
	case ui.JournalFocused:
		m.table.Blur()
		m.journal = journal.Journal(m.state, m.table.SelectedResource(), m.table.Width()-fullViewExtraPaddding, m.table.Height()+fullViewExtraHeight)
		return m.journal.Init()
	}

	return nil
}

func (m Model) updateWindowsResize(msg tea.WindowSizeMsg) Model {
	m.header.SetWidth(msg.Width)
	m.table.SetDimensions(msg.Width, msg.Height)

	w := m.table.Width()
	h := m.table.Height()

	// Resize the other components to the table's dimensions.
	m.describe.SetDimensions(w-fullViewExtraPaddding, h+fullViewExtraHeight)
	m.journal.SetDimensions(w-fullViewExtraPaddding, h+fullViewExtraHeight)
	return m
}

type Model struct {
	state   ui.ApplicationState
	focused ui.Focused

	header   header.Model
	search   search.Model
	describe describe.Model
	table    table.Model
	confirm  confirm.Model
	journal  journal.Model
}
