package docker

// ContainerState represents the observed lifecycle state of a container.
type ContainerState int

// The set of observable container lifecycle states, ordered from not-yet-started
// to terminal.
const (
	StatePending   ContainerState = iota // waiting on a dependency; not yet started
	StateStarting                        // container exists, health check running
	StateHealthy                         // running and health check passing (or no healthcheck)
	StateUnhealthy                       // health check explicitly failing
	StateExited                          // stopped with non-zero exit code
)

func (s ContainerState) String() string {
	switch s {
	case StatePending:
		return "pending"
	case StateStarting:
		return "starting"
	case StateHealthy:
		return "healthy"
	case StateUnhealthy:
		return "unhealthy"
	case StateExited:
		return "exited"
	default:
		return "unknown"
	}
}
