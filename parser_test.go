package tagparser

import (
	"reflect"
	"testing"
)

func TestBasic_strict(t *testing.T) {
	customTag := `foo:"bar" buz:"qux,foo  bar"`
	type Item struct {
		elem bool `foo:"bar" buz:"qux,foo  bar"`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	expectedDataset := map[string]string{
		"foo": field.Tag.Get("foo"),
		"buz": field.Tag.Get("buz"),
	}

	result, err := Parse(customTag, true)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 2 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestWithEscaping(t *testing.T) {
	customTag := `foo:"bar" buz:"qux\"foo\\bar" hoge:"fuga"`
	type Item struct {
		elem bool `foo:"bar" buz:"qux\"foo\\bar" hoge:"fuga"`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	expectedDataset := map[string]string{
		"foo":  field.Tag.Get("foo"),
		"buz":  field.Tag.Get("buz"),
		"hoge": field.Tag.Get("hoge"),
	}

	result, err := Parse(customTag, true)

	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); len(result) != 3 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestPragmaticExample(t *testing.T) {
	customTag := `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty" datastore:"email" goon:"id"`
	type Item struct {
		elem bool `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty" datastore:"email" goon:"id"`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	expectedDataset := map[string]string{
		"protobuf":  field.Tag.Get("protobuf"),
		"json":      field.Tag.Get("json"),
		"datastore": field.Tag.Get("datastore"),
		"goon":      field.Tag.Get("goon"),
	}

	result, err := Parse(customTag, true)

	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); len(result) != 4 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestWithMultiByte(t *testing.T) {
	customTag := `はい:"bar" buz:"qux,すごい,foobar"`
	type Item struct {
		elem bool `はい:"bar" buz:"qux,すごい,foobar"`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	expectedDataset := map[string]string{
		"はい":  field.Tag.Get("はい"),
		"buz": field.Tag.Get("buz"),
	}

	result, err := Parse(customTag, true)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); len(result) != 2 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestShouldRaiseErrorWhenKeyIsEmpty(t *testing.T) {
	if _, err := Parse(`foo:"bar" :"qux"`, true); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenValueIsEmpty(t *testing.T) {
	if _, err := Parse(`foo:"bar" buz:`, true); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenKeyValueDelimiterIsMissing(t *testing.T) {
	if _, err := Parse(`foo:"bar" buz`, true); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenValueIsNotTerminated(t *testing.T) {
	if _, err := Parse(`foo:"bar" buz:"qux`, true); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenKeyContainsWhiteSpace(t *testing.T) {
	if _, err := Parse(`foo:"bar"  bu z:"qux"`, true); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestGiveUpWhenKeyIsOmitted(t *testing.T) {
	result, err := Parse(`:"bar" buz:"qux" foobar:"hoge"`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	customTag := `foo:"bar" :"qux" foobar:"hoge"`
	type Item struct {
		elem bool `foo:"bar" :"qux" foobar:"hoge"`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	key := "foo"
	expected := field.Tag.Get(key)

	result, err = Parse(customTag, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 1 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	if got := result[key]; got != expected {
		t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
	}
}

func TestShouldRaiseErrorWhenQuoteForValueIsMissing(t *testing.T) {
	if _, err := Parse(`foo:"bar"  buz:qux"`, true); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestGiveUpWhenValueQuoteIsMissing(t *testing.T) {
	result, err := Parse(`foo:bar" buz:"qux" foobar:"hoge"`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	customTag := `foo:"bar" buz:qux" foobar:"hoge"`
	type Item struct {
		elem bool `foo:"bar" buz:qux" foobar:"hoge"`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	key := "foo"
	expected := field.Tag.Get(key)

	result, err = Parse(customTag, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 1 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	if got := result[key]; got != expected {
		t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
	}
}

func TestGiveUpWhenValueIsMissing(t *testing.T) {
	result, err := Parse(`foo:`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	customTag := `foo:"bar" buz:`
	type Item struct {
		elem bool `foo:"bar" buz:`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	key := "foo"
	expected := field.Tag.Get(key)

	result, err = Parse(customTag, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 1 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	if got := result[key]; got != expected {
		t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
	}
}

func TestGiveUpWhenKeyContainsWhiteSpace(t *testing.T) {
	result, err := Parse(`f oo:"bar" buz:"qux" foobar:"hoge"`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	customTag := `foo:"bar" b uz:"qux" foobar:"hoge"`
	type Item struct {
		elem bool `foo:"bar" b uz:"qux" foobar:"hoge"`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	key := "foo"
	expected := field.Tag.Get(key)

	result, err = Parse(customTag, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 1 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	if got := result[key]; got != expected {
		t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
	}
}

func TestShouldRaiseErrorWhenKeyContainsDoubleQuote(t *testing.T) {
	if _, err := Parse(`"foo":"bar"  buz:"qux"`, true); err == nil {
		t.Fatal("expected error has not raised")
	}

	if _, err := Parse(`foo:"bar"  "buz":"qux"`, true); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestGiveUpWhenKeyContainsDoubleQuote(t *testing.T) {
	result, err := Parse(`"foo":"bar" buz:"qux" foobar:"hoge"`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	customTag := `foo:"bar" "buz":"qux" foobar:"hoge"`
	type Item struct {
		elem bool `foo:"bar" "buz":"qux" foobar:"hoge"`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	key := "foo"
	expected := field.Tag.Get(key)

	result, err = Parse(customTag, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 1 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	if got := result[key]; got != expected {
		t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
	}
}

func TestGiveUpWhenKeyIsNotTerminated(t *testing.T) {
	result, err := Parse(`foo`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	customTag := `foo:"bar" buz`
	type Item struct {
		elem bool `foo:"bar" buz`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	key := "foo"
	expected := field.Tag.Get(key)

	result, err = Parse(customTag, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 1 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	if got := result[key]; got != expected {
		t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
	}
}

func TestGiveUpWhenValueIsNotTerminated(t *testing.T) {
	result, err := Parse(`foo:"bar`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	customTag := `foo:"bar" "buz":"qux`
	type Item struct {
		elem bool `foo:"bar" "buz":"qux`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	key := "foo"
	expected := field.Tag.Get(key)

	result, err = Parse(customTag, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 1 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	if got := result[key]; got != expected {
		t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
	}
}

func TestDuplicatedKey(t *testing.T) {
	customTag := `foo:"bar" hoge:"" foo:"buz" hoge:"test"`
	type Item struct {
		elem bool `foo:"bar" hoge:"" foo:"buz" hoge:"test"`
	}

	item := Item{elem: true}
	typ := reflect.TypeOf(item)
	field := typ.Field(0)

	expectedDataset := map[string]string{
		"foo":  field.Tag.Get("foo"),
		"hoge": field.Tag.Get("hoge"),
	}

	result, err := Parse(customTag, true)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}
