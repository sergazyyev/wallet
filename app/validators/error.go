package validators

import "strings"

//ValidationError for validation error messages
type ValidationError struct {
	messages map[string][]string
}

//Error implement error response
func (v *ValidationError) Error() string {
	var res string
	for _, m := range v.messages {
		res = res + strings.Join(m, ",")
	}
	return res
}

//NewValidationError constructor
func NewValidationError(messages map[string][]string) error {
	return &ValidationError{
		messages: messages,
	}
}

//Messages return messages
func (v *ValidationError) Messages() map[string][]string {
	return v.messages
}
