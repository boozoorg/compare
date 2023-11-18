package compare

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Compare two equal struct between each other and while finding differences
// return name of struct, old data and new data in other case error or
// "no difference" which mean the variables same
//
// example:
//
//	type Person struct {
//	    Name string        |  f:  "Jonny"       |  s:  "Bob"
//	    Age  int8          |      20            |      22
//	    Book struct {      |                    |
//	        Name string    |      "Warcraft"    |      "Dune"
//	        Returnd bool   |      true          |      false
//	    }                  |                    |
//	}
//
// compare.TwoEqualStructs(f, s)
//
// return: Person.Name was "Jonny" and now "Bob",
// Person.Age was 20 and now 22,
// Person.Book.Name was "Warcraft" and now "Dune"
// Person.Book.Returnd was true and now false
func TwoEqualStructs[K comparable](first, second K) (string, error) {
	map1, ok := isStruct(json.Marshal(first))
	if !ok {
		return "", fmt.Errorf("first input is not struct but %T", first)
	}
	map2, ok := isStruct(json.Marshal(second))
	if !ok {
		return "", fmt.Errorf("second input is not struct but %T", first)
	}

	text := twoStructsinfo(map1, map2, reflect.TypeOf(first).Name())
	if text != "" {
		return text[:len(text)-2], nil
	}

	return "no difference", nil
}

// recursion to compare structs fields
func twoStructsinfo(map1, map2 map[string]any, stName string) (text string) {
	for k, v := range map1 {
		if _, ok := v.(map[string]any); ok {
			text += twoStructsinfo(v.(map[string]any), map2[k].(map[string]any), stName+"."+k)
		} else if map2[k] != v {
			text += stName + "." + k + " was "
			switch v.(type) {
			case int, int8, int16, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, bool:
				text += fmt.Sprintf(`%v`, v)
			case string:
				text += fmt.Sprintf(`"%v"`, v)
			case nil:
				text += "nil"
			}

			text += " and now "
			switch map2[k].(type) {
			case int, int8, int16, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, bool:
				text += fmt.Sprintf(`%v, `, map2[k])
			case string:
				text += fmt.Sprintf(`"%v", `, map2[k])
			case nil:
				text += "nil, "
			}
		}
	}

	return
}

// Check is given variable is struct
func IsStruct(s any) (ok bool) {
	if strings.Contains(fmt.Sprintf("%T", s), "map") {
		return false
	}
	_, ok = isStruct(json.Marshal(s))
	return
}

// check is given []byte is struct, if yes return unmarshaled struct into map[string]any and true else nil map and false
func isStruct(st []byte, _ error) (map[string]any, bool) {
	var m map[string]any
	return m, json.Unmarshal(st, &m) == nil
}
