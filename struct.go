package compare

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// struct tag name, if want to skip field for checking, add tag with "-" value
// else will be name of field
//
// example:
//
//	type Person struct {
//	   Name string `boo:"-"`
//	}
const tag = "boo"

// Compare two equal struct between each other and while finding differences
// return name of struct, old data and new data in other case error or
// "no difference" which mean the variables same
//
// example:
//
//	type Person struct {
//	    Name string `boo:"name"`           | f:  "Jonny"       | s:  "Bob"
//	    Age  uint8                         |     20            |     22
//	    Book struct {                      |                   |
//	        Name     string                |     "Warcraft"    |     "Dune"
//	        Returned bool `boo:"returned"` |     true          |     false
//	    }                                  |                   |
//	}
//
// compare.TwoEqualStructs(f, s)
//
// return: name was "Jonny" and now "Bob",
// Person.Age was 20 and now 22,
// Person.Book.Name was "Warcraft" and now "Dune",
// returned was true and now false
func TwoEqualStructs[K comparable](first, second K) (string, error) {
	// check if given values are struct
	st1, ok := isStruct(json.Marshal(first))
	if !ok {
		return "", fmt.Errorf("first input is not struct but %T", first)
	}
	st2, ok := isStruct(json.Marshal(second))
	if !ok {
		return "", fmt.Errorf("second input is not struct but %T", first)
	}

	// start comparing
	text := twoStructsinfo(st1, st2, reflect.TypeOf(first).Name(), reflect.TypeOf(first))
	if text != "" {
		return text[:len(text)-2], nil
	}

	return "no difference", nil
}

// recursion to compare structs fields
func twoStructsinfo(st1, st2 map[string]any, stName string, rt reflect.Type) (text string) {
	// get data from first struct
	for key, val := range st1 {
		f, _ := rt.FieldByName(key)
		// check if field of this struct is anther struct,
		// if so, recursively continue
		if _, ok := val.(map[string]any); ok {
			text += twoStructsinfo(val.(map[string]any), st2[key].(map[string]any), stName+"."+key, f.Type)
		} else if st2[key] != val {
			// check tag
			tagVal := f.Tag.Get(tag)
			if tagVal == "-" {
				continue
			}
			if tagVal != "" {
				text += tagVal + " was"
			} else {
				text += stName + "." + key + " was"
			}

			// check type of first field, and set the past value
			switch val.(type) {
			case nil:
				text += " nil"
			case string:
				text += fmt.Sprintf(` "%v"`, val)
			default:
				text += fmt.Sprintf(` %v`, val)
			}

			// check type of second field, and set the present value
			text += " and now"
			switch st2[key].(type) {
			case nil:
				text += " nil, "
			case string:
				text += fmt.Sprintf(` "%v", `, st2[key])
			default:
				text += fmt.Sprintf(` %v, `, st2[key])
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
