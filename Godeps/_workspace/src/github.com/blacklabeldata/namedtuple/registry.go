package namedtuple

import (
	"hash/fnv"
	"sync"
)

var DefaultRegistry Registry

func init() {
	DefaultRegistry = NewRegistry()
}

func NewRegistry() Registry {
	return Registry{content: make(map[uint64]TupleType), hasher: NewHasher(fnv.New32a())}
}

type Registry struct {
	content map[uint64]TupleType
	hasher  SynchronizedHash
	mutex   sync.Mutex
}

func (r *Registry) Contains(t TupleType) bool {
	return r.ContainsHash(r.typeSignature(t.Namespace, t.Name))
}

func (r *Registry) ContainsHash(hash uint64) bool {
	// lock registry
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.content[hash]
	return exists
}

func (r *Registry) ContainsName(namespace, name string) bool {
	return r.ContainsHash(r.typeSignature(namespace, name))
}

func (r *Registry) Get(namespace, name string) (tupleType TupleType, exists bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	tupleType, exists = r.content[r.typeSignature(namespace, name)]
	return
}

func (r *Registry) GetWithHash(namespace, name uint32) (tupleType TupleType, exists bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	tupleType, exists = r.content[r.typeSignatureHash(namespace, name)]
	return
}

func (r *Registry) Register(t TupleType) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	hash := r.typeSignature(t.Namespace, t.Name)
	if _, exists := r.content[hash]; !exists {
		r.content[hash] = t
	}
}

func (r *Registry) Unregister(t TupleType) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	hash := r.typeSignature(t.Namespace, t.Name)
	if _, exists := r.content[hash]; exists {
		delete(r.content, hash)
	}
}

func (r *Registry) Size() int {
	return len(r.content)
}

func (r *Registry) typeSignature(namespace, typename string) (hash uint64) {
	// Calculate hashes
	nhash := r.hasher.Hash([]byte(namespace))
	thash := r.hasher.Hash([]byte(typename))

	hash = r.typeSignatureHash(nhash, thash)
	return
}

func (r *Registry) typeSignatureHash(nhash, thash uint32) (hash uint64) {

	// Combine hashes
	hash = uint64(nhash) << 32
	hash |= uint64(thash)
	return
}
