package cerebrum

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/serf/serf"
)

// GetFullEventName computes the full name of a event
func GetFullEventName(name string) string {
	return CerebrumEventPrefix + name
}

// IsCerebrumEvent checks if a serf event is a cerebrum event
func IsCerebrumEvent(name string) bool {
	return strings.HasPrefix(name, CerebrumEventPrefix)
}

// GetRawEventName is used to get the raw event name
func GetRawEventName(name string) string {
	return strings.TrimPrefix(name, CerebrumEventPrefix)
}

// ValidateNode determines whether a node is a known server and returns
//  its data center and role.
func ValidateNode(member serf.Member) (ok bool, role, dc string) {
	if _, ok = member.Tags["id"]; !ok {
		return false, "", ""
	}

	// Get role name
	if role, ok = member.Tags["role"]; !ok {
		return false, "", ""
	} else if role != CerebrumRole {
		return false, "", ""
	}

	// Get datacenter name
	if dc, ok = member.Tags["dc"]; ok {
		return true, role, dc
	}
	return false, "", ""
}

// GetNodeDetails should validate all the Serf tags for the given member and return
// NodeDetails or any error that occurred.
func GetNodeDetails(m serf.Member) (n *NodeDetails, err error) {

	// Validate server node
	ok, role, dc := ValidateNode(m)
	if !ok {
		return nil, errors.New("Invalid server node")
	}

	// Get services
	// Services are meant to be encoded in this format: svc-1:port;svc-2:port
	services := make([]NodeService, 0)
	if svcs, ok := m.Tags["services"]; ok {
		for _, svc := range strings.Split(svcs, ";") {
			if strings.Contains(svc, ":") {
				attrs := strings.Split(svc, ":")
				if len(attrs) >= 2 {

					// Convert port to int
					p, e := strconv.Atoi(attrs[1])
					if e != nil {
						err = fmt.Errorf("service port cannot be converted to string: '%s'", p)
						return
					}

					services = append(services, NodeService{
						Name: attrs[0],
						Port: p,
					})
				} else {
					return nil, fmt.Errorf("Invalid service attributes: '%s'", svc)
				}
			} else {
				return nil, fmt.Errorf("Invalid service attributes: '%s'", svc)
			}
		}
	}

	// All nodes which have this tag are bootstrapped
	_, bootstrap := m.Tags["bootstrap"]

	n = &NodeDetails{
		Bootstrap:  bootstrap,
		ID:         m.Tags["id"],
		Name:       m.Name,
		Role:       role,
		DataCenter: dc,
		Addr:       m.Addr,
		Services:   services,
		Status:     m.Status,
		Tags:       m.Tags,
	}
	return
}

// NodeDetails stores details about a single serf.Member
type NodeDetails struct {
	Bootstrap  bool
	ID         string
	Name       string
	Role       string
	DataCenter string
	Addr       net.IP
	Services   []NodeService
	Tags       map[string]string
	Status     serf.MemberStatus
}

func (n NodeDetails) String() (s string) {
	// NodeDetails{Name: "somename", Role: "role", Cluster: "cluster", Addr: "127.0.0.1:9000"}
	if b, err := n.Addr.MarshalText(); err != nil {
		s = fmt.Sprintf("NodeDetails{Name: \"%s\", Role: \"%s\", DataCenter: \"%s\"}", n.Name, n.Role, n.DataCenter)
	} else {
		s = fmt.Sprintf("NodeDetails{Name: \"%s\", Role: \"%s\", DataCenter: \"%s\", Addr: \"%s:%d\"}", n.Name, n.Role, n.DataCenter, string(b))
	}
	return
}

type NodeService struct {
	Name string
	Port int
}
