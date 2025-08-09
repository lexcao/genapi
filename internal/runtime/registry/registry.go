package registry

import (
	"fmt"
	"reflect"
	"sync"
)

type key struct {
	pkg  string
	name string
}

type value struct {
	typ    reflect.Type
	config []any
}

var registry = &struct {
	registration map[key]value
	lock         sync.RWMutex
}{
	registration: make(map[key]value),
}

func Register[api any, client any](config ...any) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	clientType := reflect.TypeOf(new(client)).Elem()
	if clientType.Kind() == reflect.Pointer {
		clientType = clientType.Elem()
	}

	registry.registration[getKey[api]()] = value{
		typ:    clientType,
		config: config,
	}
}

func New[api any]() (api, any) {
	registry.lock.RLock()
	defer registry.lock.RUnlock()

	key := getKey[api]()
	value, ok := registry.registration[key]
	if !ok {
		panic(fmt.Sprintf(`
genapi: no registration found for interface %s.%s

This usually means:
1. You forgot to run 'go generate' on your API package
2. The generated *.gen.go file wasn't imported
3. There's a bug in code generation

Run: go generate ./...
`, key.pkg, key.name))
	}

	clientImpl := reflect.New(value.typ).Interface().(api)

	if len(value.config) == 0 {
		return clientImpl, nil
	}

	return clientImpl, value.config[0]
}

func getKey[api any]() key {
	apiType := reflect.TypeOf(new(api)).Elem()
	return key{
		pkg:  apiType.PkgPath(),
		name: apiType.Name(),
	}
}
