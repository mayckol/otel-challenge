package utils

import "regexp"

type ZipCode string

func (z ZipCode) IsValid() bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(string(z))
}

func (z ZipCode) Raw() string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(string(z), "")
}
