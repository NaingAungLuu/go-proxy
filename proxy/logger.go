package proxy

import (
	"github.com/charmbracelet/log"
)

type Logger struct{}

func (l *Logger) Write(data []byte) (n int, err error) {
	message := string(data)
	log.Info(message)
	n = len(message)

	return
}

func NewLogger() *Logger {
	return &Logger{}
}
