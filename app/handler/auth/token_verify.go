package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"paotui.sg/app/jwt"
	"time"
)

type TokenVerifyRequest struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

type TokenVerifyResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func TokenVerify(w http.ResponseWriter, r *http.Request) {
	var tokenVerifyRequest TokenVerifyRequest
	var tokenVerifyResponse TokenVerifyResponse
	var err error
	var privateKey *rsa.PrivateKey
	var object *jose.JSONWebEncryption
	var decrypted []byte
	var userClaim UserClaim
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	fmt.Printf("token verify->request URI:%v\n", r.RequestURI)

	err = decoder.Decode(&tokenVerifyRequest)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	fmt.Printf("verifyRequest:%v\n", tokenVerifyRequest)

	privateKey = jwt.SignKey
	object, err = jose.ParseEncrypted(tokenVerifyRequest.Token)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	decrypted, err = object.Decrypt(privateKey)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	err = json.Unmarshal(decrypted, &userClaim)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	fmt.Printf("user Claim is %v \n", userClaim)
	if tokenVerifyRequest.UserId != userClaim.UserId {
		tokenVerifyResponse.Status = "error"
		tokenVerifyResponse.Msg = "userId is not correct"
		goto Label1
	}
	if time.Now().After(time.Unix(userClaim.ExpireDateTime, 0)) {
		tokenVerifyResponse.Status = "error"
		tokenVerifyResponse.Msg = "the token is expired"
		goto Label1
	}
	tokenVerifyResponse.Status = "success"
	tokenVerifyResponse.Msg = fmt.Sprintf("token is verified")

Label0:
	if tokenVerifyResponse.Status != "success" {
		tokenVerifyResponse.Status = "error"
		tokenVerifyResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(tokenVerifyResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
