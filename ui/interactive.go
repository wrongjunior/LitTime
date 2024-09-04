package ui

import (
	"LitTime/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
)

var (
	appStyle = lipgloss.NewStyle().Margin(1, 2)

	titleStyleIter = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	inputStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
)

type userInputs struct {
	FilePath     string
	ReadingSpeed int
	HasVisuals   bool
	Workers      int
}

type modelIter struct {
	inputs     []textinput.Model
	focusIndex int
	userInputs userInputs
	config     *config.Config
}

func RunInteractive(cfg *config.Config) (userInputs, error) {
	m := initialModel(cfg)

	p := tea.NewProgram(m)
	finishedModel, err := p.Run()
	if err != nil {
		return userInputs{}, err
	}

	finalModel := finishedModel.(modelIter)
	return finalModel.userInputs, nil
}

func initialModel(cfg *config.Config) modelIter {
	inputs := make([]textinput.Model, 4)
	labels := []string{"File Path:", "Reading Speed (wpm):", "Has Visuals (y/n):", "Workers:"}
	defaults := []string{"", strconv.Itoa(cfg.DefaultReadingSpeed), "n", strconv.Itoa(cfg.DefaultWorkers)}

	for i := range inputs {
		t := textinput.New()
		t.Placeholder = labels[i]
		t.SetValue(defaults[i])
		t.CharLimit = 70
		inputs[i] = t
	}

	inputs[0].Focus()

	return modelIter{
		inputs: inputs,
		config: cfg,
	}
}

func (m modelIter) Init() tea.Cmd {
	return textinput.Blink
}

func (m modelIter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Проверяем завершение ввода
			if s == "enter" && m.focusIndex == len(m.inputs) {
				// Проверяем, что путь к файлу не пустой
				if strings.TrimSpace(m.inputs[0].Value()) == "" {
					m.inputs[0].SetValue("Error: file path cannot be empty")
					m.inputs[0].Focus()
					m.focusIndex = 0
					return m, nil
				}

				// Если все поля заполнены корректно, сохраняем данные и выходим
				m.userInputs.FilePath = m.inputs[0].Value()
				m.userInputs.ReadingSpeed, _ = strconv.Atoi(m.inputs[1].Value())
				m.userInputs.HasVisuals = strings.ToLower(m.inputs[2].Value()) == "y"
				m.userInputs.Workers, _ = strconv.Atoi(m.inputs[3].Value())

				return m, tea.Quit
			}

			// Навигация между полями ввода
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

			// Обновляем фокусировку на полях ввода
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
	}

	// Обновляем поля ввода
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *modelIter) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m modelIter) View() string {
	var b strings.Builder

	b.WriteString(titleStyleIter.Render("LitTime Setup"))
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

	b.WriteString(dimStyle.Render("tab: next • enter: submit • q: quit"))

	return appStyle.Render(b.String())
}
