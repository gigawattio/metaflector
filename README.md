# metaflector

[![Documentation](https://godoc.org/github.com/gigawattio/metaflector?status.svg)](https://godoc.org/github.com/gigawattio/metaflector)
[![Build Status](https://travis-ci.org/gigawattio/metaflector.svg?branch=master)](https://travis-ci.org/gigawattio/metaflector)
[![Report Card](https://goreportcard.com/badge/github.com/gigawattio/metaflector)](https://goreportcard.com/report/github.com/gigawattio/metaflector)

### About

Go (golang) package which provides reflection abstractions beyond the Go standard library.

Metaflector makes it easy to inspect objects and structs, and programmatically access structural data.

To be precise, this package currently provides:

* Generating a dot-separated list of struct hierarchical properties

```go
metaflector.TerminalFields(&http.Server{})
// Output: []string{"Addr", "IdleTimeout", "MaxHeaderBytes", "ReadHeaderTimeout", "ReadTimeout", "WriteTimeout"}
```

* Iterating over a struct or slice or array objects' fields

```go
metaflector.EachField(http.Server{}, func(obj interface{}, name string, kind reflect.Kind) {
    fmt.Printf("obj=%v name=%v kind=%v\n", obj, name, kind)
})

// Output:
// obj= name=Addr kind=string
// obj=<nil> name=TLSConfig kind=struct
// obj=0 name=ReadTimeout kind=int64
// obj=0 name=ReadHeaderTimeout kind=int64
// obj=0 name=WriteTimeout kind=int64
// obj=0 name=IdleTimeout kind=int64
// obj=0 name=MaxHeaderBytes kind=int
// obj=map[] name=TLSNextProto kind=map
// obj=<nil> name=ErrorLog kind=struct
```

* Dynamic property extraction based on dot-paths

e.g.
```go
Get(myVar, "A.Nested.Property")
```

I've found this functionality useful for automatically applying user input as search filters against arbitrary structs in command-line progreams.

See the [docs](https://godoc.org/github.com/gigawattio/metaflector) for more info.

Created by [Jay Taylor](https://jaytaylor.com/) and used by [Gigawatt](https://gigawatt.io/).

#### A word about current limitations

* For heterogeneous collections (i.e. this is possible via `[]interface{}`), only the structure of the first non-nil slice or array element will be considered.

* No maps support _[yet]_.

### Requirements

* Go version 1.6 or newer

### Running the test suite

    go test ./...

### Example Usage

[examples/basic.go](examples/basic.go)

```go
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

// Output:
// foo TerminalFields: []string{"Bar.ID", "Name"}
// Bar.ID resolved to: 2017
// obj={2017 this isn't exported} == Get(obj, "Bar") ? true
// obj=meta == Get(obj, "Name") ? true
// Get()=meta
```

See the [tests](https://github.com/gigawattio/metaflector/blob/master/terminal_fields_test.go#L56-L100) for more examples.

### Related Work

The [reflections](https://github.com/oleiade/reflections) package is related and complementary.

#### License

Permissive MIT license, see the [LICENSE](LICENSE) file for more information.

