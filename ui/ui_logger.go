package ui

import (
	"go-proxy/proxy"

	tea "github.com/charmbracelet/bubbletea"
)

type UILogger struct {
	program tea.Program
}

func (l *UILogger) Log(log proxy.RequestLog) {
	l.program.Send(LogEvent{log: log})
}

func NewUiLogger(program tea.Program) UILogger {
	return UILogger{program: program}
}
