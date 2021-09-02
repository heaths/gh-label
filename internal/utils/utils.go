package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ColorE(s string) (string, error) {
	if match, _ := regexp.MatchString("^#?[A-Fa-f0-9]{6}$", s); !match {
		return "", fmt.Errorf(`colors must include 6 hexadecimal digits for RGB with optional "#" prefix`)
	}

	if strings.HasPrefix(s, "#") {
		return s[1:], nil
	}

	return s, nil
}
