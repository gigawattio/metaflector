# metaflect

[![Documentation](https://godoc.org/github.com/gigawattio/metaflect?status.svg)](https://godoc.org/github.com/gigawattio/metaflect)
[![Build Status](https://travis-ci.org/gigawattio/metaflect.svg?branch=master)](https://travis-ci.org/gigawattio/metaflect)
[![Report Card](https://goreportcard.com/badge/github.com/gigawattio/metaflect)](https://goreportcard.com/report/github.com/gigawattio/metaflect)

### About

Go (golang) package which provides reflection abstractions beyond the Go standard library.

Metaflect makes it easy to inspect objects and structs, and programmatically access structural metadata.

To be precise, this package currently provides:

* Generating a dot-separated list of struct hierarchy

* Iterating over a struct or slice / array object's fields

### Requirements

* Go version 1.6 or newer

### Running the test suite

    go test ./...

### Example Usage

```go
package main

import (
    "fmt"

    "github.com/gigawattio/metaflect"
)

type (
    Foo struct {
        Bar Bar
        Name string
    }

    Bar struct {
        ID string
        private string
    }
)

var foo = &Foo{
    Bar: Bar{
        ID: "2017",
        private: "this isn't exported",
    },
    Name: "meta",
}

func main() {
    fmt.Printf("%# v\n", metaflect.TerminalFields(foo))
}
```

Output:

```
[]string{"Bar.ID", "Name"}
```

### Related Work

The [reflections](https://github.com/oleiade/reflections) package is related and complementary.

#### License

Permissive MIT license, see the [LICENSE](LICENSE) file for more information.

