package ui

import (
	"context"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hyturing/compose-pulse/internal/docker"
)

const logTailLines = 200

// fetchLogsCmd returns an async Bubble Tea Cmd that fetches logs for service.
func fetchLogsCmd(dc *docker.Client, service string) tea.Cmd {
	return func() tea.Msg {
		raw, err := dc.Logs(context.Background(), service, logTailLines)
		lines := strings.Split(strings.TrimRight(raw, "\n"), "\n")
		return logsMsg{service: service, lines: lines, err: err}
	}
}

// filterLines returns lines matching pattern (regex, with substring fallback).
func filterLines(lines []string, pattern string) []string {
	if pattern == "" {
		return lines
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		// Invalid regex — fall back to plain substring match.
		var out []string
		for _, l := range lines {
			if strings.Contains(l, pattern) {
				out = append(out, l)
			}
		}
		return out
	}
	var out []string
	for _, l := range lines {
		if re.MatchString(l) {
			out = append(out, l)
		}
	}
	return out
}

// renderModal renders the log overlay centered over the blurred tree.
func renderModal(m Model) string {
	modalW := m.width * 80 / 100
	modalH := m.height * 80 / 100
	if modalH < 5 {
		modalH = 5
	}

	lines := filterLines(m.logs, m.filter)

	// Apply scroll offset.
	maxScroll := len(lines) - (modalH - 4)
	if maxScroll < 0 {
		maxScroll = 0
	}
	scroll := m.logScroll
	if scroll > maxScroll {
		scroll = maxScroll
	}
	visible := lines
	if scroll < len(lines) {
		visible = lines[scroll:]
	}
	if len(visible) > modalH-4 {
		visible = visible[:modalH-4]
	}

	content := strings.Join(visible, "\n")

	// Search bar.
	var searchBar string
	if m.searching {
		searchBar = "\n" + styleSearchPrompt.Render("/") + styleSearchInput.Render(m.filter+"█")
	} else if m.filter != "" {
		searchBar = "\n" + styleSearchPrompt.Render("filter: ") + styleSearchInput.Render(m.filter)
	}

	title := lipgloss.NewStyle().Bold(true).Render("logs: " + m.modalSvc)
	body := title + "\n" + strings.Repeat("─", modalW-6) + "\n" + content + searchBar

	modal := styleModal.Width(modalW).Render(body)

	// Center horizontally and vertically.
	hPad := (m.width - modalW) / 2
	vPad := (m.height - modalH) / 2
	pad := strings.Repeat("\n", vPad)
	linePad := strings.Repeat(" ", hPad)

	var sb strings.Builder
	sb.WriteString(pad)
	for _, line := range strings.Split(modal, "\n") {
		sb.WriteString(linePad + line + "\n")
	}
	return sb.String()
}
