package main

import (
	"fmt"
	"reflect"

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
	fmt.Printf("foo TerminalFields: %# v\n", metaflector.TerminalFields(foo))
	fmt.Printf("Bar.ID resolved to: %v\n", metaflector.Get(foo, "Bar.ID"))

	metaflector.EachField(foo, func(obj interface{}, name string, kind reflect.Kind) {
		fmt.Printf("obj=%v == Get(obj, %q) ? %v\n", obj, name, reflect.DeepEqual(obj, metaflector.Get(foo, name)))
	})

	fmt.Printf("Get()=%v\n", metaflector.Get(foo, "Name"))
}
