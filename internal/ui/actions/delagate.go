package actions

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/ui"
)

//nolint:gochecknoglobals // will integrate into styles
var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type Action resource.Action

func (a Action) FilterValue() string {
	return a.Name
}

type ActionResultMsg struct {
	Err error
}

func (a Action) Command(state ui.ApplicationState) tea.Cmd {
	return func() tea.Msg {
		return ActionResultMsg{Err: a.Do(context.Background(), resource.NewIndex(state.Store, state.Search), state.ScwClient)}
	}
}

type actionDelegate struct{}

func (d actionDelegate) Height() int { return 1 }

func (d actionDelegate) Spacing() int { return 0 }

func (d actionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d actionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	action, ok := listItem.(Action)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, action.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
