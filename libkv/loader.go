package libkv

import (
	"strings"

	"github.com/docker/libkv/store"

	"github.com/go-mixins/loader"
)

// Loader implements loader.Loader
type Loader struct {
	store   kvStore
	prefix  string
	changes chan struct{}
}

var _ loader.Loader = (*Loader)(nil)

type kvStore interface {
	// Get a value given its key
	Get(key string) (*store.KVPair, error)
	// List the content of a given prefix
	List(directory string) ([]*store.KVPair, error)
	// Close the store connection
	Close()
}

// Close closes underlying changes channel
func (l *Loader) Close() error {
	l.store.Close()
	close(l.changes)
	return nil
}

// Changes provides source of config change events. For environment variables
// there will never be any.
func (l *Loader) Changes() <-chan struct{} {
	return l.changes
}

// New creates loader initialized with KV store prefix
func New(prefix string, store kvStore) *Loader {
	return &Loader{
		store:   store,
		prefix:  strings.Trim(prefix, "/"),
		changes: make(chan struct{}),
	}
}
