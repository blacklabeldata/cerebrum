package namedtuple

import (
	"hash"
	"sync"
)

func NewHasher(hasher hash.Hash32) SynchronizedHash {
	var mutex sync.Mutex
	return SynchronizedHash{hasher, mutex}
}

type SynchronizedHash struct {
	hasher hash.Hash32
	mutex  sync.Mutex
}

func (s *SynchronizedHash) Hash(data []byte) uint32 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.hasher.Reset()
	s.hasher.Write(data)
	return s.hasher.Sum32()
}
