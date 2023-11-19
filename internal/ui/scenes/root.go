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
	"github.com/cyclimse/scwtui/internal/ui/actions"
	"github.com/cyclimse/scwtui/internal/ui/confirm"
	"github.com/cyclimse/scwtui/internal/ui/describe"
	"github.com/cyclimse/scwtui/internal/ui/header"
	"github.com/cyclimse/scwtui/internal/ui/journal"
	"github.com/cyclimse/scwtui/internal/ui/search"
	"github.com/cyclimse/scwtui/internal/ui/table"
)

const refreshInterval = 3 * time.Second

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

type refreshPeriodicallyMsg struct {
	Resources []resource.Resource
}

func refreshEvery(state ui.ApplicationState, d time.Duration) tea.Cmd {
	return tea.Every(d, func(t time.Time) tea.Msg {
		ctx, cancel := context.WithDeadline(context.Background(), t.Add(d))
		defer cancel()

		resources, err := state.Store.ListAllResources(ctx)
		if err != nil {
			state.Logger.Error("failed to list resources", slog.String("error", err.Error()))
			return refreshPeriodicallyMsg{Resources: []resource.Resource{}}
		}

		return refreshPeriodicallyMsg{Resources: resources}
	})
}

type refreshOnceMsg struct {
	Resources []resource.Resource
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		// on startup, refresh the resources quicker than the default interval.
		refreshEvery(m.state, time.Second),
		m.header.Init(),
		m.search.Init(),
		m.table.Init(),
	)
}

//nolint:funlen // root component, hard to split.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.updateWindowsResize(msg), nil
	case refreshPeriodicallyMsg:
		// if the search input is dirty, don't update the table.
		if !m.search.Dirty() {
			m.table.UpdateResources(msg.Resources)
		}
		return m, refreshEvery(m.state, refreshInterval)
	case refreshOnceMsg:
		m.table.UpdateResources(msg.Resources)
		return m, nil
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
				cmd = m.setFocused(ui.TableFocused)
				return m, cmd
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

		// absorb all keys when those components are focused.
		if m.focused != ui.TableFocused {
			return m.updateFocusedOnKeyMsg(msg)
		}

		switch {
		case key.Matches(msg, m.state.Keys.Search):
			cmd = m.setFocused(ui.SearchFocused)
			return m, cmd
		case key.Matches(msg, m.state.Keys.Describe):
			cmd = m.setFocused(ui.DescribeFocused)
			return m, cmd
		case key.Matches(msg, m.state.Keys.Logs):
			canViewLogs := m.table.SelectedResource().CockpitMetadata().CanViewLogs
			if canViewLogs {
				cmd = m.setFocused(ui.JournalFocused)
				return m, cmd
			}
		case key.Matches(msg, m.state.Keys.Delete):
			cmd = m.setFocused(ui.ConfirmFocused)
			return m, cmd
		case key.Matches(msg, m.state.Keys.Actions):
			_, ok := m.table.SelectedResource().(resource.Actionable)
			if ok {
				cmd = m.setFocused(ui.ActionsFocused)
				return m, cmd
			}
		}

		m.setFocused(ui.TableFocused)
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	// absorb all messages when those components are focused.
	// for instance, the "blink" message from the text input component
	// will get forwarded to the search component.
	return m.updateFocusedOnAnyMsg(msg)
}

func (m Model) updateFocusedOnKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.focused {
	case ui.SearchFocused:
		m.search, cmd = m.search.Update(msg)
		if !m.search.Dirty() {
			cmd = m.setFocused(ui.TableFocused)
		}
	case ui.DescribeFocused:
		m.describe, cmd = m.describe.Update(msg)
	case ui.ConfirmFocused:
		m.confirm, cmd = m.confirm.Update(msg)
	case ui.JournalFocused:
		m.journal, cmd = m.journal.Update(msg)
	case ui.ActionsFocused:
		m.actions, cmd = m.actions.Update(msg)
	}

	return m, cmd
}

func (m Model) updateFocusedOnAnyMsg(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

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
	case ui.ActionsFocused:
		if _, ok := msg.(actions.ActionResultMsg); ok {
			cmd = tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
				return ui.TableFocused
			})
		} else {
			m.actions, cmd = m.actions.Update(msg)
		}
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
	case ui.ConfirmFocused: // confirm is a modal, so we need to render it on top of the table.
		b.WriteString("\n\n")
		b.WriteString(lipgloss.PlaceHorizontal(m.table.Width(), lipgloss.Center, m.confirm.View()))
	case ui.JournalFocused:
		b.WriteString(m.journal.View())
	case ui.ActionsFocused: // actions is a modal, so we need to render it on top of the table.
		b.WriteString("\n\n")
		b.WriteString(lipgloss.PlaceHorizontal(m.table.Width(), lipgloss.Center, m.actions.View()))
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
	var cmd tea.Cmd

	switch focused {
	case ui.TableFocused:
		m.table.Focus()

		// special case: if we come back from the search, we need to update the
		// resources so that the table is up to date.
		// this will be handled by the refreshEvery command, but we need to
		// do it quicker than the default interval to make the UI more responsive.
		if m.focused == ui.SearchFocused && !m.search.Dirty() {
			cmd = func() tea.Msg {
				ctx := context.Background()
				resources, err := m.state.Store.ListAllResources(ctx)
				if err != nil {
					m.state.Logger.Error("failed to list resources", slog.String("error", err.Error()))
					return refreshOnceMsg{Resources: []resource.Resource{}}
				}
				return refreshOnceMsg{Resources: resources}
			}
		}
	case ui.SearchFocused:
		m.table.Blur()
		cmd = m.search.Focus()
	case ui.DescribeFocused:
		m.table.Blur()
		m.describe = describe.Describe(m.state, m.table.SelectedResource(), m.table.Width()-fullViewExtraPaddding, m.table.Height()+fullViewExtraHeight)
		cmd = m.describe.Init()
	case ui.ConfirmFocused:
		m.table.Blur()
		m.confirm = confirm.Confirm(m.state, m.table.SelectedResource(), m.table.Width(), m.table.Height())
		cmd = m.confirm.Init()
	case ui.JournalFocused:
		m.table.Blur()
		m.journal = journal.Journal(m.state, m.table.SelectedResource(), m.table.Width()-fullViewExtraPaddding, m.table.Height()+fullViewExtraHeight)
		cmd = m.journal.Init()
	case ui.ActionsFocused:
		m.table.Blur()
		m.actions = actions.Actions(m.state, m.table.SelectedResource().(resource.Actionable), m.table.Width(), m.table.Height())
		cmd = m.actions.Init()
	}

	m.focused = focused
	m.header.SetFocused(focused)

	return cmd
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
	actions  actions.Model
}
