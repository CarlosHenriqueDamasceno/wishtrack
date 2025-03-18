package validation

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ErrorCollection map[string][]*ValidationError

func (c *ErrorCollection) Append(err *ValidationError) {
	if *c == nil {
		*c = make(ErrorCollection)
	}
	(*c)[err.Field] = append((*c)[err.Field], err)
}

func (c *ErrorCollection) AppendCollection(collection *ErrorCollection) {
	if *c == nil {
		*c = *collection
		return
	}

	for _, errors := range *collection {
		for _, err := range errors {
			c.Append(err)
		}
	}
}

func (c *ErrorCollection) WithMessage(field, message string) {
	c.Append(&ValidationError{
		Field:   field,
		Message: message,
	})
}

// Returns if a error collection has some error attached
// Should be used instead of nil comparison
func (c ErrorCollection) HasError() bool {
	return len(c) > 0
}

func (c ErrorCollection) Get(field string) []*ValidationError {
	return c[field]
}

func (c ErrorCollection) Error() string {
	var msg string
	for _, field := range c {
		for _, err := range field {
			msg += err.Error() + "; "
		}
	}
	return msg
}

func (c ErrorCollection) MarshalJSON() ([]byte, error) {
	var errors = make(map[string][]string)
	for field, list := range c {
		for _, err := range list {
			errors[field] = append(errors[field], err.Message)
		}
	}

	return json.Marshal(struct {
		Errors map[string][]string `json:"errors"`
	}{errors})
}
