package tagparser

import (
	"reflect"
	"testing"

	"github.com/moznion/go-struct-custom-tag-parser"
)

const tag = `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty" datastore:"email" goon:"id"`

func BenchmarkParser(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tagValue, _ := tagparser.Parse(tag, true)
		nop(tagValue["protobuf"])
	}
}

func BenchmarkReflect(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tagValue := reflect.StructTag(tag)
		nop(tagValue.Get("protobuf"))
	}
}

func nop(s string) {
}
