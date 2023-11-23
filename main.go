package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	probing "github.com/prometheus-community/pro-bing"
)

const timeLimit = 1 * time.Second

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("120"))
)

type model struct {
	err         error
	pinger      *probing.Pinger
	inputs      []textinput.Model
	log         string
	isSubmmited bool
}

type pingMsg *probing.Statistics
type errMsg struct{ err error }

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 1),
	}

	var t textinput.Model
	t = textinput.New()
	t.CharLimit = 32
	t.Placeholder = "Enter a URL"
	t.Focus()

	m.inputs[0] = t
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case pingMsg:
		if isPacketRecv(msg) {
			m.log = fmt.Sprintf("🐐%v", msg.Rtts)
		} else {
			m.log = "🐇 Failed"
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			if !m.inputs[0].Focused() {
				// quit if not focused
				return m, tea.Quit
			} else {
				// reset input if focused
				if m.inputs[0].Value() != "" {
					m.inputs[0].Reset()
					m.inputs[0].Focus()
					return m, nil
				}
			}
			return m, tea.Quit

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			// allowing 'q' to be used as input
			if !m.inputs[0].Focused() {
				return m, tea.Quit
			}

		case "ctrl+r":
			// reset input
			m.inputs[0].Reset()
			m.inputs[0].Focus()
			return m, nil

		case " ", "enter":
			url := m.inputs[0].Value()

			// validate input
			if !isValidInput(url) {
				m.inputs[0].Reset()
				m.err = fmt.Errorf("invalid input")
				return m, nil
			}
			// unfocused & submit
			m.log = "🌀 pinging..."
			m.err = nil
			m.inputs[0].Blur()
			m.isSubmmited = true

			return m, ping(url)
		}

	case errMsg:
		fmt.Printf("⛔ : %v", msg.err)
		return m, tea.Quit
	}

	// handle character input & blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	// only text inputs with Focus() set will respond
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(helpStyle.Render("hit space to ping 📡"))
	b.WriteRune('\n')
	b.WriteString(m.log)

	if m.isSubmmited == false {
		for i := range m.inputs {
			b.WriteString(m.inputs[i].View())
		}
	}

	b.WriteRune('\n')

	if m.err != nil {
		b.WriteString(fmt.Sprintf("🚫 error: %v", m.err))
	}
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, theres been an error: %v", err)
		os.Exit(1)
	}
}
