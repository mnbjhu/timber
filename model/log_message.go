package model

type LogMessage struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Prefix  string `json:"prefix"`
	File    string `json:"file"`
	Line    string `json:"line"`
	Message string `json:"message"`
}
