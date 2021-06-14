package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"os"
)

type EncryptedString struct {
	Data string
}

func (es *EncryptedString) Scan(value interface{}) error {
	key := []byte(os.Getenv("DB_SEC_ENC"))
	ciphertext := value.([]byte)
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
		return err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
		return errors.New("poink")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	es.Data = string(plaintext)
	return nil
}

func (es EncryptedString) Value() (driver.Value, error) {
	text := []byte(es.Data)
	key := []byte(os.Getenv("DB_SEC_ENC"))
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return gcm.Seal(nonce, nonce, text, nil), nil
}
