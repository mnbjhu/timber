package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	m "github.com/mnbjhu/timber/model"
)

var (
	rows    = []table.Row{}
	updated = false
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
	sub   chan AddLog
}

func listenForActivity(sub chan AddLog) tea.Cmd {
	decoder := json.NewDecoder(os.Stdin)
	return func() tea.Msg {
		for {
			var obj m.LogMessage
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

type AddLog struct {
	Log m.LogMessage
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenForActivity(m.sub),
		waitForActivity(m.sub),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case AddLog:
		rows := append(m.table.Rows(), table.Row{msg.Log.Time, msg.Log.Level, msg.Log.Prefix, msg.Log.File, msg.Log.Line, msg.Log.Message})
		m.table.GotoBottom()
		m.table.SetRows(rows)
		return m, waitForActivity(m.sub)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {
	columns := []table.Column{
		{Title: "Time", Width: 10},
		{Title: "Level", Width: 8},
		{Title: "Prefix", Width: 6},
		{Title: "File", Width: 10},
		{Title: "Line", Width: 4},
		{Title: "Message", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	m := model{
		table: t,
		sub:   make(chan AddLog),
	}
	// go readLog(&rows, &m)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func readLog(rows *[]table.Row, t *model) {
	decoder := json.NewDecoder(os.Stdin)

	for {
		var obj m.LogMessage
		err := decoder.Decode(&obj)
		if err != nil {
			// If we've reached the end of input or encountered an error, break the loop
			break
		}
	}
}
