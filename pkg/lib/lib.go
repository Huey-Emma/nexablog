package lib

import (
	"net/mail"
	"strings"
)

type H[T any] map[string]T

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func WhiteSpace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func NonWhiteSpace(s string) bool {
	return !WhiteSpace(s)
}
