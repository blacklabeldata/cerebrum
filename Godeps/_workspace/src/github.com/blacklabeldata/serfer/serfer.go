package serfer

import (
	"github.com/hashicorp/serf/serf"
	tomb "gopkg.in/tomb.v2"
)

// Serfer processes Serf.Events and is meant to be ran in a goroutine.
type Serfer interface {

	// Start starts the serfer goroutine.
	Start()

	// Stop stops all event processing and blocks until finished.
	Stop() error
}

// NewSerfer returns a new Serfer implementation that uses the given channel and event handlers.
func NewSerfer(c chan serf.Event, handler EventHandler) Serfer {
	var t tomb.Tomb
	return &serfer{handler, c, t}
}

type serfer struct {
	handler EventHandler
	channel chan serf.Event
	t       tomb.Tomb
}

func (s *serfer) Start() {
	s.t.Go(func() error {
		// Start event processing
		for {
			select {

			// Handle context close
			case <-s.t.Dying():
				return nil

			// Handle serf events
			case evt := <-s.channel:
				s.handler.HandleEvent(evt)
			}
		}
	})
}

func (s *serfer) Stop() error {
	s.t.Kill(nil)
	return s.t.Wait()
}
