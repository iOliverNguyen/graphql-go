package utilities

import (
	"errors"
	"fmt"
	"reflect"
)

func IsNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

func throw(format string, args ...interface{}) {
	panic(errors.New(fmt.Sprintf(format, args...)))
}

func Find(list []interface{}, predicate func(interface{}) bool) interface{} {
	for i := 0; i < len(list); i++ {
		if predicate(list[i]) {
			return list[i]
		}
	}
	return nil
}
