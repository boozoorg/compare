package compare

import (
	"fmt"
	"reflect"
)

// compare to comparable type
func IsEqual[K comparable](first, second K) bool {
	if reflect.ValueOf(first).Kind() == reflect.Pointer {
		return fmt.Sprint(getDataFromPointer(reflect.ValueOf(first))) == fmt.Sprint(getDataFromPointer(reflect.ValueOf(second)))
	}
	return first == second
}

func getDataFromPointer(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer {
		return getDataFromPointer(v.Elem())
	}
	return v
}
