package ui

import "github.com/charmbracelet/lipgloss"

type TitleBar struct {
	title string
	style lipgloss.Style
}

func (t *TitleBar) Render(content ...string) string {
	return t.style.Render(content...)
}

func (t *TitleBar) GetFrameSize() (x, y int) {
	return t.style.GetFrameSize()
}

func NewTitleBar(title string) TitleBar {
	titleStyle := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Left).
		Foreground(lipgloss.NoColor{}).
		Background(lipgloss.Color("62")).
		Padding(0, 1).
		MarginLeft(2).
		MarginTop(1)

	return TitleBar{
		title: title,
		style: titleStyle,
	}
}
