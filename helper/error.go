package helper

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

type GroupedError struct {
	Errors []gojsonschema.ResultError
}

func (e *GroupedError) Error() string {
	var b strings.Builder // Using strings.Builder for efficient concatenation

	for _, desc := range e.Errors {
		// Each person's information is appended as a formatted string
		b.WriteString(fmt.Sprintf("- %s\n", desc))
	}

	return b.String()
}

func (e *GroupedError) GetErrors() []gojsonschema.ResultError {
	return e.Errors
}
