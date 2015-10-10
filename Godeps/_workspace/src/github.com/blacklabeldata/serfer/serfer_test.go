package serfer

import (
	"testing"
	"time"

	"github.com/hashicorp/serf/serf"
	"github.com/stretchr/testify/assert"
)

func TestRunSerfer(t *testing.T) {

	// Create event
	evt := &MockEvent{}

	// Create handler
	handler := &MockEventHandler{}
	handler.On("HandleEvent", evt).Return()

	// Create channel and serfer
	ch := make(chan serf.Event, 1)
	serfer := NewSerfer(ch, handler)

	// Start serfer
	serfer.Start()
	// death.Go(func() error {
	// 	return serfer.Run(ctx)
	// })

	// Send events
	select {
	case ch <- evt:
	case <-time.After(time.Second):
		t.Fatal("Event was not sent over channel")
	}
	ch <- evt

	// Verify stopped without error
	assert.Nil(t, serfer.Stop(), "Error should be nil")

	// Validate event was prcoessed
	handler.AssertCalled(t, "HandleEvent", evt)

}
