package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"LitTime/estimator"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))
)

type model struct {
	viewport viewport.Model
	result   *estimator.Result
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))) {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 1
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s", m.viewport.View(), infoStyle.Render("Press q to quit"))
}

func RunUI(result *estimator.Result) error {
	content := strings.Builder{}
	content.WriteString(titleStyle.Render("LitTime Results\n\n"))
	content.WriteString(fmt.Sprintf("Reading time: %s\n", highlightStyle.Render(fmt.Sprintf("%.2f min", result.ReadingTime))))
	content.WriteString(fmt.Sprintf("Words: %s\n", highlightStyle.Render(fmt.Sprintf("%d", result.WordCount))))
	content.WriteString(fmt.Sprintf("Sentences: %s\n", highlightStyle.Render(fmt.Sprintf("%d", result.SentenceCount))))
	content.WriteString(fmt.Sprintf("Syllables: %s\n", highlightStyle.Render(fmt.Sprintf("%d", result.SyllableCount))))
	content.WriteString(fmt.Sprintf("Flesch-Kincaid Index: %s\n", highlightStyle.Render(fmt.Sprintf("%.2f", result.FleschKincaidIndex))))

	vp := viewport.New(80, 20)
	vp.SetContent(content.String())

	p := tea.NewProgram(model{viewport: vp, result: result})
	_, err := p.Run()
	return err
}
