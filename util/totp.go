package util

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image"
)

func GenerateTOTP(email string) (string, *image.Image, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Nistagram",
		AccountName: email,
	})
	if err != nil {
		return "", nil, err
	}
	img, err := key.Image(200, 200)
	if err != nil {
		return "", nil, err
	}
	return key.URL(), &img, nil
}

func ValidateTOTP(totpURL string, passcode string) bool {
	key, err := otp.NewKeyFromURL(totpURL)
	if err != nil {
		return false
	}
	return totp.Validate(passcode, key.Secret())
}