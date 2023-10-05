package model

import "github.com/charmbracelet/bubbles/help"

type AppModel struct {
	LogsModel  LogsModel
	HelpModel  help.Model
	PagerModel PagerModel
}
