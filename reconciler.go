package cerebrum

import "github.com/hashicorp/serf/serf"

const (
	// StatusReap is used to update the status of a node if we
	// are handling a EventMemberReap
	StatusReap = serf.MemberStatus(-1)
)

// Reconciler dispatches membership changes to Raft.
type Reconciler struct {
	ReconcileCh chan serf.Member
	IsLeader    func() bool
}

// Reconcile is used to reconcile Serf events with the strongly
// consistent store if we are the current leader
func (r *Reconciler) Reconcile(m serf.Member) {
	if r.IsLeader() {
		select {
		case r.ReconcileCh <- m:
		default:
		}
	}
}
