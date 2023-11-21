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

func ping() tea.Msg {
	// create a new pinger
	pinger, err := probing.NewPinger("google.ca")
	if err != nil {
		panic(err)
	}

	pinger.Count = 1
	pinger.Timeout = timeLimit

	stats, err := getStats(pinger)
	if err != nil {
		return errMsg{err}
	}
	return pingMsg(stats)

}

func getStats(pinger *probing.Pinger) (*probing.Statistics, error) {
	err := pinger.Run()
	if err != nil {
		return nil, err
	}
	return pinger.Statistics(), err

}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pingMsg:
		m.log = fmt.Sprintf("🐐%v - %d", msg.Rtts, len(msg.Rtts))
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case " ", "enter":
			m.log += "pinging..."
			return m, ping
		}
	case errMsg:
		m.err = msg.err
		m.log = "⛔"
		return m, ping
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
