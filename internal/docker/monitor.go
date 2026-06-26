package docker

import (
	"context"
	"time"
)

const pollInterval = 500 * time.Millisecond

// StateMsg is a Bubble Tea message emitted each poll cycle.
type StateMsg struct {
	States map[string]ContainerState
}

// StartCh launches a background goroutine that polls Docker every 500 ms and
// sends a StateMsg on the returned channel. The goroutine exits when ctx is
// cancelled and the channel is then closed.
func (c *Client) StartCh(ctx context.Context, services []string) <-chan StateMsg {
	ch := make(chan StateMsg, 1)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				states, err := c.FetchStates(ctx, services)
				if err != nil {
					continue
				}
				select {
				case ch <- StateMsg{States: states}:
				default: // drop frame if UI is still processing the previous one
				}
			}
		}
	}()
	return ch
}
