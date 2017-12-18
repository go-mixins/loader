package env

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/go-mixins/loader"
)

// Loader implements loader.Loader
type Loader struct {
	prefix string
}

var _ loader.Loader = (*Loader)(nil)

// Load loads the target from environment
func (l *Loader) Load(dest interface{}) error {
	return errors.Wrap(envconfig.Process(l.prefix, dest), "unmarshal yml")
}

// New creates loader initialized with environment variable prefix
func New(prefix string) *Loader {
	return &Loader{
		prefix: strings.ToUpper(prefix),
	}
}
