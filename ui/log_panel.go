package ui

import (
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogPanel struct {
	panel Panel
	Logs  []string
}

type LogEvent struct {
	Request *http.Request
}

/**
* tea.Msg
 */
type newLog []byte

func (p *LogPanel) Init() tea.Cmd {
	return nil
}

/**
** Style Configurations
**/

var (
	requestPathStyle = lipgloss.NewStyle().
				Bold(false).
				Foreground(lipgloss.Color("#FAFAFA"))

	httpMethodStyle = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Left).
			Bold(true)

	requestHostStyle = lipgloss.NewStyle().
				Bold(false).
				Foreground(lipgloss.Color("#FFFFFF"))
)

func (p *LogPanel) Update(message tea.Msg) (tea.Msg, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := message.(type) {
	case LogEvent:
		p.LogRequest(msg.Request)
	}
	_, cmd = p.panel.Update(message)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p *LogPanel) View() string {
	return p.panel.View()
}

func (p *LogPanel) SetFrameSize(width, height int) {
	p.panel.SetFrameSize(width, height)
}

/**
* Member Functions
**/
func (p *LogPanel) LogRequest(request *http.Request) {
	p.Logs = append(p.Logs, getRequestLogUi(request))
	vpcontent := ""
	for _, message := range p.Logs {
		vpcontent += message + "\n"
	}
	p.panel.SetContent(vpcontent)
	p.panel.ScrollDown(len(p.Logs))
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
		color = greenColor
	case "POST", "PUT":
		color = yellowColor
	case "DELETE":
		color = redColor
	default:
		color = blueColor
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

func NewLogPanel() *LogPanel {
	return &LogPanel{
		panel: *NewPanel("Logs"),
	}
}
