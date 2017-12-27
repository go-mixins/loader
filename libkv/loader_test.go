package libkv_test

import (
	"testing"

	"github.com/docker/libkv/store"
	"github.com/go-mixins/loader/libkv"
)

type kvMock struct{}

func (kvm *kvMock) Close() {}

func (kvm *kvMock) Get(key string) (res *store.KVPair, err error) {
	switch key {
	case "unknown":
		err = store.ErrKeyNotFound
	case "a/b/c/":
		res = &store.KVPair{
			Key:   key,
			Value: []byte("1"),
		}
	}
	return
}

func (kvm *kvMock) List(directory string) (res []*store.KVPair, err error) {
	switch directory {
	case "a/":
		res = []*store.KVPair{
			&store.KVPair{Key: "a/b/"},
		}
	case "a/b/":
		res = []*store.KVPair{
			&store.KVPair{Key: "a/b/c/"},
		}
	case "a/b/c/":
		res = []*store.KVPair{}
	default:
		err = store.ErrKeyNotFound
	}
	return
}

func TestLoad(t *testing.T) {
	loader := libkv.New("a", &kvMock{})
	var dest struct {
		B struct {
			C int
		}
	}
	defer loader.Close()
	err := loader.Load(&dest)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("%+v", dest)
}
