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

var registry = &struct {
	registration map[key]reflect.Type
	lock         sync.RWMutex
}{
	registration: make(map[key]reflect.Type),
}

func Register[api any, client any]() {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	clientType := reflect.TypeOf(new(client)).Elem()
	if clientType.Kind() == reflect.Pointer {
		clientType = clientType.Elem()
	}

	registry.registration[getKey[api]()] = clientType
}

func New[api any]() api {
	registry.lock.RLock()
	defer registry.lock.RUnlock()

	key := getKey[api]()
	clientType := registry.registration[key]
	if clientType == nil {
		panic(fmt.Sprintf("no registration for key: %s", key))
	}

	return reflect.New(clientType).Interface().(api)
}

func getKey[api any]() key {
	apiType := reflect.TypeOf(new(api)).Elem()
	return key{
		pkg:  apiType.PkgPath(),
		name: apiType.Name(),
	}
}
