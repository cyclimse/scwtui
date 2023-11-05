package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scaleway-dangling/internal/resource"
	"github.com/cyclimse/scaleway-dangling/internal/ui/header"
)

const (
	defaultWidth = 80
)

// nolint: gochecknoglobals
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func Table(projectsIDsToNames map[string]string) Model {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	b := NewBuilder(s)

	return Model{
		builder:            b,
		lastWidthBuilt:     defaultWidth,
		projectsIDsToNames: projectsIDsToNames,
	}
}

func (m Model) Init() tea.Cmd { return nil }

const (
	additionalHorizontalPadding = 8
	tableHeaderHeight           = 4
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width != m.lastWidthBuilt {
			m.lastWidthBuilt = msg.Width - baseStyle.GetHorizontalFrameSize() - additionalHorizontalPadding
			m.lastHeight = msg.Height - baseStyle.GetVerticalFrameSize() - header.MaxHeight - tableHeaderHeight
			m.rebuildTable()
			m.table.SetWidth(m.lastWidthBuilt + additionalHorizontalPadding)
			m.table.SetHeight(m.lastHeight)
		}
		return m, cmd
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyF2:
			m.ToggleAltView()
			return m, cmd
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) rebuildTable() {
	previous := struct {
		width  int
		height int
		cursor int
	}{
		width:  m.table.Width(),
		height: m.table.Height(),
		cursor: m.table.Cursor(),
	}

	m.table = m.builder.Build(BuildParams{
		Width:             m.lastWidthBuilt,
		AltView:           m.showingAltView,
		Resources:         m.resources,
		ProjectIDsToNames: m.projectsIDsToNames,
	})

	m.table.SetWidth(previous.width)
	m.table.SetHeight(previous.height)
	m.table.SetCursor(previous.cursor)
}

func (m *Model) UpdateResources(resources []resource.Resource) {
	m.resources = resources
	m.rebuildTable()
}

func (m *Model) SelectedResource() resource.Resource {
	return m.resources[m.table.Cursor()]
}

func (m Model) View() string {
	return baseStyle.Render(m.table.View())
}

func (m *Model) ToggleAltView() {
	m.showingAltView = !m.showingAltView
	m.rebuildTable()
}

func (m *Model) Focus() {
	m.table.Focus()
}

func (m *Model) Blur() {
	m.table.Blur()
}

func (m Model) Width() int {
	return m.table.Width()
}

func (m Model) Height() int {
	return m.table.Height()
}

type Model struct {
	table   table.Model
	builder *Build

	resources          []resource.Resource
	projectsIDsToNames map[string]string

	lastWidthBuilt int
	lastHeight     int

	showingAltView bool
}
