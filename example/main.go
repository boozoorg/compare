package main

import (
	"fmt"

	"github.com/boozoorg/compare"
)

type Person struct {
	ID   uint64 `boo:"-"`
	Name string `boo:"name"`
	Age  uint8
	Book Book
}

type Book struct {
	Name     string `boo:"book_name"`
	Returned bool
}

func main() {
	var f = Person{
		ID: 1, Name: "boozoorg", Age: 22, Book: Book{Name: "WWW", Returned: true},
	}
	var s = Person{
		ID: 2, Name: "buzurg", Age: 23, Book: Book{Name: "XXX", Returned: false},
	}
	resp, _ := compare.TwoEqualStructs(f, s)
	fmt.Printf("%+v\n", resp)
}
