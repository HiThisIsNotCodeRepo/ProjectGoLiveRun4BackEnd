package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"paotui.sg/app/jwt"
	"time"
)

type UserClaim struct {
	UserId         string
	ExpireDateTime int64
}

func main() {
	var err error
	publicKey := jwt.VerifyKey
	encryptor, _ := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: publicKey}, nil)

	userClaim, _ := json.Marshal(&UserClaim{UserId: "abc", ExpireDateTime: time.Now().Add(1 * time.Hour).Unix()})

	object, _ := encryptor.Encrypt(userClaim)

	token, _ := object.CompactSerialize()
	fmt.Println(token)
	privateKey := jwt.SignKey
	object, err = jose.ParseEncrypted(token)
	log.Println(err)
	decrypted, err := object.Decrypt(privateKey)
	_ = json.Unmarshal(decrypted, &userClaim)

	fmt.Printf("Claim user is %s \n", userClaim)
}
