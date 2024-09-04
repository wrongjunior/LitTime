package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"LitTime/estimator"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Margin(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	inputStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	focusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	blurredStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	dimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	highlightStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
)

type model struct {
	inputs          []textinput.Model
	focusIndex      int
	result          *estimator.Result
	processingTime  time.Duration
	processingError error
}

func initialModel() model {
	inputs := make([]textinput.Model, 4)
	labels := []string{"File Path:", "Reading Speed (wpm):", "Has Visuals (t/f):", "Workers:"}
	defaults := []string{"", "180", "false", fmt.Sprintf("%d", runtime.NumCPU())}

	for i := range inputs {
		t := textinput.New()
		t.Placeholder = labels[i]
		t.SetValue(defaults[i])
		t.CharLimit = 50
		inputs[i] = t
	}

	inputs[0].Focus()

	return model{inputs: inputs}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, processInput(m)
			}
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}
			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		}
	case processResult:
		m.result = &msg.result
		m.processingTime = msg.processingTime
		m.processingError = nil
		return m, nil
	case processError:
		m.processingError = msg.err
		return m, nil
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func processInput(m model) tea.Cmd {
	return func() tea.Msg {
		filePath := strings.TrimSpace(m.inputs[0].Value())
		readingSpeed, _ := strconv.ParseFloat(strings.TrimSpace(m.inputs[1].Value()), 64)
		hasVisuals, _ := strconv.ParseBool(strings.TrimSpace(m.inputs[2].Value()))
		workerCount, _ := strconv.Atoi(strings.TrimSpace(m.inputs[3].Value()))

		if filePath == "" {
			return processError{fmt.Errorf("file path cannot be empty")}
		}

		start := time.Now()
		text, err := estimator.ReadTextFromFile(filePath)
		if err != nil {
			return processError{err}
		}

		result, err := estimator.EstimateReadingTimeParallel(text, readingSpeed, hasVisuals, workerCount)
		if err != nil {
			return processError{err}
		}

		duration := time.Since(start)

		return processResult{result, duration}
	}
}

type processResult struct {
	result         estimator.Result
	processingTime time.Duration
}

type processError struct {
	err error
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("LitTime"))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(inputStyle.Render(m.inputs[i].View()))
		b.WriteString("\n")
	}

	button := blurredStyle.Render("[ Submit ]")
	if m.focusIndex == len(m.inputs) {
		button = focusedStyle.Render("[ Submit ]")
	}
	b.WriteString(button + "\n\n")

	if m.processingError != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n\n", m.processingError))
	}

	if m.result != nil {
		b.WriteString(dimStyle.Render("Results:") + "\n")
		b.WriteString(fmt.Sprintf("Reading time: %s\n", highlightStyle.Render(fmt.Sprintf("%.2f min", m.result.ReadingTime))))
		b.WriteString(fmt.Sprintf("Words: %s\n", highlightStyle.Render(fmt.Sprintf("%d", m.result.WordCount))))
		b.WriteString(fmt.Sprintf("Sentences: %s\n", highlightStyle.Render(fmt.Sprintf("%d", m.result.SentenceCount))))
		b.WriteString(fmt.Sprintf("Syllables: %s\n", highlightStyle.Render(fmt.Sprintf("%d", m.result.SyllableCount))))
		b.WriteString(fmt.Sprintf("Flesch-Kincaid: %s\n", highlightStyle.Render(fmt.Sprintf("%.2f", m.result.FleschKincaidIndex))))
		b.WriteString(fmt.Sprintf("Processing: %s\n", dimStyle.Render(m.processingTime.String())))
	}

	b.WriteString("\n" + dimStyle.Render("tab: next • enter: submit • q: quit"))

	return appStyle.Render(b.String())
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
