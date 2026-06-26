package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/hyturing/compose-pulse/internal/dag"
	"github.com/hyturing/compose-pulse/internal/docker"
)

// renderTree renders the dependency tree in pstree style.
// cursor is an index into g.Visual, which exactly matches visual row order.
func renderTree(g *dag.Graph, cursor, spinFrame, width int) string {
	if len(g.Visual) == 0 {
		return "No services found."
	}

	// Map each node to its position in the visual (depth-first) order.
	// This is what the cursor indexes into.
	visualIdx := make(map[*dag.Node]int, len(g.Visual))
	for i, n := range g.Visual {
		visualIdx[n] = i
	}

	var sb strings.Builder
	for i, root := range g.Roots {
		isLast := i == len(g.Roots)-1
		renderNode(&sb, root, "", isLast, true, cursor, spinFrame, visualIdx, width)
	}

	sb.WriteString("\n")
	sb.WriteString(styleStatusBar.Render("↑/↓ navigate   enter open logs   q quit"))
	return sb.String()
}

func renderNode(
	sb *strings.Builder,
	n *dag.Node,
	prefix string,
	isLast bool,
	isRoot bool,
	cursor, spinFrame int,
	visualIdx map[*dag.Node]int,
	width int,
) {
	var linePrefix string
	var childPrefix string

	switch {
	case isRoot:
		linePrefix = "  "
		childPrefix = "  "
	case isLast:
		linePrefix = prefix + "└─ "
		childPrefix = prefix + "   "
	default:
		linePrefix = prefix + "├─ "
		childPrefix = prefix + "│  "
	}

	indicator := stateIndicator(n.State, spinFrame)
	name := styleName.Render(n.Name)
	stateTxt := styleDim.Render("(" + n.State.String() + ")")

	// Show extra deps (beyond the primary/first one) inline.
	var extraDeps string
	if len(n.Deps) > 1 {
		extraDeps = styleDim.Render(" also←" + strings.Join(n.Deps[1:], ","))
	}

	line := fmt.Sprintf("%s%s %s %s%s", linePrefix, indicator, name, stateTxt, extraDeps)

	if visualIdx[n] == cursor {
		line = styleSelected.Width(width).Render(line)
	}

	sb.WriteString(line + "\n")

	for i, child := range n.TreeChildren {
		renderNode(sb, child, childPrefix, i == len(n.TreeChildren)-1, false, cursor, spinFrame, visualIdx, width)
	}
}

// stateIndicator returns a colored symbol for the given state.
func stateIndicator(s docker.ContainerState, frame int) string {
	switch s {
	case docker.StateHealthy:
		return lipgloss.NewStyle().Foreground(colorHealthy).Render("●")
	case docker.StateStarting:
		spinner := spinnerFrames[frame%len(spinnerFrames)]
		return lipgloss.NewStyle().Foreground(colorStarting).Render(spinner)
	case docker.StateUnhealthy:
		return lipgloss.NewStyle().Foreground(colorUnhealthy).Render("●")
	case docker.StatePending:
		return lipgloss.NewStyle().Foreground(colorPending).Render("○")
	case docker.StateExited:
		return lipgloss.NewStyle().Foreground(colorUnhealthy).Render("✕")
	default:
		return "?"
	}
}
