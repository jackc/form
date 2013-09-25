package form

import (
	"net/url"
	"testing"
)

func hasError(t *testing.T, actualErrors []error, expected error) {
	if len(actualErrors) == 0 {
		t.Error("Expected a validation error but none occurred")
		return
	}

	var found bool
	for _, e := range actualErrors {
		if e == expected {
			found = true
		}
	}

	if !found {
		t.Errorf("Expected MissingError but it was %T", actualErrors[0])
	}
}

func FieldLint(t *testing.T, f Field) {
	expected := "abcd"
	f.Parse(expected)
	if actual := f.Submission(); actual != expected {
		t.Errorf("Expected Submission() (%#v) to equal value passed to Parse() (%#v)", actual, expected)
	}
}

func TestStringField(t *testing.T) {
	var f StringField

	FieldLint(t, &f)

	f = StringField{MinLength: 0, MaxLength: 100}
	f.Parse("foo")
	if f.Value != "foo" {
		t.Errorf("Expected \"foo\" to parse as \"foo\" but it was %v", f.Value)
	}

	f = StringField{Required: true}
	f.Validate()
	hasError(t, f.Errors(), MissingError)
}

func TestIntField(t *testing.T) {
	var f IntField

	FieldLint(t, &f)

	f = IntField{Min: 0, Max: 100}
	f.Parse("42")
	if f.Value != 42 {
		t.Errorf("Expected \"42\" to parse as 42 but it was %v", f.Value)
	}

	f = IntField{Min: 0, Max: 100, Required: true}
	f.Validate()
	hasError(t, f.Errors(), MissingError)

	f = IntField{Min: 0, Max: 100}
	f.Parse("asdf")
	f.Validate()
	hasError(t, f.Errors(), NotIntError)

	f = IntField{Min: 0, Max: 100}
	f.Parse("-1")
	f.Validate()
	hasError(t, f.Errors(), TooSmallError{Min: 0})

	f = IntField{Min: 0, Max: 100}
	f.Parse("101")
	f.Validate()
	hasError(t, f.Errors(), TooBigError{Max: 100})
}

func TestFormParse(t *testing.T) {
	f := NewForm()
	f.AddField("name", &StringField{MaxLength: 30})
	f.AddField("age", &IntField{Max: 1000})

	values := make(url.Values)
	values["name"] = []string{"David"}
	values["age"] = []string{"30"}
	f.Parse(values)

	if len(f.Errors) > 0 {
		t.Errorf("Unexpected errors parsing form: %v", f.Errors)
	}

	if f.Fields["name"].(*StringField).Value != "David" {
		t.Errorf("Expected \"name\" to be parsed to \"David\" but it was \"%v\"", f.Fields["name"].(*StringField).Value)
	}

	if f.Fields["age"].(*IntField).Value != 30 {
		t.Errorf("Expected \"age\" to be parsed to 30 but it was %v", f.Fields["age"].(*IntField).Value)
	}
}

func TestFormParseWithErrors(t *testing.T) {
	f := NewForm()
	f.AddField("age", &IntField{Max: 1000})

	values := make(url.Values)
	values["age"] = []string{"-1"}
	f.Parse(values)

	if len(f.Errors) != 1 {
		t.Fatalf("Expected 1 error but had %v", len(f.Errors))
	}
}
