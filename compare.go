package compare

import (
	"fmt"
	"reflect"
)

// compare to comparable type
func IsEqual[K comparable](first, second K) bool {
	if reflect.ValueOf(first).Kind() == reflect.Pointer {
		f, s := getDataFromPointer(reflect.ValueOf(first)), getDataFromPointer(reflect.ValueOf(second))
		return fmt.Sprint(f) == fmt.Sprint(s) && (f.IsValid() == s.IsValid())
	}
	return first == second
}

// recursively get value from pointer
func getDataFromPointer(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer {
		return getDataFromPointer(v.Elem())
	}
	return v
}
