package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	probing "github.com/prometheus-community/pro-bing"
)

const url = "google.ca"
const timeLimit = 2 * time.Second

type model struct {
	err    error
	pinger *probing.Pinger
	log    string
}

type pingMsg *probing.Statistics
type errMsg struct{ err error }

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pingMsg:
		if isPacketRecv(msg) {
			m.log = fmt.Sprintf("🐐%v", msg.Rtts)
		} else {
			m.log = "🐇 Failed to ping"
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case " ", "enter":
			m.log = "🌀 pinging..."
			return m, ping
		}
	case errMsg:
		fmt.Printf("⛔ : %v", msg.err)
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	s := "hit space to ping google.ca\n\n"
	s += m.log
	if m.err != nil {
		s += fmt.Sprintf("\n\nerror: %v", m.err)
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, theres been an error: %v", err)
		os.Exit(1)
	}
}
