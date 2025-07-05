package ui

import (
	"fmt"
	"go-proxy/proxy"
	"net/http"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogPanel struct {
	panel Panel
	Logs  []string
}

type LogEvent struct {
	log proxy.RequestLog
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

	requestTimeStyle = lipgloss.NewStyle().
				Bold(false).
				Foreground(lipgloss.Color("#FAFAFA")).
				AlignHorizontal(lipgloss.Right)
)

func (p *LogPanel) Update(message tea.Msg) (tea.Msg, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := message.(type) {
	case LogEvent:
		p.LogRequest(msg.log)
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
func (p *LogPanel) LogRequest(log proxy.RequestLog) {
	p.Logs = append(p.Logs, getRequestLogUi(log))
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
func getRequestLogUi(log proxy.RequestLog) string {
	return getHttpMethodUi(log.Request) +
		" " +
		getRequestHostUi(log.Request) +
		" " +
		getRequestTimeUi(log.Time) +
		" " +
		getRequestPathUi(log.Request)
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

func getRequestTimeUi(time time.Duration) string {
	milliseconds := time.Abs().Milliseconds()
	durationString := strconv.Itoa(int(milliseconds))
	label := fmt.Sprintf("%vms", durationString)
	return requestTimeStyle.Render(label)
}

func NewLogPanel() *LogPanel {
	return &LogPanel{
		panel: *NewPanel("Logs"),
	}
}
