package cerebrum

import (
	"fmt"

	"github.com/hashicorp/serf/serf"
	log "github.com/mgutz/logxi/v1"
)

// mergeDelegate is used to handle a cluster merge on the gossip
// ring. We check that the peers are known servers and abort the merge
// otherwise.
type mergeDelegate struct {
	logger log.Logger
}

func (md *mergeDelegate) NotifyMerge(members []*serf.Member) error {
	for _, m := range members {
		_, err := GetNodeDetails(*m)
		if err != nil {
			md.logger.Warn("unknown server attempted to join network", "name", m.Name)
			return fmt.Errorf("member '%s' is not a server", m.Name)
		}
	}
	return nil
}
