package utils

import (
  "crypto/aes"
  "crypto/cipher"
  "encoding/base64"
  "os"
)

func Encode(b []byte) string {
  return base64.StdEncoding.EncodeToString(b)
}
func Decode(s string) []byte {
  data, err := base64.StdEncoding.DecodeString(s)
  if err != nil {
    panic(err)
  }
  return data
}
func Encrypt(text string) (string, error) {
  ebstr := os.Getenv("EBSTR")
  bytes := []byte(ebstr)
  secret := os.Getenv("ESECRET")
  block, err := aes.NewCipher([]byte(secret))
  if err != nil {
    return "", err
  }
  plainText := []byte(text)
  cfb := cipher.NewCFBEncrypter(block, bytes)
  cipherText := make([]byte, len(plainText))
  cfb.XORKeyStream(cipherText, plainText)
  return Encode(cipherText), nil
}

func Decrypt(text string) (string, error) {
  ebstr := os.Getenv("EBSTR")
  bytes := []byte(ebstr)
  secret := os.Getenv("ESECRET")
  block, err := aes.NewCipher([]byte(secret))
  if err != nil {
    return "", err
  }
  cipherText := Decode(text)
  cfb := cipher.NewCFBDecrypter(block, bytes)
  plainText := make([]byte, len(cipherText))
  cfb.XORKeyStream(plainText, cipherText)
  return string(plainText), nil
}
