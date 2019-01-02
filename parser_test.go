package tagparser

import (
	"testing"
)

func TestBasic(t *testing.T) {
	result, err := Parse(`foo:"bar" buz:"qux,foobar"`)
	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); len(result) != 2 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	expectedDataset := map[string]string{
		"foo": "bar",
		"buz": "qux,foobar",
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestWithEscaping(t *testing.T) {
	result, err := Parse(`foo:"bar" buz:"qux\"foobar" hoge:"fuga"`)

	if err != nil {
		t.Fatalf("unexpected error has come: %s", err)
	}

	if mapLen := len(result); len(result) != 3 {
		t.Fatalf("unexpected length of got map: got = %d", mapLen)
	}

	expectedDataset := map[string]string{
		"foo":  "bar",
		"buz":  `qux"foobar`,
		"hoge": "fuga",
	}

	for key, expected := range expectedDataset {
		if got := result[key]; got != expected {
			t.Errorf(`parsed result of "%s" is not correct: expected = "%s", got = "%s"`, key, expected, got)
		}
	}
}

func TestPragmaticExample(t *testing.T) {
	result, err := Parse(`protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty" datastore:"email" goon:"id"`)

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
	result, err := Parse(`foo:"bar" buz:"qux,すごい,foobar"`)
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
	if _, err := Parse(`foo:"bar" :"qux"`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenValueIsEmpty(t *testing.T) {
	if _, err := Parse(`foo:"bar" buz:`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenKeyValueDelimiterIsMissing(t *testing.T) {
	if _, err := Parse(`foo:"bar" buz`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenValueIsNotTerminated(t *testing.T) {
	if _, err := Parse(`foo:"bar" buz:"qux`); err == nil {
		t.Fatal("expected error has not raised")
	}
}

func TestShouldRaiseErrorWhenKeyContainsWhiteSpace(t *testing.T) {
	if _, err := Parse(`foo:"bar"  bu z:"qux"`); err == nil {
		t.Fatal("expected error has not raised")
	}
}
