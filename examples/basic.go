package main

import (
	"fmt"

	"github.com/gigawattio/metaflector"
)

type (
	Foo struct {
		Bar  Bar
		Name string
	}

	Bar struct {
		ID      string
		private string
	}
)

var foo = &Foo{
	Bar: Bar{
		ID:      "2017",
		private: "this isn't exported",
	},
	Name: "meta",
}

func main() {
	fmt.Printf("%# v\n", metaflector.TerminalFields(foo))
	fmt.Printf("Bar.ID resolved to: %v\n", metaflector.Get(foo, "Bar.ID"))
}
