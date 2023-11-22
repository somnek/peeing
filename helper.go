package main

import (
	tea "github.com/charmbracelet/bubbletea"
	probing "github.com/prometheus-community/pro-bing"
)

func isPacketRecv(msg *probing.Statistics) bool {
	return msg.PacketsRecv == 1
}

func ping(url string) tea.Cmd {
	return func() tea.Msg {
		// create a new pinger
		pinger, err := probing.NewPinger(url)
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

}

func getStats(pinger *probing.Pinger) (*probing.Statistics, error) {
	err := pinger.Run()
	if err != nil {
		return nil, err
	}
	return pinger.Statistics(), err

}
