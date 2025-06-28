package ui

import (
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	vp     viewport.Model
	Logs   []string
	marker int
}

type LogEvent struct {
	Request *http.Request
}

/**
* tea.Msg
 */
type newLog []byte

type tickMsg time.Time

/**
* Style Configurations
**/
const (
	defaultTextColor = "#FAFAFA"
)

var (
	viewPortStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Margin(1).
			MarginTop(1).
			PaddingLeft(1)

	requestPathStyle = lipgloss.NewStyle().
				Bold(false).
				Foreground(lipgloss.Color("#FAFAFA"))

	httpMethodStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			Bold(true)

	requestHostStyle = lipgloss.NewStyle().
				Bold(false).
				Foreground(lipgloss.Color("#FFFFFF"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1).
			MarginLeft(2)

	titleBar = NewTitleBar("Logs")
)

/**
* Tea Model Functions: Init(), Update() and View()
**/
func (m *Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
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
		m.updateWindowSize(msg.Width, msg.Height)
		return m, nil

	case LogEvent:
		m.LogRequest(msg.Request)
	}

	// Handle keyboard and mouse events in the viewport
	m.vp, cmd = m.vp.Update(message)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return titleBar.Render("Logs") +
		m.vp.View() +
		helpStyle.Render("↑/↓: Navigate • q or esc: Quit")
}

/**
* Member Functions
**/
func (m *Model) LogRequest(request *http.Request) {
	m.Logs = append(m.Logs, getRequestLogUi(request))
	vpcontent := ""
	for _, message := range m.Logs {
		vpcontent += message + "\n"
	}
	m.vp.SetContent(vpcontent)
	m.vp.ScrollDown(len(m.Logs))
}

func (m *Model) updateWindowSize(width, height int) {
	_, titleHeight := titleBar.GetFrameSize()
	helperHeight := helpStyle.GetVerticalFrameSize()
	m.vp.Width = width
	m.vp.Height = height - (titleHeight + helperHeight)
	m.vp.Style.Width(width)
	m.vp.Style.Height(height - (titleHeight + helperHeight))
}

/**
* UI Utilities
**/
func getRequestLogUi(request *http.Request) string {
	return getHttpMethodUi(*request) + " " + getRequestHostUi(*request) + " " + getRequestPathUi(*request)
}

func getHttpMethodUi(request http.Request) string {
	color := defaultTextColor
	switch request.Method {
	case "GET":
		color = "#42f5b6"
	case "POST", "PUT":
		color = "#ffbe57"
	case "DELETE":
		color = "#ff5757"
	default:
		color = "#57d8ff"
	}

	style := httpMethodStyle.Foreground(lipgloss.Color(color))
	return style.Render(request.Method)
}

func getRequestHostUi(request http.Request) string {
	return requestHostStyle.Render(request.Host)

}

func getRequestPathUi(request http.Request) string {
	return requestPathStyle.Render(request.URL.String())
}

/**
* Constructor
**/
func NewModel() *Model {
	vp := viewport.New(78, 20)
	vp.Style = viewPortStyle
	return &Model{
		vp: vp,
	}
}
