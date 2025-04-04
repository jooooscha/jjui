package status

import (
	"github.com/charmbracelet/bubbles/key"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/context"
)

var cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "dismiss"))

type Model struct {
	context context.AppContext
	spinner spinner.Model
	help    help.Model
	keyMap  help.KeyMap
	command string
	running bool
	output  string
	error   error
	width   int
	mode    string
}

const CommandClearDuration = 3 * time.Second

type clearMsg string

func (m *Model) Width() int {
	return m.width
}

func (m *Model) Height() int {
	return 1
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m *Model) SetHeight(int) {}
func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case clearMsg:
		if m.command == string(msg) {
			m.command = ""
			m.error = nil
			m.output = ""
		}
		return m, nil
	case common.CommandRunningMsg:
		m.command = string(msg)
		m.running = true
		return m, m.spinner.Tick
	case common.CommandCompletedMsg:
		m.running = false
		m.output = msg.Output
		m.error = msg.Err
		if m.error == nil {
			commandToBeCleared := m.command
			return m, tea.Tick(CommandClearDuration, func(time.Time) tea.Msg {
				return clearMsg(commandToBeCleared)
			})
		}
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, cancel) && m.error != nil:
			m.error = nil
			m.output = ""
			m.command = ""
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m *Model) View() string {
	s := common.DefaultPalette.Normal.Render(" ")
	if m.running {
		s = common.DefaultPalette.Normal.Render(m.spinner.View())
	} else if m.error != nil {
		s = common.DefaultPalette.StatusError.Render("✗ ")
	} else if m.command != "" {
		s = common.DefaultPalette.StatusSuccess.Render("✓ ")
	} else {
		s = m.help.View(m.keyMap)
	}
	ret := common.DefaultPalette.Normal.Render(m.command)
	mode := common.DefaultPalette.StatusMode.Width(10).Render("", m.mode)
	ret = lipgloss.JoinHorizontal(lipgloss.Left, mode, " ", s, ret)
	if m.error != nil {
		k := cancel.Help().Key
		return lipgloss.JoinVertical(0,
			ret,
			common.DefaultPalette.StatusError.Render(strings.Trim(m.output, "\n")),
			common.DefaultPalette.ChangeId.Render("press ", k, " to dismiss"))
	}
	return ret
}

func (m *Model) SetHelp(keyMap help.KeyMap) {
	m.keyMap = keyMap
}

func (m *Model) SetMode(mode string) {
	m.mode = mode
}

func New(context context.AppContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	h := help.New()
	h.Styles.ShortKey = common.DefaultPalette.ChangeId
	h.Styles.ShortDesc = common.DefaultPalette.Dimmed
	h.ShortSeparator = " "
	return Model{
		context: context,
		spinner: s,
		help:    h,
		command: "",
		running: false,
		output:  "",
		keyMap:  nil,
	}
}
