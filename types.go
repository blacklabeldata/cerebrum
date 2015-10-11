package cerebrum

import (
	log "github.com/mgutz/logxi/v1"

	"github.com/blacklabeldata/namedtuple"
)

var (
	nodeStatus namedtuple.TupleType
)

type NodeStatus uint8

const (
	StatusAlive NodeStatus = iota
	StatusFailed
	StatusLeft
	StatusReaped
)

func init() {

	// Node registration type
	nodeStatus = namedtuple.New("cerebrum", "NodeStatus")
	nodeStatus.AddVersion(
		namedtuple.Field{"ID", true, namedtuple.StringField},
		namedtuple.Field{"Name", true, namedtuple.StringField},
		namedtuple.Field{"DataCenter", true, namedtuple.StringField},
		namedtuple.Field{"Status", true, namedtuple.Uint8Field},
		namedtuple.Field{"Addr", true, namedtuple.StringField},
		namedtuple.Field{"Port", true, namedtuple.Int32Field})
	namedtuple.DefaultRegistry.Register(nodeStatus)
}

type NodeStatusUpdater interface {
	Update(*NodeDetails, NodeStatus) error
}

func NewNodeStatusUpdater(a Applier, l log.Logger) NodeStatusUpdater {
	return &nodeStatusUpdater{applier: a, logger: l}
}

type nodeStatusUpdater struct {
	buffer  [512]byte
	applier Applier
	logger  log.Logger
}

func (n *nodeStatusUpdater) Update(details *NodeDetails, status NodeStatus) (err error) {
	builder := namedtuple.NewBuilder(nodeStatus, n.buffer[:])
	if _, err = builder.PutString("ID", details.ID); err != nil {
		n.logger.Warn("Failed to encode NodeStatus", "field", "ID", "err", err)
		return
	}
	if _, err = builder.PutString("Name", details.Name); err != nil {
		n.logger.Warn("Failed to encode NodeStatus", "field", "Name", "err", err)
		return
	}
	if _, err = builder.PutString("DataCenter", details.DataCenter); err != nil {
		n.logger.Warn("Failed to encode NodeStatus", "field", "DataCenter", "err", err)
		return
	}
	if _, err = builder.PutUint8("Status", uint8(status)); err != nil {
		n.logger.Warn("Failed to encode NodeStatus", "field", "Status", "err", err)
		return
	}
	if _, err = builder.PutString("Addr", details.Addr.String()); err != nil {
		n.logger.Warn("Failed to encode NodeStatus", "field", "Addr", "err", err)
		return
	}
	if _, err = builder.PutInt32("Port", int32(details.Port)); err != nil {
		n.logger.Warn("Failed to encode NodeStatus", "field", "Port", "err", err)
		return
	}

	// encode data
	tuple, err := builder.Build()
	if err != nil {
		n.logger.Warn("Failed to build NodeStatus", "err", err)
		return
	}
	return n.applier.Apply(tuple)
}
