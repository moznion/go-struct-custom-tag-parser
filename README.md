go-struct-custom-tag-parser
==

[![Build Status](https://travis-ci.org/moznion/go-struct-custom-tag-parser.svg?branch=master)](https://travis-ci.org/moznion/go-struct-custom-tag-parser) [![GoDoc](https://godoc.org/github.com/moznion/go-struct-custom-tag-parser?status.svg)](https://godoc.org/github.com/moznion/go-struct-custom-tag-parser)

A simple parser for golang's struct custom tags.

STOP!!!!
--

You don't have to use this library. You should use [reflect.StructTag](https://golang.org/pkg/reflect/#StructTag) instead of this.

This library's processing time is slower than reflect's one.

```
goos: darwin
goarch: amd64
pkg: github.com/moznion/go-struct-custom-tag-parser/author
BenchmarkParser-4        1000000              1932 ns/op             848 B/op         11 allocs/op
BenchmarkReflect-4      10000000               140 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/moznion/go-struct-custom-tag-parser/author   4.484s
```

(benchmark code is here: [author/parser_bench_test.go](/author/parser_bench_test.go))

__There is only one reason for using this library, this has a feature to check the syntax of custom tag (in strict mode, `Parse()` returns an error when a custom tag is an invalid syntax).__

Synopsis
--

### with strict mode

It raises an error when an unacceptable custom tag is given on parsing.

```go
package main

import (
	"fmt"
	"log"

	"github.com/moznion/go-struct-custom-tag-parser"
)

func main() {
	result, err := tagparser.Parse(`foo:"bar" buz:"qux,foobar"`, true)
	if err != nil {
		log.Fatalf("unexpected error has come: %s", err)
	}
	fmt.Printf("%v\n", result) // => map[foo:bar buz:qux,foobar]
}
```

### with no strict mode

It immediately returns the processed results until just before the invalid custom tag syntax. It never raises any error.

```go
package main

import (
	"fmt"

	"github.com/moznion/go-struct-custom-tag-parser"
)

func main() {
	result, _ := tagparser.Parse(`foo:"bar" buz:"qux,foobar"`, false)
	fmt.Printf("%v\n", result) // => map[foo:bar buz:qux,foobar]
}
```

License
--

```
The MIT License (MIT)
Copyright © 2019 moznion, http://moznion.net/ <moznion@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```

