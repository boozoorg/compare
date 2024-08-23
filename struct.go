package compare

import (
	"encoding/json"
	"fmt"
	"reflect"
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

type Difference struct {
	Name string
	Old  string
	New  string
}

type Differences []Difference

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
// return: [{Name:"name" Old:"Jonny" New:"Bob"}
// {Name:"Person.Age" Old:"20" New:"22"}
// {Name:"Person.Book.Name" Old:"Warcraft" New:"Dune"}
// {Name:"returned" Old:"true" New:"false"}]
func TwoEqualStructs[K comparable](first, second K) (Differences, error) {
	// check if given values are struct
	st1, ok := isStruct(json.Marshal(first))
	if !ok {
		return Differences{}, fmt.Errorf("input is not struct but: %T", first)
	}
	st2, _ := isStruct(json.Marshal(second))
	// start comparing
	return twoEqualStructs(st1, st2, reflect.TypeOf(first).Name(), reflect.TypeOf(second)), nil
}

// recursion to compare structs fields
func twoEqualStructs(st1, st2 map[string]any, stN string, rt reflect.Type) (difs Differences) {
	// get data from first struct
	for k, v := range st1 {
		f, _ := rt.FieldByName(k)
		// check if field of this struct is anther struct,
		// if so, recursively continue
		if _, ok := v.(map[string]any); ok {
			difs = append(difs, twoEqualStructs(v.(map[string]any), st2[k].(map[string]any), stN+"."+k, f.Type)...)
		} else if st2[k] != v {
			// check tag
			tn := f.Tag.Get(tag)
			switch tn {
			case "-":
				continue
			case "":
				tn += stN + "." + k
			}
			// check type of first field, and set the old value
			// check type of second field, and set the new value
			difs = append(difs, Difference{
				Name: tn,
				Old:  field2Text(v),
				New:  field2Text(st2[k]),
			})
		}
	}
	return
}

// format struct value
func field2Text(val any) string {
	switch getDataFromPointer(reflect.ValueOf(val)).Kind() {
	case reflect.Invalid:
		return "nil"
	default:
		return fmt.Sprintf("%v", val)
	}
}

// check is given []byte is struct, if yes return unmarshaled struct
// into map[string]any and true else nil map and false
func isStruct(st []byte, _ error) (m map[string]any, ok bool) {
	return m, json.Unmarshal(st, &m) == nil
}
