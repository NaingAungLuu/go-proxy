package ui

import (
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type UILogger struct {
	program tea.Program
}

func (l *UILogger) Log(request *http.Request) {
	l.program.Send(LogEvent{Request: request})
}

func NewUiLogger(program tea.Program) UILogger {
	return UILogger{program: program}
}
