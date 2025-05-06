//go:build !solution

package reversemap

import (
	"reflect"
)

func ReverseMap(forward interface{}) interface{} {
	forwardVal := reflect.ValueOf(forward)
	forwardType := reflect.TypeOf(forward)

	key := forwardType.Key()
	value := forwardType.Elem()
	newMap := reflect.MakeMap(reflect.MapOf(value, key))

	for _, keyOldMap := range forwardVal.MapKeys() {
		newMap.SetMapIndex(forwardVal.MapIndex(keyOldMap), keyOldMap)
	}

	return newMap.Interface()
}
