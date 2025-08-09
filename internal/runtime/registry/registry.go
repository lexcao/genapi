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

func New[api any]() (api, any, error) {
	registry.lock.RLock()
	defer registry.lock.RUnlock()

	key := getKey[api]()
	value, ok := registry.registration[key]
	if !ok {
		var zero api
		return zero, nil, fmt.Errorf("no registration found for interface %s.%s", key.pkg, key.name)
	}

	clientImpl := reflect.New(value.typ).Interface().(api)

	if len(value.config) == 0 {
		return clientImpl, nil, nil
	}

	return clientImpl, value.config[0], nil
}

func getKey[api any]() key {
	apiType := reflect.TypeOf(new(api)).Elem()
	return key{
		pkg:  apiType.PkgPath(),
		name: apiType.Name(),
	}
}
