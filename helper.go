package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	probing "github.com/prometheus-community/pro-bing"
)

const (
	OUTPUT_FILE = "output.log"
)

func insertHistory(history []string, t time.Time, rtt time.Duration) []string {
	history = append(history, rtt.String())
	history = history[1:]
	return history
}

// check if the ping request is successful
func isPacketRecv(msg *probing.Statistics) bool {
	return msg.PacketsRecv == 1
}

// ping() sends a ping request to the specified URL and returns a tea.Cmd.
// The ping request measures the round-trip time (RTT) and collects statistics about the network connection.
// url: parameter specifies the target URL to ping.
// - returns a tea.Cmd, which is a command that can be executed by a tea.Program.
// - tea.Cmd function is executed asynchronously and returns a tea.Msg when completed.
// - tea.Msg contains the ping statistics and the duration of the ping request.
func ping(url string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		// create a new pinger
		pinger, err := probing.NewPinger(url)
		if err != nil {
			return errMsg{err}
		}

		// ping options
		pinger.Timeout = TimeLimit
		pinger.Count = 1

		stats, err := getStats(pinger)
		if err != nil {
			return errMsg{err}
		}

		duration := time.Since(start)
		t := timing{
			dur:   duration,
			start: start,
			end:   time.Now(),
		}
		return pingMsg{
			stats:  stats,
			timing: t,
		}
	}
}

// returns the ping statistics
func getStats(pinger *probing.Pinger) (*probing.Statistics, error) {
	err := pinger.Run()
	if err != nil {
		return nil, err
	}
	return pinger.Statistics(), err

}

// this is just a basic input validate
// for url validation, let probing.NewPinger do it
func isValidInput(u string) bool {
	return strings.Contains(u, ".") && !strings.Contains(u, " ")
}

// we have 8 block symbols to represent the RTT
func convertToBlockUnit(dur time.Duration) rune {
	if dur == -1*time.Millisecond {
		return SHADED_BLOCKS[0]
	}

	unitMap := []time.Duration{
		50 * time.Millisecond,
		100 * time.Millisecond,
		150 * time.Millisecond,
		200 * time.Millisecond,
		250 * time.Millisecond,
		300 * time.Millisecond,
		350 * time.Millisecond,
		400 * time.Millisecond,
	}

	for i, v := range unitMap {
		if dur <= v {
			return BARS[i]
		}
	}

	// if the duration is greater than 400ms, we return the last block symbol
	return BARS[7]
}

func save(history []record) error {
	// format content
	var content string
	for _, rec := range history {
		if rec == (record{}) {
			// trim the place holder
			continue
		}
		url := rec.url
		timestamp := rec.timestamp.Format(time.RFC3339)
		rtt := rec.rtt.String()
		content += fmt.Sprintf("%s: [%s]  %s\n", timestamp, url, rtt)
	}

	// to file
	f, err := os.Create(OUTPUT_FILE)
	if err != nil {
		return err
	}

	_, err = io.WriteString(f, content)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}
