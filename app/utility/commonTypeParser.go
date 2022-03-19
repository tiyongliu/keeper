package utility

import (
	"regexp"
)

func IsTypeNumeric(dataType string) bool {
	return regexp.MustCompile("(?i)numeric|decimal").MatchString(dataType)
}

func IsTypeString(dataType string) bool {
	return regexp.MustCompile("(?i)char|binary").MatchString(dataType)
}
