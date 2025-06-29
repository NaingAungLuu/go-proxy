package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Panel struct {
	titleBar TitleBar
	vp       viewport.Model
}

var (
	panelStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Margin(1).
		PaddingLeft(1).
		PaddingRight(1)
)

func (p *Panel) Init() tea.Cmd {
	return nil
}

func (p *Panel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	p.vp, cmd = p.vp.Update(message)
	cmds = append(cmds, cmd)
	return p, tea.Batch(cmds...)
}

func (p *Panel) View() string {
	return p.titleBar.Render() +
		p.vp.View()
}

func (p *Panel) SetFrameSize(w, h int) {
	_, titleHeight := p.titleBar.GetFrameSize()
	p.vp.Width = w
	p.vp.Height = h - titleHeight
	p.vp.Style.Width(w)
	p.vp.Style.Height(h - titleHeight)
}

func (p *Panel) SetContent(content string) {
	p.vp.SetContent(content)
}

func (p *Panel) ScrollDown(n int) (lines []string) {
	return p.vp.ScrollDown(n)
}

func NewPanel(title string) *Panel {
	titleBar := NewTitleBar(title)
	viewport := viewport.New(20, 20)
	viewport.Style = panelStyle
	return &Panel{titleBar: *titleBar, vp: viewport}
}
