package datatype

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	// regular expression to validation unix timestamp
	RegexUnixTimestamp = `^([\d]{10}|0)$`
)

// IsTimestamp will validate if the value is a valid unix timestamp or not
func IsTimestamp(value interface{}) error {
	if value == nil {
		return nil
	}

	return MatchRegex(fmt.Sprintf("%v", value), RegexUnixTimestamp)
}

// MatchRegex checks if given input matches a given regex or not
func MatchRegex(value string, regex string) error {
	if validString, err := regexp.Compile(regex); err != nil {
		return errors.New("invalid regex")
	} else if !validString.MatchString(value) {
		return errors.New("not a valid input")
	}

	return nil
}
