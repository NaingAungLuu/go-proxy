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

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {

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
		return m, nil
	}
	return m, nil
}

func (m *Model) updateWindowSize(width, height int) {
	m.vp.Width = width
	m.vp.Height = height
	m.vp.Style.Width(width)
	m.vp.Style.Height(height)
	// m.vp.Style = lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.RoundedBorder()).
	// 	BorderForeground(lipgloss.Color("62")).
	// 	PaddingRig
}

func NewModel() *Model {
	vp := viewport.New(78, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return &Model{
		vp: vp,
	}
}

func (m Model) View() string {
	vpContent := ""
	for _, message := range m.Logs {
		vpContent += message + "\n"
	}
	m.vp.SetContent(vpContent)
	return m.vp.View()
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
