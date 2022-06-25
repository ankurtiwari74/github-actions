package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Has(field string, r *http.Request) bool {
	x := f.Get(field)
	fmt.Println(f.Get(field))
	if x == "" {
		f.Errors.Add(field, "This is required field!")
		return false
	}
	return true
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty!!!")
		}
	}
}

func (f *Form) Minlength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsNumeric(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if !govalidator.IsInt(x) {
		f.Errors.Add(field, fmt.Sprintf("%q doesn't looks like a valid phone number.\n", x))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if !govalidator.IsEmail(x) {
		f.Errors.Add(field, fmt.Sprintf("%q is not a valid email.\n", x))
		return false
	}
	return true
}
