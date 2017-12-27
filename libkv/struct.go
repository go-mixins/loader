package libkv

import (
	"fmt"
	"strings"

	"github.com/docker/libkv/store"
	"github.com/go-mixins/loader"
	"github.com/mitchellh/mapstructure"
)

// Load loads the target from environment
func (l *Loader) Load(dest interface{}) error {
	cfg := &mapstructure.DecoderConfig{
		Result:           dest,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return loader.Errors.Wrap(err, "creating map decoder")
	}
	data, err := l.getRecursive(l.prefix)
	if err != nil {
		return err
	}
	fmt.Printf("*** %+v\n", data)
	return loader.Errors.Wrap(decoder.Decode(data), "decoding values")
}

func put(dest map[string]interface{}, prefix string, val interface{}) {
	path := strings.Split(strings.Trim(prefix, "/"), "/")
	fmt.Printf("&&& %s %#v\n", prefix, path)
	if len(path) == 0 {
		return
	}
	for i := len(path) - 1; i >= 1; i-- {
		val = map[string]interface{}{
			path[i]: val,
		}
	}
	dest[path[0]] = val
}

func (l *Loader) getRecursive(prefix string) (interface{}, error) {
	fmt.Printf("### %s\n", prefix)
	pairs, err := l.store.List(prefix)
	if err == store.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, loader.Errors.Wrap(err, "getting KV list")
	}
	if len(pairs) == 0 {
		val, err := l.store.Get(prefix)
		if err != nil {
			return nil, loader.Errors.Wrap(err, "getting key value")
		}
		if val != nil {
			return string(val.Value), nil
		}
		return nil, nil
	}
	res := make(map[string]interface{})
	for _, p := range pairs {
		val, err := l.getRecursive(p.Key)
		if err != nil {
			return res, err
		}
		put(res, strings.TrimPrefix(p.Key, prefix), val)
	}
	return res, nil
}
