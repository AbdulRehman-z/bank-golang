package gapi

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func validateString(value string, minLen int8, maxLen int8) error {
	x := len(value)
	if x < int(minLen) || x > int(maxLen) {
		return fmt.Errorf("string length must be between %d and %d", minLen, maxLen)
	}
	return nil

}

func ValidateUsername(value string) error {
	if err := validateString(value, 3, 32); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("username can only contain lowercase letters, numbers and underscores")
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := validateString(value, 3, 50); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return fmt.Errorf("full name can only contain letters and spaces")
	}

	return nil
}

func ValidatePassword(value string) error {
	return validateString(value, 8, 20)
}

func ValidateEmail(value string) error {
	if err := validateString(value, 5, 50); err != nil {
		return err
	}

	err, _ := mail.ParseAddress(value)
	if err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}
