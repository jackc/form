package form

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Field interface {
	Parse(value string)
	Submission() string
	Validate()
	Errors() []error
}

var MissingError = errors.New("is missing")
var NotIntError = errors.New("is not an integer")

type TooSmallError struct {
	Min int64
}

func (e TooSmallError) Error() string {
	return fmt.Sprintf("is too small (minimum: %d)", e.Min)
}

type TooBigError struct {
	Max int64
}

func (e TooBigError) Error() string {
	return fmt.Sprintf("is too big (maximum: %d)", e.Max)
}

type IntField struct {
	submission    string
	submitted     bool
	unconvertable bool
	Value         int64
	Required      bool
	Min           int64
	Max           int64
	errors        []error
}

func (f *IntField) Parse(value string) {
	var err error
	f.submission = value
	f.submitted = true
	f.Value, err = strconv.ParseInt(value, 10, 64)
	if err != nil {
		f.unconvertable = true
	}
}

func (f *IntField) Submission() string {
	return f.submission
}

func (f *IntField) Validate() {
	f.errors = make([]error, 0)

	if f.Required && !f.submitted {
		f.errors = append(f.errors, MissingError)
		return
	}

	if f.unconvertable {
		f.errors = append(f.errors, NotIntError)
		return
	}

	if f.Value < f.Min {
		f.errors = append(f.errors, TooSmallError{Min: f.Min})
	}

	if f.Value > f.Max {
		f.errors = append(f.errors, TooBigError{Max: f.Max})
	}
}

func (f *IntField) Errors() []error {
	return f.errors
}

type StringField struct {
	submitted bool
	Value     string
	Required  bool
	MinLength int64
	MaxLength int64
	errors    []error
}

func (f *StringField) Parse(value string) {
	f.Value = value
	f.submitted = true
}

func (f *StringField) Submission() string {
	return f.Value
}

func (f *StringField) Validate() {
	f.errors = make([]error, 0)

	if f.Required && !f.submitted {
		f.errors = append(f.errors, MissingError)
		return
	}
}

func (f *StringField) Errors() []error {
	return f.errors
}

type Form struct {
	Fields map[string]Field
	Errors []error
}

type FieldError struct {
	FieldName  string
	FieldError error
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %v", e.FieldName, e.FieldError)
}

func NewForm() (f *Form) {
	f = &Form{}
	f.Fields = make(map[string]Field)
	f.Errors = make([]error, 0)
	return
}

func (f *Form) AddField(name string, field Field) {
	f.Fields[name] = field
}

func (f *Form) Parse(values url.Values) {
	f.Errors = make([]error, 0)

	for name, field := range f.Fields {
		if v, ok := values[name]; ok {
			field.Parse(v[0])
		}
		field.Validate()
		for _, e := range field.Errors() {
			f.Errors = append(f.Errors, FieldError{FieldName: name, FieldError: e})
		}
	}
}
