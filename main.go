package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mnbjhu/timber/model"
)

func main() {
	columns := []table.Column{
		{Title: "Time", Width: 10},
		{Title: "Level", Width: 10},
		{Title: "Prefix", Width: 6},
		{Title: "File", Width: 10},
		{Title: "Line", Width: 4},
		{Title: "Message", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
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
	m := model.LogsModel{
		Table: t,
		Sub:   make(chan model.AddLog),
		Help:  help.New(),
		Keys:  model.Keys,
	}
	// go readLog(&rows, &m)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
