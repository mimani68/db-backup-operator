package utils

import (
	"errors"
	"regexp"
	"strings"
)

func FindKeyword(input string) (string, error) {
	re := regexp.MustCompile(`(?i)(mysql|postgresql|postgres)`)
	// return re.FindAllString(input, -1)
	for _, match := range re.FindAllString(input, -1) {
		return strings.ToUpper(match), nil
	}
	return "", errors.New("invalid regex match")
}
