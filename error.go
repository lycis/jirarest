package jirarest

import (
	"encoding/json"
	"fmt"
)

// This type of error represents an error returned by jira
type JiraError struct {
	ErrorMessages []string          `json:"errorMessages"`
	Errors        map[string]string `json:"errors"`
}

// Provides the description of the error
func (je JiraError) Error() string {
	var msg string
	for _, e := range je.ErrorMessages {
		msg = e + ";" + msg
	}

	var errs string
	for _, e := range je.Errors {
		errs = e + ";" + errs
	}

	return fmt.Sprintf("errorMessage: %s errors: %s", msg, errs)
}

// internal: returns a jira error if the given JSON contains an error message
func toJiraError(content []byte) (bool, JiraError) {
	var je JiraError
	err := json.Unmarshal(content, &je)
	if err != nil {
		return false, je
	}

	if len(je.ErrorMessages) == 0 && len(je.Errors) == 0 {
		return false, je
	}

	return true, je
}
