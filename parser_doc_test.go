package tagparser

import (
	"fmt"
	"log"
)

func ExampleParse() {
	result, err := Parse(`foo:"bar" buz:"qux,foobar"`, true)
	if err != nil {
		log.Fatalf("unexpected error has come: %s", err)
	}
	fmt.Printf("%v\n", result) // => map[foo:bar buz:qux,foobar]
}
