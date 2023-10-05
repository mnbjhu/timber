package model

import (
	"encoding/json"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type LogsModel struct {
	Table table.Model
	Sub   chan AddLog
	Help  help.Model
	Keys  keyMap
}

func listenForActivity(sub chan AddLog) tea.Cmd {
	decoder := json.NewDecoder(os.Stdin)
	return func() tea.Msg {
		for {
			var obj LogMessage
			decoder.Decode(&obj)
			sub <- AddLog{Log: obj}
		}
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan AddLog) tea.Cmd {
	return func() tea.Msg {
		return AddLog(<-sub)
	}
}

func (m LogsModel) Init() tea.Cmd {
	return tea.Batch(
		listenForActivity(m.Sub),
		waitForActivity(m.Sub),
	)
}

func (m LogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case AddLog:
		dateTime, _ := time.Parse(time.RFC3339, msg.Log.Time)
		msg.Log.Time = dateTime.Format("15:04:05")
		rows := append(m.Table.Rows(), table.Row{msg.Log.Time, msg.Log.Level, msg.Log.Prefix, msg.Log.File, msg.Log.Line, msg.Log.Message})
		goToBottom := false
		if m.Table.Cursor() == len(m.Table.Rows())-1 {
			goToBottom = true
		}
		m.Table.SetRows(rows)
		if goToBottom {
			m.Table.GotoBottom()
		}
		return m, waitForActivity(m.Sub)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.Table.SelectedRow()[1]),
			)
		case "?":
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m LogsModel) View() string {
	return baseStyle.Render(m.Table.View()) + "\n" + m.Help.View(m.Keys)
}

var Keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}
