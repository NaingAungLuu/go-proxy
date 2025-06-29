package ui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InfoPanel struct {
	panel          Panel
	destinationUrl string
	port           int
}

type InfoUpdate struct {
	destinationUrl string
	port           int
}

/**
** Style Configurations
**/
var (
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(greenColor)).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(defaultTextColor))
)

/**
* Tea Model Functions: Init(), Update() and View()
**/
func (m *InfoPanel) Init() tea.Cmd {
	return nil
}

func (m *InfoPanel) Update(message tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := message.(type) {
	case InfoUpdate:
		m.updateInfo(msg.destinationUrl, msg.port)
	}

	_, cmd = m.panel.Update(message)
	cmds = append(cmds, cmd)

	m.updateInfo(m.destinationUrl, m.port)
	return m, tea.Batch(cmds...)
}

func (m *InfoPanel) View() string {
	return m.panel.View()
}

func (m *InfoPanel) SetContent(content string) {
	m.panel.SetContent(content)
}

func (m *InfoPanel) SetFrameSize(width, height int) {
	m.panel.SetFrameSize(width, height)
}

func UpdateInfo(destionationUrl string, port int) tea.Msg {
	return InfoUpdate{
		destinationUrl: destionationUrl,
		port:           port,
	}
}

func (m *InfoPanel) updateInfo(destinationUrl string, port int) {
	m.destinationUrl = destinationUrl
	m.port = port
	m.panel.SetContent(m.getInfoContent())
}

func (m *InfoPanel) getInfoContent() string {
	return labelStyle.Render("Destionation URL:") + " " + valueStyle.Render(m.destinationUrl) +
		"\n" +
		labelStyle.Render("Port Number:") + " " + valueStyle.Render(strconv.Itoa(m.port))
}

/**
* Constructor
**/
func NewInfoPanel(destinationUrl string, port int) *InfoPanel {
	panel := &InfoPanel{
		panel:          *NewPanel("Info"),
		destinationUrl: destinationUrl,
		port:           port,
	}
	return panel
}
