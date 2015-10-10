package cerebrum

import "github.com/hashicorp/serf/serf"

func (c *cerebrum) HandleMemberJoin(e serf.MemberEvent) {
	for _, m := range e.Members {
		c.logger.Info("member joined", "name", m.Name, "addr", m.Addr, "port", m.Port)
	}
	if c.config.NodeJoined != nil {
		c.config.NodeJoined.HandleMemberJoin(e)
	}
}
func (c *cerebrum) HandleMemberUpdate(e serf.MemberEvent) {
	for _, m := range e.Members {
		c.logger.Info("member updated", "name", m.Name, "addr", m.Addr, "port", m.Port)
	}
	if c.config.NodeUpdated != nil {
		c.config.NodeUpdated.HandleMemberUpdate(e)
	}
}
func (c *cerebrum) HandleMemberLeave(e serf.MemberEvent) {
	for _, m := range e.Members {
		c.logger.Info("member left", "name", m.Name, "addr", m.Addr, "port", m.Port)
	}
	if c.config.NodeLeft != nil {
		c.config.NodeLeft.HandleMemberLeave(e)
	}
}
func (c *cerebrum) HandleMemberFailure(e serf.MemberEvent) {
	for _, m := range e.Members {
		c.logger.Info("member failed", "name", m.Name, "addr", m.Addr, "port", m.Port)
	}
	if c.config.NodeFailed != nil {
		c.config.NodeFailed.HandleMemberFailure(e)
	}
}
func (c *cerebrum) HandleMemberReap(e serf.MemberEvent) {
	for _, m := range e.Members {
		c.logger.Info("member reaped", "name", m.Name, "addr", m.Addr, "port", m.Port)
	}
	if c.config.NodeReaped != nil {
		c.config.NodeReaped.HandleMemberReap(e)
	}
}

func (c *cerebrum) HandleLeaderElection(evt serf.UserEvent) {
	c.setLeader(string(evt.Payload))

	if c.config.LeaderElectionHandler != nil {
		c.config.LeaderElectionHandler.HandleLeaderElection(evt)
	}
}
