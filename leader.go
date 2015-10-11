package cerebrum

import (
	"net"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

// const (
// 	// SerfCheckID           = "serfHealth"
// 	// SerfCheckName         = "Serf Health Status"
// 	// SerfCheckAliveOutput  = "Agent alive and reachable"
// 	// SerfCheckFailedOutput = "Agent not live or unreachable"
// 	KappaServiceID   = "kappa"
// 	KappaServiceName = "kappa"
// 	LeaderEventName  = "kappa:new-leader"
// )

func (c *cerebrum) IsLeader() bool {
	return c.raft.State() == raft.Leader
}

func (c *cerebrum) setLeader(name string) {
	c.leader = name
}

// monitorLeadership is used to monitor if we acquire or lose our role
// as the leader in the Raft cluster. There is some work the leader is
// expected to do, so we must react to changes
func (c *cerebrum) monitorLeadership() {
	leaderCh := c.raft.LeaderCh()
	var stopCh chan struct{}
	for {
		select {
		case isLeader := <-leaderCh:
			if isLeader {
				stopCh = make(chan struct{})
				go c.leaderLoop(stopCh)
				c.logger.Info("cluster leadership acquired")
			} else if stopCh != nil {
				close(stopCh)
				stopCh = nil
				c.logger.Info("cluster leadership lost")
			}
		case <-c.context.Done():
			return
		}
	}
}

// leaderLoop runs as long as we are the leader to run various
// maintenance activities
func (c *cerebrum) leaderLoop(stopCh chan struct{}) {
	// Ensure we revoke leadership on stepdown
	defer c.revokeLeadership()

	// Fire a user event indicating a new leader
	payload := []byte(c.config.NodeName)
	if err := c.serf.UserEvent(CerebrumEventPrefix+":new-leader", payload, false); err != nil {
		c.logger.Warn("failed to broadcast new leader event: %v", err)
	}

	// Reconcile channel is only used once initial reconcile
	// has succeeded
	var reconcileCh chan serf.Member
	establishedLeader := false

RECONCILE:
	// Setup a reconciliation timer
	reconcileCh = nil
	interval := time.After(c.config.ReconcileInterval)

	// Apply a raft barrier to ensure our FSM is caught up
	barrier := c.raft.Barrier(0)
	if err := barrier.Error(); err != nil {
		c.logger.Error("failed to wait for barrier", err)
		goto WAIT
	}

	// Check if we need to handle initial leadership actions
	if !establishedLeader {
		if err := c.establishLeadership(); err != nil {
			c.logger.Error("failed to establish leadership", err)
			goto WAIT
		}
		establishedLeader = true
	}

	// Initial reconcile worked, now we can process the channel
	// updates
	reconcileCh = c.reconcileCh

WAIT:
	// Periodically reconcile as long as we are the leader,
	// or when Serf events arrive
	for {
		select {
		case <-stopCh:
			return
		case <-c.context.Done():
			return
		case <-interval:
			goto RECONCILE
		case member := <-reconcileCh:
			c.reconcileMember(member)
		}
	}
}

// establishLeadership is invoked once we become leader and are able
// to invoke an initial barrier. The barrier is used to ensure any
// previously inflight transactions have been committed and that our
// state is up-to-date.
func (c *cerebrum) establishLeadership() error {
	if c.config.EstablishLeadership != nil {
		return c.config.EstablishLeadership()
	}
	return nil
}

// revokeLeadership is invoked once we step down as leader.
// This is used to cleanup any state that may be specific to a leader.
func (c *cerebrum) revokeLeadership() error {
	if c.config.RevokeLeadership != nil {
		return c.config.RevokeLeadership()
	}
	return nil
}

// reconcileMember is used to do an async reconcile of a single
// serf member
func (c *cerebrum) reconcileMember(member serf.Member) (err error) {

	// Check if this is a member we should handle
	var node *NodeDetails
	if node, err = GetNodeDetails(member); err != nil {
		return err
	}

	// Validate node is in the same data center
	if node.DataCenter != c.config.DataCenter {
		c.logger.Warn("skipping reconcile of node", "node", member)
		return nil
	}

	switch member.Status {
	case serf.StatusAlive:
		err = c.handleAliveMember(member, node)
	case serf.StatusFailed:
		err = c.handleFailedMember(member, node)
	case serf.StatusLeft:
		err = c.handleLeftMember(member, node)
	case StatusReap:
		err = c.handleReapMember(member, node)
	}
	if err != nil {
		c.logger.Error("failed to reconcile member", "member", member, "err", err)
		return err
	}
	return nil
}

// handleAliveMember is used to ensure the node
// is registered, with a passing health check.
func (c *cerebrum) handleAliveMember(member serf.Member, details *NodeDetails) error {

	// Attempt to join the consul server
	if err := c.joinConsulServer(member, details); err != nil {
		return err
	}

	c.logger.Info("member joined, marking health alive", "member", member.Name)
	return c.updateNodeStatus(details, StatusAlive)
}

// handleFailedMember is used to mark the node's status
// as being critical, along with all checks as unknown.
func (c *cerebrum) handleFailedMember(member serf.Member, details *NodeDetails) error {
	c.logger.Info("member failed, marking health critical", "member", member.Name)
	return c.updateNodeStatus(details, StatusFailed)
}

// handleLeftMember is used to handle members that gracefully
// left. They are deregistered if necessary.
func (c *cerebrum) handleLeftMember(member serf.Member, details *NodeDetails) error {
	return c.handleDeregisterMember(StatusLeft, member, details)
}

// handleReapMember is used to handle members that have been
// reaped after a prolonged failure. They are deregistered.
func (c *cerebrum) handleReapMember(member serf.Member, details *NodeDetails) error {
	return c.handleDeregisterMember(StatusReaped, member, details)
}

// handleDeregisterMember is used to deregister a member of a given reason
func (c *cerebrum) handleDeregisterMember(reason NodeStatus, member serf.Member, details *NodeDetails) error {
	// Do not deregister ourself. This can only happen if the current leader
	// is leaving. Instead, we should allow a follower to take-over and
	// deregister us later.
	if member.Name == c.config.NodeName {
		c.logger.Warn("deregistering self should be done by follower", "node", c.config.NodeName)
		return nil
	}

	// Remove from Raft peers if this was a server
	if err := c.removeConsulServer(member, details.Port); err != nil {
		return err
	}

	// Deregister the node
	c.logger.Info("deregistering member", "name", member.Name, "reason", reason)
	return c.updateNodeStatus(details, reason)
}

// joinConsulServer is used to try to join another consul server
func (c *cerebrum) joinConsulServer(m serf.Member, details *NodeDetails) error {
	// Do not join ourself
	if m.Name == c.config.NodeName {
		return nil
	}

	// Check for possibility of multiple bootstrap nodes
	if details.Bootstrap {
		members := c.serf.Members()
		for _, member := range members {
			det, err := GetNodeDetails(member)
			if err != nil && member.Name != m.Name && det.Bootstrap {
				c.logger.Error("Two nodes are both in bootstrap mode. Only one"+
					" node should be in bootstrap mode, not adding Raft peer.",
					"node-1", m.Name, "node-2", member.Name)
				return nil
			}
		}
	}

	// Attempt to add as a peer
	var addr net.Addr = &net.TCPAddr{IP: m.Addr, Port: details.Port}
	future := c.raft.AddPeer(addr.String())
	if err := future.Error(); err != nil && err != raft.ErrKnownPeer {
		c.logger.Error("failed to add raft peer", "err", err)
		return err
	}
	return nil
}

// removeConsulServer is used to try to remove a consul server that has left
func (c *cerebrum) removeConsulServer(m serf.Member, port int) error {
	// Attempt to remove as peer
	peer := &net.TCPAddr{IP: m.Addr, Port: port}
	future := c.raft.RemovePeer(peer.String())
	if err := future.Error(); err != nil && err != raft.ErrUnknownPeer {
		c.logger.Error("failed to remove raft peer",
			"peer", peer, "err", err)
		return err
	} else if err == nil {
		c.logger.Info("removed server as peer", "name", m.Name)
	}
	return nil
}
