package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func main() {
	h := sha256.New()
	h.Write([]byte("hello world"))
	textInBase64 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(textInBase64)
}
