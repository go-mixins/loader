package yaml

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/go-mixins/loader"
)

// Loader implements loader.Loader for a YAML file
type Loader struct {
	fn string
}

var _ loader.Loader = (*Loader)(nil)

// New creates Loader initialized with a file name
func New(fn string) *Loader {
	return &Loader{fn}
}

// Load target object from a YAML file
func (l *Loader) Load(dest interface{}) error {
	data, err := ioutil.ReadFile(l.fn)
	if err != nil {
		return errors.Wrap(err, "read file")
	}
	return errors.Wrap(yaml.Unmarshal(data, dest), "unmarshal yml")
}
