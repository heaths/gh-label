package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

func RandomColor() string {
	r := rand.Int31n(256)
	g := rand.Int31n(256)
	b := rand.Int31n(256)

	return fmt.Sprintf("%02X%02X%02X", r, g, b)
}

func ValidateColor(s string) (string, error) {
	if match, _ := regexp.MatchString("^#?[A-Fa-f0-9]{6}$", s); !match {
		return "", fmt.Errorf(`colors must include 6 hexadecimal digits for RGB with optional "#" prefix`)
	}

	if strings.HasPrefix(s, "#") {
		return s[1:], nil
	}

	return s, nil
}
