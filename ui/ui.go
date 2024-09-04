package ui

import (
	"fmt"
	_ "strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"LitTime/estimator"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Align(lipgloss.Center) // Центрируем заголовок

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Align(lipgloss.Left) // Выравниваем данные по левому краю

	// Стиль для строк результата
	resultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Align(lipgloss.Left)
)

type model struct {
	viewport viewport.Model
	result   *estimator.Result
	ready    bool
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
		// Адаптируем размеры viewport к размеру окна терминала
		headerHeight := 3 // Высота заголовка
		footerHeight := 2 // Высота для строки с информацией о выходе
		verticalMargin := headerHeight + footerHeight

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMargin
		m.ready = true
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "Инициализация..."
	}

	// Собираем вывод на экран
	content := titleStyle.Render("LitTime Results") + "\n\n"
	content += resultStyle.Render(fmt.Sprintf("Reading time: %s", highlightStyle.Render(fmt.Sprintf("%.2f min", m.result.ReadingTime)))) + "\n"
	content += resultStyle.Render(fmt.Sprintf("Words: %s", highlightStyle.Render(fmt.Sprintf("%d", m.result.WordCount)))) + "\n"
	content += resultStyle.Render(fmt.Sprintf("Sentences: %s", highlightStyle.Render(fmt.Sprintf("%d", m.result.SentenceCount)))) + "\n"
	content += resultStyle.Render(fmt.Sprintf("Syllables: %s", highlightStyle.Render(fmt.Sprintf("%d", m.result.SyllableCount)))) + "\n"
	content += resultStyle.Render(fmt.Sprintf("Flesch-Kincaid Index: %s", highlightStyle.Render(fmt.Sprintf("%.2f", m.result.FleschKincaidIndex)))) + "\n"

	// Обновляем контент viewport
	m.viewport.SetContent(content)

	return fmt.Sprintf("%s\n%s", m.viewport.View(), infoStyle.Render("Press q to quit"))
}

func RunUI(result *estimator.Result) error {
	// Создаем модель с начальным состоянием
	vp := viewport.Model{Width: 80, Height: 20} // Стартовые размеры по умолчанию

	// Запускаем программу с динамическим изменением размеров
	p := tea.NewProgram(model{viewport: vp, result: result})
	_, err := p.Run()
	return err
}
