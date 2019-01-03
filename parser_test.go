package tagparser

import (
	"testing"
)

func TestBasic(t *testing.T) {
	result, err := Parse(`foo:"bar" buz:"qux,foo  bar"`, true)
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
	result, err := Parse(`foo:"bar" buz:"qux\"foo\\bar" hoge:"fuga"`, true)

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
	result, err := Parse(`protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty" datastore:"email" goon:"id"`, true)

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
	result, err := Parse(`foo:"bar" buz:"qux,すごい,foobar"`, true)
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

	result, err = Parse(`foo:"bar" :"qux" foobar:"hoge"`, false)
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

	result, err = Parse(`foo:"bar" buz:qux" foobar:"hoge"`, false)
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

func TestGiveUpWhenValueIsMissing(t *testing.T) {
	result, err := Parse(`foo:`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	result, err = Parse(`foo:"bar" buz:`, false)
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

func TestGiveUpWhenKeyContainsWhiteSpace(t *testing.T) {
	result, err := Parse(`f oo:"bar" buz:"qux" foobar:"hoge"`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	result, err = Parse(`foo:"bar" b uz:"qux" foobar:"hoge"`, false)
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

	result, err = Parse(`foo:"bar" "buz":"qux" foobar:"hoge"`, false)
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

func TestGiveUpWhenKeyIsNotTerminated(t *testing.T) {
	result, err := Parse(`foo`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	result, err = Parse(`foo:"bar" buz`, false)
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

func TestGiveUpWhenValueIsNotTerminated(t *testing.T) {
	result, err := Parse(`foo:"bar`, false)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); mapLen != 0 {
		t.Fatal("got non empty map")
	}

	result, err = Parse(`foo:"bar" "buz":"qux`, false)
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
