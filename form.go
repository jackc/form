package form

import (
	"errors"
	"net/url"
	"strconv"
)

type FieldTemplate interface {
	GetName() string
	Parse(unparsed string) (interface{}, error)
	Validate(value interface{}) error
}

var MissingError = errors.New("is missing")

type TooShortError struct {
	Minimum int
}

func (e TooShortError) Error() string {
	return "Too short"
}

type TooLongError struct {
	Maximum int
}

func (e TooLongError) Error() string {
	return "Too long"
}

type StringTemplate struct {
	Name      string
	Required  bool
	MinLength int
	MaxLength int
}

func (f *StringTemplate) GetName() string {
	return f.Name
}

func (f *StringTemplate) Parse(unparsed string) (interface{}, error) {
	return unparsed, nil
}

func (f *StringTemplate) Validate(value interface{}) (err error) {
	if value == nil || value == "" {
		if f.Required {
			return MissingError
		} else {
			return nil
		}
	}

	v := value.(string)

	if len(v) < f.MinLength {
		return TooShortError{Minimum: f.MinLength}
	}

	if f.MaxLength < len(v) {
		return TooLongError{Maximum: f.MaxLength}
	}

	return
}

type IntTemplate struct {
	Name     string
	Required bool
	Minimum  int64
	Maximum  int64
}

func (f *IntTemplate) GetName() string {
	return f.Name
}

func (f *IntTemplate) Parse(unparsed string) (interface{}, error) {
	if unparsed == "" {
		return nil, nil
	}

	if parsed, err := strconv.ParseInt(unparsed, 10, 64); err == nil {
		return parsed, err
	} else {
		return nil, err
	}
}

func (f *IntTemplate) Validate(value interface{}) error {
	if f.Required && value == nil {
		return MissingError
	}

	v := value.(int64)

	if v < f.Minimum {
		return errors.New("Too small")
	}

	if f.Maximum < v {
		return errors.New("Too big")
	}

	return nil
}

type FormTemplate struct {
	fieldTemplates map[string]FieldTemplate
	CustomValidate func(*Form)
}

func NewFormTemplate() (f *FormTemplate) {
	f = &FormTemplate{}
	f.fieldTemplates = make(map[string]FieldTemplate)
	return
}

func (f *FormTemplate) AddField(fieldTemplate FieldTemplate) {
	f.fieldTemplates[fieldTemplate.GetName()] = fieldTemplate
}

func (f *FormTemplate) Parse(values url.Values) (s *Form) {
	s = new(Form)
	s.Fields = make(map[string]*Field, len(f.fieldTemplates))

	for name, field := range f.fieldTemplates {
		var sf Field

		if fieldValues, ok := values[name]; ok {
			unparsed := fieldValues[len(fieldValues)-1]
			sf.Unparsed = unparsed

			parsed, err := field.Parse(unparsed)
			sf.Parsed = parsed
			sf.Error = err
		}

		s.Fields[name] = &sf
	}

	return
}

func (f *FormTemplate) New() (s *Form) {
	s = new(Form)
	s.Fields = make(map[string]*Field, len(f.fieldTemplates))

	for name, field := range f.fieldTemplates {
		var sf Field
		sf.Parsed, _ = field.Parse("")
		s.Fields[name] = &sf
	}

	return
}

func (f *FormTemplate) Validate(s *Form) {
	for name, fieldTemplate := range f.fieldTemplates {
		s.Fields[name].Error = fieldTemplate.Validate(s.Fields[name].Parsed)
	}

	if f.CustomValidate != nil {
		f.CustomValidate(s)
	}
}

type Field struct {
	Unparsed string
	Parsed   interface{}
	Error    error
}

type Form struct {
	Fields map[string]*Field
}

func (f *Form) IsValid() bool {
	for _, field := range f.Fields {
		if field.Error != nil {
			return false
		}
	}
	return true
}
