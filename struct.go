package main

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	probing "github.com/prometheus-community/pro-bing"
)

type model struct {
	err         error
	inputs      []textinput.Model
	log         string
	isSubmitted bool
	isPinging   bool
	rttList     []time.Duration
	history     []record
	help        string
}

type timing struct {
	dur   time.Duration
	start time.Time
	end   time.Time
}

type pingMsg struct {
	stats *probing.Statistics
	timing
}

type record struct {
	timestamp time.Time
	rtt       time.Duration
	url       string
}

type errMsg struct{ err error }
