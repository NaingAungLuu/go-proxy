package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TestModel struct {
	infoPanel InfoPanel
	logPanel  LogPanel
}

/**
* Style Configurations
**/
const (
	defaultTextColor = "#FAFAFA"
	greenColor       = "#42f5b6"
	yellowColor      = "#ffbe57"
	redColor         = "#ff5757"
	blueColor        = "#57d8ff"
)

var (
	viewPortStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Margin(1).
			MarginTop(1).
			PaddingLeft(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1).
			MarginLeft(2)

	titleBar = NewTitleBar("Logs")
)

/**
* Tea Model Functions: Init(), Update() and View()
**/
func (m *TestModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m *TestModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.updateFrameSize(msg.Width, msg.Height)
		return m, nil
	}

	_, cmd = m.infoPanel.Update(message)
	cmds = append(cmds, cmd)

	_, cmd = m.logPanel.Update(message)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)

}

func (m *TestModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.infoPanel.View(), m.logPanel.View()) +
		helpStyle.Render("↑/↓: Navigate • q or esc: Quit")
}

func (m *TestModel) updateFrameSize(width, height int) {
	panelHeight := height - helpStyle.GetVerticalFrameSize()
	// Width = half of the screen - (each panel's horizontal margin)
	panelWidth := (width / 2) - (2 * panelStyle.GetHorizontalMargins())

	m.infoPanel.SetFrameSize(panelWidth, panelHeight/2)
	m.logPanel.SetFrameSize(panelWidth, panelHeight)
}

/**
* Constructor
**/
func NewTestModel(destinationUrl string, port int) *TestModel {
	return &TestModel{
		infoPanel: *NewInfoPanel(destinationUrl, port),
		logPanel:  *NewLogPanel(),
	}
}
