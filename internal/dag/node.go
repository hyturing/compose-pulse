package dag

import "github.com/hyturing/compose-pulse/internal/docker"

// Node represents a single service in the dependency graph.
type Node struct {
	Name         string
	Deps         []string // all dependency names (sorted, for display)
	Children     []*Node  // all services that depend on this node (full edge set, used by Kahn's)
	TreeChildren []*Node  // dependants for which this node is the primary parent (used by renderer)
	State        docker.ContainerState
	Level        int // render depth: root = 0
}
