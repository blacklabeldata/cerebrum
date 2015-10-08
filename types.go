package cerebrum

import "github.com/blacklabeldata/namedtuple"

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
	nodeStatus.AddVersion([]namedtuple.Field{
		namedtuple.Field{"ID", true, namedtuple.StringField},
		namedtuple.Field{"Name", true, namedtuple.StringField},
		namedtuple.Field{"DataCenter", true, namedtuple.StringField},
		namedtuple.Field{"Status", true, namedtuple.Uint8Field},
		namedtuple.Field{"Addr", true, namedtuple.StringField},
		namedtuple.Field{"Port", true, namedtuple.Int32Field},
	}...)
	namedtuple.DefaultRegistry.Register(nodeStatus)
}

func (c *cerebrum) updateNodeStatus(details *NodeDetails, status NodeStatus) (err error) {
	buffer := make([]byte, 512)
	builder := namedtuple.NewBuilder(nodeStatus, buffer)
	if _, err = builder.PutString("ID", details.ID); err != nil {
		c.logger.Warn("Failed to encode NodeStatus", "field", "ID", "err", err)
		return
	}
	if _, err = builder.PutString("Name", details.Name); err != nil {
		c.logger.Warn("Failed to encode NodeStatus", "field", "Name", "err", err)
		return
	}
	if _, err = builder.PutString("DataCenter", details.DataCenter); err != nil {
		c.logger.Warn("Failed to encode NodeStatus", "field", "DataCenter", "err", err)
		return
	}
	if _, err = builder.PutUint8("Status", uint8(status)); err != nil {
		c.logger.Warn("Failed to encode NodeStatus", "field", "Status", "err", err)
		return
	}
	if _, err = builder.PutString("Addr", details.Addr.String()); err != nil {
		c.logger.Warn("Failed to encode NodeStatus", "field", "Addr", "err", err)
		return
	}
	if _, err = builder.PutInt32("Port", int32(details.Port)); err != nil {
		c.logger.Warn("Failed to encode NodeStatus", "field", "Port", "err", err)
		return
	}

	// encode data
	tuple, err := builder.Build()
	if err != nil {
		c.logger.Warn("Failed to build NodeStatus", "err", err)
		return
	}
	return c.apply(tuple)
}
