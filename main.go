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

const (
	PingInterval = 500 * time.Millisecond
	TimeLimit    = 1 * time.Second
	MaxCol       = 25
	FailedRttVal = -1 * time.Millisecond
)

var (
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("120"))
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("120")).Bold(true)
)

type model struct {
	err         error
	pinger      *probing.Pinger
	inputs      []textinput.Model
	log         string
	isSubmitted bool
	isPinging   bool
	rttList     []time.Duration
}

// type pingMsg *probing.Statistics
type pingMsg struct {
	stats *probing.Statistics
	dur   time.Duration
}
type errMsg struct{ err error }

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 1),
	}

	var t textinput.Model
	t = textinput.New()
	t.CharLimit = 32
	t.Placeholder = "Enter a URL to ping..."
	t.Focus()

	m.inputs[0] = t
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// receive result from ping()
	case pingMsg:
		dur := msg.dur
		url := m.inputs[0].Value()
		stats := msg.stats
		rtt := stats.Rtts[0]

		if isPacketRecv(stats) {
			// display Rtts
			m.log = fmt.Sprintf("🐐%v", rtt)
			m.rttList = append(m.rttList, rtt)
		} else {
			m.log = "🐇 Failed, retrying..."
			m.rttList = append(m.rttList, FailedRttVal)
		}

		// should wait at least 1 second before ping again
		spareTime := PingInterval - dur
		time.Sleep(spareTime)

		m.log += "  🌀 pinging..."
		return m, ping(url)

	// handle shortcut keys (not character input)
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
			if m.isPinging {
				return m, nil
			}

			url := m.inputs[0].Value()

			// validate input
			if !isValidInput(url) {
				m.inputs[0].Reset()
				m.err = fmt.Errorf("invalid input")
				return m, nil
			}
			// unfocused & submit
			m.log += "🌀 pinging..."
			m.err = nil
			m.inputs[0].Blur()
			m.isSubmitted = true
			m.isPinging = true

			return m, ping(url)
		}

	case errMsg:
		m.err = msg.err
		m.inputs[0].Reset()
		m.inputs[0].Focus()
		m.isSubmitted = false
		m.log = ""
		m.isPinging = false
		return m, nil
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

	b.WriteString(titleStyle.Render("Peeing! 📡"))
	b.WriteRune('\n')
	b.WriteString(m.log)

	if !m.isSubmitted {
		for i := range m.inputs {
			b.WriteString(m.inputs[i].View())
		}
	}

	b.WriteRune('\n')

	if m.err != nil {
		b.WriteString(fmt.Sprintf("🚫 error: %v", m.err))
	}

	// bar chart
	// slide the bar chart if exceed maxCol
	relevantRtts := m.rttList
	if len(m.rttList) > MaxCol {
		relevantRtts = m.rttList[len(m.rttList)-MaxCol:]
	}

	for _, rtt := range relevantRtts {
		block := convertToBlockUnit(rtt)
		b.WriteRune(block)
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
