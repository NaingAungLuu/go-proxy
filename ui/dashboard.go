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

type newLog []byte

type tickMsg time.Time

type LogEvent struct {
	Request *http.Request
}

var (
	titleStyle = lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).
		// BorderStyle(lipgloss.RoundedBorder()).
		// BorderBottom(false).BorderTop(true).BorderLeft(true).BorderRight(true).
		// BorderForeground(lipgloss.Color("62")).
		Foreground(lipgloss.NoColor{}).
		Background(lipgloss.Color("62")).
		ColorWhitespace(true).
		Padding(0, 1).
		MarginLeft(1).
		MarginTop(1)
)

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
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
		tea.Println("Log Event received!")
		m.Logs = append(m.Logs, LogRequest(msg.Request))
		vpContent := ""
		for _, message := range m.Logs {
			vpContent += message + "\n"
		}
		m.vp.SetContent(vpContent)
		m.vp.ScrollDown(len(m.Logs))
	}

	// Handle keyboard and mouse events in the viewport
	m.vp, cmd = m.vp.Update(message)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateWindowSize(width, height int) {
	_, titleHeight := titleStyle.GetFrameSize()
	m.vp.Width = width
	m.vp.Height = height - (titleHeight) // - titleHeight
	m.vp.Style.Width(width)
	m.vp.Style.Height(height - (titleHeight))
}

func NewModel() *Model {
	vp := viewport.New(78, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		MarginTop(1).
		PaddingLeft(1)
	return &Model{
		vp: vp,
	}
}

func (m Model) View() string {
	return titleStyle.Render("Logs") + m.vp.View()
}

func (m Model) LogRequest(request *http.Request) {
	m.Logs = append(m.Logs, LogRequest(request))
}

func LogRequest(request *http.Request) string {
	return getHttpMethodUi(*request) + " " + getRequestHostUi(*request) + " " + getRequestPathUi(*request)
}

func getHttpMethodUi(request http.Request) string {
	style := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Left).
		Bold(true)

	switch request.Method {
	case "GET":
		style = style.Foreground(lipgloss.Color("#42f5b6"))
	case "POST", "PUT":
		style = style.Foreground(lipgloss.Color("#ffbe57"))
	case "DELETE":
		style = style.Foreground(lipgloss.Color("#ff5757"))
	default:
		style = style.Foreground(lipgloss.Color("#57d8ff"))
	}

	return style.Render(request.Method)
}

func getRequestHostUi(request http.Request) string {
	style := lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color("#FFFFFF"))

	return style.Render(request.Host)

}

func getRequestPathUi(request http.Request) string {
	style := lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color("#FAFAFA"))

	return style.Render(request.URL.String())
}
