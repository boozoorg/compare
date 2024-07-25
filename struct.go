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

type Difference struct {
	Name string `json:"name"`
	Old  string `json:"old"`
	New  string `json:"new"`
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
func TwoEqualStructs[K comparable](first, second K) (difs Differences, err error) {
	// check if given values are struct
	st1, ok := isStruct(json.Marshal(first))
	if !ok {
		return difs, fmt.Errorf("input is not struct but: %T", first)
	}
	st2, _ := isStruct(json.Marshal(second))

	// start comparing
	difs = twoStructsInfo(st1, st2, reflect.TypeOf(first).Name(), reflect.TypeOf(first))

	return difs, nil
}

// recursion to compare structs fields
func twoStructsInfo(st1, st2 map[string]any, stName string, rt reflect.Type) (difs Differences) {
	// get data from first struct
	for key, val := range st1 {
		f, _ := rt.FieldByName(key)
		// check if field of this struct is anther struct,
		// if so, recursively continue
		if _, ok := val.(map[string]any); ok {
			difs = append(difs, twoStructsInfo(val.(map[string]any), st2[key].(map[string]any), stName+"."+key, f.Type)...)
		} else if st2[key] != val {
			// check tag
			_name := f.Tag.Get(tag)
			switch _name {
			case "-":
				continue
			case "":
				_name += stName + "." + key
			}
			// check type of first field, and set the old value
			// check type of second field, and set the new value
			difs = append(difs, Difference{
				Name: _name,
				Old:  field2Text(val),
				New:  field2Text(st2[key]),
			})
		}
	}

	return
}

func field2Text(val any) (text string) {
	rf := getDataFromPointer(reflect.ValueOf(val))
	switch rf.Kind() {
	case reflect.Invalid:
		text = "nil"
	default:
		text = fmt.Sprintf("%v", val)
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
