package utils

import "github.com/matoous/go-nanoid/v2"

const (
	alphaLowercase   = "abcdefghijklmnopqrstuvwxyz"
	alphaUppercase   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numeric          = "0123456789"
	alphaNumericBoth = alphaLowercase + alphaUppercase + numeric
)

func NewID() string {
	return gonanoid.MustGenerate(alphaNumericBoth, 8)
}
