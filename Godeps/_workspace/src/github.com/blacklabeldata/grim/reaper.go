package grim

import (
	"sync"

	"golang.org/x/net/context"
)

// GrimReaper is a task runner which will wait until all the tasks
// are complete until comtinuing. Reapers can also kill tasks prematurely if
// they are listening to the given context.
type GrimReaper interface {
	// New creates a sub-reaper which is attached to the parent context. If
	// the parent context is killed, so are the children. However, Wait must
	// still be called for child reapers.
	New() GrimReaper

	// SpawnFunc starts a new goroutine for the given function. The task should
	// return as soon as possible after the context completes.
	SpawnFunc(TaskFunc)

	// Spawn starts a new goroutine for the given Task. The task should
	// return as soon as possible after the context completes.
	Spawn(Task)

	// Kill sends a message to all running tasks to stop. If child reapers
	// have also been created, they will be triggered to stop as well.
	Kill()

	// Wait will block until all tasks have completed. This will NOT block
	// until chlid reapers are finished. Each reaper must call wait
	// independently.
	Wait()
}

// Reaper returns an implementation of the GrimReaper interface.
func Reaper() GrimReaper {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	return &reaper{wg, ctx, cancel}
}

// ReaperWithContext creates a new GrimReaper implementation and uses the given
// context as the parent context.
func ReaperWithContext(c context.Context) GrimReaper {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(c)
	return &reaper{wg, ctx, cancel}
}

// TaskFunc is a killable function which runs in a separate go routine. In order to fulfill the contract the function MUST listen to the context and exit if it fires.
type TaskFunc func(context.Context)

type Task interface {
	Execute(context.Context)
}

// reaper implements the GrimReaper interface.
type reaper struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// SpawnFunc runs a task.
func (r *reaper) SpawnFunc(t TaskFunc) {
	r.wg.Add(1)
	c, _ := context.WithCancel(r.ctx)
	go func(ctx context.Context) {
		defer r.wg.Done()
		t(ctx)
	}(c)
}

// Spawn runs a task.
func (r *reaper) Spawn(t Task) {
	r.wg.Add(1)
	c, _ := context.WithCancel(r.ctx)
	go func(ctx context.Context) {
		defer r.wg.Done()
		t.Execute(ctx)
	}(c)
}

// Wait blocks until all tasks have completed.
func (r *reaper) Wait() {
	r.wg.Wait()
}

// Kill cancels the context and waits for all tasks to exit.
func (r *reaper) Kill() {
	r.cancel()
	r.Wait()
}

// New creates a new reaper with the current context as the parent context of
// the new reaper.
func (r *reaper) New() GrimReaper {
	return ReaperWithContext(r.ctx)
}
