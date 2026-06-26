package ui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hyturing/compose-pulse/internal/dag"
	"github.com/hyturing/compose-pulse/internal/docker"
)

// Model is the root Bubble Tea model.
type Model struct {
	graph   *dag.Graph
	docker  *docker.Client
	stateCh <-chan docker.StateMsg
	cancel  context.CancelFunc

	// Navigation: index into graph.Visual (depth-first render order).
	cursor int

	// Log modal state.
	modalOpen bool
	modalSvc  string
	logs      []string
	logScroll int
	searching bool
	filter    string
	spinFrame int

	// Terminal dimensions.
	width, height int
}

// New creates the initial model and starts the Docker polling goroutine.
func New(g *dag.Graph, dc *docker.Client) Model {
	ctx, cancel := context.WithCancel(context.Background())
	services := make([]string, 0, len(g.ByName))
	for name := range g.ByName {
		services = append(services, name)
	}
	return Model{
		graph:   g,
		docker:  dc,
		stateCh: dc.StartCh(ctx, services),
		cancel:  cancel,
	}
}

// waitForState returns a Cmd that blocks on the next StateMsg from the monitor channel.
// Each call to Update that receives a StateMsg must re-dispatch this Cmd to keep
// the listen loop alive.
func waitForState(ch <-chan docker.StateMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		return msg
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Init starts the spinner ticker and begins listening for state updates.
func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), waitForState(m.stateCh))
}

// Update routes incoming messages, mutating model state and returning any
// follow-up commands.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.spinFrame = (m.spinFrame + 1) % len(spinnerFrames)
		return m, tickCmd()

	case docker.StateMsg:
		for name, state := range msg.States {
			if node, ok := m.graph.ByName[name]; ok {
				node.State = state
			}
		}
		return m, waitForState(m.stateCh) // re-arm listener

	case logsMsg:
		if msg.err == nil {
			m.logs = msg.lines
		}

	case tea.KeyMsg:
		if m.searching {
			return m.updateSearch(msg)
		}
		if m.modalOpen {
			return m.updateModal(msg)
		}
		return m.updateTree(msg)
	}

	return m, nil
}

func (m Model) updateTree(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		m.cancel()
		return m, tea.Quit

	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.graph.Visual)-1 {
			m.cursor++
		}

	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case key.Matches(msg, keys.Enter):
		if len(m.graph.Visual) > 0 {
			svc := m.graph.Visual[m.cursor].Name
			m.modalOpen = true
			m.modalSvc = svc
			m.logs = nil
			m.logScroll = 0
			m.filter = ""
			return m, fetchLogsCmd(m.docker, svc)
		}
	}
	return m, nil
}

func (m Model) updateModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Esc), key.Matches(msg, keys.Quit):
		m.modalOpen = false
		m.logs = nil
	case key.Matches(msg, keys.Search):
		m.searching = true
	case key.Matches(msg, keys.Down):
		m.logScroll++
	case key.Matches(msg, keys.Up):
		if m.logScroll > 0 {
			m.logScroll--
		}
	}
	return m, nil
}

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.searching = false
		m.filter = ""
	case tea.KeyEnter:
		m.searching = false
	case tea.KeyBackspace:
		if len(m.filter) > 0 {
			m.filter = m.filter[:len(m.filter)-1]
		}
	default:
		if msg.Type == tea.KeyRunes {
			m.filter += string(msg.Runes)
		}
	}
	return m, nil
}

// View renders the current model state to a string for display.
func (m Model) View() string {
	if m.modalOpen {
		return renderModal(m)
	}
	return renderTree(m.graph, m.cursor, m.spinFrame, m.width)
}
