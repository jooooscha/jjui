package preview

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/context"
)

type Model struct {
	tag     int
	view    viewport.Model
	help    help.Model
	width   int
	height  int
	content string
	context context.AppContext
	keyMap  config.KeyMappings[key.Binding]
}

const DebounceTime = 10 * time.Millisecond

type refreshPreviewContentMsg struct {
	Tag int
}

type updatePreviewContentMsg struct {
	Content string
}

func (m *Model) ShortHelp() []key.Binding {
	return []key.Binding{
		m.view.KeyMap.HalfPageUp,
		m.view.KeyMap.HalfPageDown,
	}
}

func (m *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

func (m *Model) Width() int {
	return m.width
}

func (m *Model) Height() int {
	return m.height
}

func (m *Model) SetWidth(w int) {
	content := lipgloss.NewStyle().MaxWidth(w - 2).Render(m.content)
	m.view.SetContent(content)
	m.view.Width = w
	m.width = w
}

func (m *Model) SetHeight(h int) {
	m.view.Height = h
	m.height = h
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updatePreviewContentMsg:
		m.content = msg.Content
		content := lipgloss.NewStyle().MaxWidth(m.Width() - 4).Render(msg.Content)
		m.view.SetContent(content)
		m.view.GotoTop()
	case common.SelectionChangedMsg, common.RefreshMsg:
		m.tag++
		tag := m.tag
		return m, tea.Tick(DebounceTime, func(t time.Time) tea.Msg {
			return refreshPreviewContentMsg{Tag: tag}
		})
	case refreshPreviewContentMsg:
		if m.tag == msg.Tag {
			switch msg := m.context.SelectedItem().(type) {
			case context.SelectedFile:
				return m, func() tea.Msg {
					output, _ := m.context.RunCommandImmediate(jj.Diff(msg.ChangeId, msg.File))
					return updatePreviewContentMsg{Content: string(output)}
				}
			case context.SelectedRevision:
				return m, func() tea.Msg {
					output, _ := m.context.RunCommandImmediate(jj.Show(msg.ChangeId))
					return updatePreviewContentMsg{Content: string(output)}
				}
			case context.SelectedOperation:
				return m, func() tea.Msg {
					output, _ := m.context.RunCommandImmediate(jj.OpShow(msg.OperationId))
					return updatePreviewContentMsg{Content: string(output)}
				}
			}
		}
	}
	var cmd tea.Cmd
	m.view, cmd = m.view.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.view.View()
}

func viewPortKeyMap(km config.KeyMappings[key.Binding]) viewport.KeyMap {
	return viewport.KeyMap{
		PageDown:     key.NewBinding(key.WithDisabled()),
		PageUp:       key.NewBinding(key.WithDisabled()),
		HalfPageUp:   km.Preview.HalfPageUp,
		HalfPageDown: km.Preview.HalfPageDown,
		Up:           km.Preview.ScrollUp,
		Down:         km.Preview.ScrollDown,
	}
}

func New(context context.AppContext) Model {
	view := viewport.New(0, 0)
	view.Style = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	keyMap := context.KeyMap()
	view.KeyMap = viewPortKeyMap(keyMap)
	return Model{
		context: context,
		keyMap:  keyMap,
		view:    view,
		help:    help.New(),
	}
}
