package tagparser

import (
	"testing"
)

func TestBasic(t *testing.T) {
	result, err := ParseStrict(`foo:"bar" buz:"qux,foo  bar"`)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); len(result) != 2 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	expectedDataset := map[string]string{
		"foo": "bar",
		"buz": "qux,foo  bar",
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestWithEscaping(t *testing.T) {
	result, err := ParseStrict(`foo:"bar" buz:"qux\"foo\\bar" hoge:"fuga"`)

	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); len(result) != 3 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	expectedDataset := map[string]string{
		"foo":  "bar",
		"buz":  `qux"foo\bar`,
		"hoge": "fuga",
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestPragmaticExample(t *testing.T) {
	result, err := ParseStrict(`protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty" datastore:"email" goon:"id"`)

	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); len(result) != 4 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	expectedDataset := map[string]string{
		"protobuf":  "bytes,1,opt,name=email,proto3",
		"json":      "email,omitempty",
		"datastore": "email",
		"goon":      "id",
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestWithMultiByte(t *testing.T) {
	result, err := ParseStrict(`foo:"bar" buz:"qux,すごい,foobar"`)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); len(result) != 2 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	expectedDataset := map[string]string{
		"foo": "bar",
		"buz": "qux,すごい,foobar",
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestShouldRaiseErrorWhenKeyIsEmpty(t *testing.T) {
	if _, err := ParseStrict(`foo:"bar" :"qux"`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenValueIsEmpty(t *testing.T) {
	if _, err := ParseStrict(`foo:"bar" buz:`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenKeyValueDelimiterIsMissing(t *testing.T) {
	if _, err := ParseStrict(`foo:"bar" buz`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenValueIsNotTerminated(t *testing.T) {
	if _, err := ParseStrict(`foo:"bar" buz:"qux`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenKeyContainsWhiteSpace(t *testing.T) {
	if _, err := ParseStrict(`foo:"bar"  bu z:"qux"`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestGiveUpWhenKeyIsOmitted(t *testing.T) {
	result, err := Parse(`:"bar" buz:"qux" foobar:"hoge"`)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	result, err = Parse(`foo:"bar" :"qux" foobar:"hoge"`)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 1 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	key := "foo"
	expected := "bar"
	if got := result[key]; got != expected {
		t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
	}
}

func TestShouldRaiseErrorWhenQuoteForValueIsMissing(t *testing.T) {
	if _, err := ParseStrict(`foo:"bar"  buz:qux"`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestGiveUpWhenValueQuoteIsMissing(t *testing.T) {
	result, err := Parse(`foo:bar" buz:"qux" foobar:"hoge"`)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	result, err = Parse(`foo:"bar" buz:qux" foobar:"hoge"`)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 1 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	key := "foo"
	expected := "bar"
	if got := result[key]; got != expected {
		t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
	}
}
