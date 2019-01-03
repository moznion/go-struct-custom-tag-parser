go-custom-tag-parser
==

[![Build Status](https://travis-ci.org/moznion/go-custom-tag-parser.svg?branch=master)](https://travis-ci.org/moznion/go-custom-tag-parser) [![GoDoc](https://godoc.org/github.com/moznion/go-custom-tag-parser?status.svg)](https://godoc.org/github.com/moznion/go-custom-tag-parser)

A simple parser for golang's custom tags.

Synopsis
--

### with strict mode

It raises an error when an unacceptable custom tag is given.

```go
package main

import "github.com/moznion/go-custom-tag-parser"

func main() {
	result, err := Parse(`foo:"bar" buz:"qux,foobar"`, true)
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

import "github.com/moznion/go-custom-tag-parser"

func main() {
	result, _ := Parse(`foo:"bar" buz:"qux,foobar"`, false)
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

