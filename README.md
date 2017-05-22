# metaflector

[![Documentation](https://godoc.org/github.com/gigawattio/metaflector?status.svg)](https://godoc.org/github.com/gigawattio/metaflector)
[![Build Status](https://travis-ci.org/gigawattio/metaflector.svg?branch=master)](https://travis-ci.org/gigawattio/metaflector)
[![Report Card](https://goreportcard.com/badge/github.com/gigawattio/metaflector)](https://goreportcard.com/report/github.com/gigawattio/metaflector)

### About

Go (golang) package which provides reflection abstractions beyond the Go standard library.

Metaflector makes it easy to inspect objects and structs, and programmatically access structural metadata.

To be precise, this package currently provides:

* Generating a dot-separated list of struct hierarchy

* Iterating over a struct or slice / array object's fields

I've found this functionality useful for automatically applying user input as search filters against arbitrary structs.

See the [docs](https://godoc.org/github.com/gigawattio/metaflector) for more info.

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
```

Output:

```
[]string{"Bar.ID", "Name"}
Bar.ID resolved to: 2017
```

See the [tests](https://github.com/gigawattio/metaflector/blob/master/terminal_fields_test.go#L56-L100) for more examples.

### Related Work

The [reflections](https://github.com/oleiade/reflections) package is related and complementary.

#### License

Permissive MIT license, see the [LICENSE](LICENSE) file for more information.

