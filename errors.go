package cerebrum

import "errors"

var ErrNoLeader = errors.New("No cluster leader")

var ErrUnknownConnType = errors.New("Unknown connection type")
