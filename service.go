package cerebrum

import (
	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
	"golang.org/x/net/context"
)

type Service interface {
	Name() string
	Start(*Context) error
	Stop()
}

type Context struct {
	Context context.Context
	Serf    serf.Serf
	Raft    raft.Raft
}
