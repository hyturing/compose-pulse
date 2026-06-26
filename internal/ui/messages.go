package ui

// logsMsg carries fetched log lines for the currently selected service.
type logsMsg struct {
	service string
	lines   []string
	err     error
}

// tickMsg is sent by the spinner animation ticker.
type tickMsg struct{}
