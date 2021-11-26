# go-dynamic

go-dynamic is a library that provides dynamic features of Go language.

## Installation

To install go-dynamic, use go get:

```shell
go get -u github.com/garenchan/go-dynamic
```

## Features

- dynamic function calls

## Usage

### Demo 1: dynamic function calls

```go
package main

import (
	"fmt"

	"github.com/garenchan/go-dynamic"
)

type Endpoint struct {
}

func (ep *Endpoint) Add(a, b int) int {
	return a + b
}

func main() {
	ep := &Endpoint{}
	result, _ := dynamic.Call(ep, "Add", 1, 2)

	// Here will print 3.
	fmt.Println(result[0])
}
```
