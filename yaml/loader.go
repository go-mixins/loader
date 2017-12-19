package yaml

import (
	"io/ioutil"
	"time"

	"github.com/go-fsnotify/fsnotify"
	"gopkg.in/yaml.v2"

	"github.com/go-mixins/loader"
)

// DebounceTimeout defines change event settle time in milliseconds
var DebounceTimeout = 1000

// Loader implements loader.Loader for a YAML file
type Loader struct {
	name          string
	stop, changes chan struct{}
	close         chan error
	watcher       *fsnotify.Watcher
	err           error
}

var _ loader.Loader = (*Loader)(nil)

// New creates Loader initialized with a file name
func New(name string) (res *Loader) {
	res = &Loader{
		name:    name,
		stop:    make(chan struct{}),
		changes: make(chan struct{}),
		close:   make(chan error, 1),
	}
	if res.watcher, res.err = fsnotify.NewWatcher(); res.err != nil {
		res.err = loader.Errors.Wrap(res.err, "creating fsnotify watcher")
		return
	}
	go func() {
		defer close(res.changes)
		defer func() {
			res.close <- res.watcher.Close()
			close(res.close)
		}()
		for {
			select {
			case <-res.stop:
				return
			case <-res.watcher.Events:
			loop:
				for {
					select {
					case <-res.watcher.Events:
						break
					case <-time.After(time.Duration(DebounceTimeout) * time.Millisecond):
						break loop
					}
				}
				res.changes <- struct{}{}
			}
		}
	}()
	if res.err = res.watcher.Add(res.name); res.err != nil {
		res.err = loader.Errors.Wrapf(res.err, "adding %q to watcher", res.name)
		return
	}
	return
}

// Close stops background process and releases FS watcher
func (l *Loader) Close() error {
	close(l.stop)
	return loader.Errors.Wrap(<-l.close, "closing watcher")
}

// Changes provides source of config change events. For environment variables
// there will never be any.
func (l *Loader) Changes() <-chan struct{} {
	return l.changes
}

// Load target object from a YAML file
func (l *Loader) Load(dest interface{}) error {
	if l.err != nil {
		return l.err
	}
	data, err := ioutil.ReadFile(l.name)
	if err != nil {
		return loader.Errors.Wrap(err, "read file")
	}
	return loader.Errors.Wrap(yaml.Unmarshal(data, dest), "unmarshal yml")
}
