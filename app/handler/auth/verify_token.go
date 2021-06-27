package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"paotui.sg/app/jwt"
	"strings"
	"time"
)

type VerifyTokenRequest struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

type VerifyTokenResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var verifyTokenRequest VerifyTokenRequest
	var verifyTokenResponse VerifyTokenResponse
	var err error
	var privateKey *rsa.PrivateKey
	var object *jose.JSONWebEncryption
	var decrypted []byte
	var claimUser UserClaim
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	fmt.Printf("request URI:%v\n", r.RequestURI)

	userID := mux.Vars(r)["userID"]
	err = decoder.Decode(&verifyTokenRequest)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	fmt.Printf("verifyRequest:%v\n", verifyTokenRequest)
	if strings.TrimSpace(userID) == "" {
		verifyTokenResponse.Status = "error"
		verifyTokenResponse.Msg = "no userID"
		goto Label1
	}

	privateKey = jwt.SignKey
	object, err = jose.ParseEncrypted(verifyTokenRequest.Token)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	decrypted, err = object.Decrypt(privateKey)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	err = json.Unmarshal(decrypted, &claimUser)
	if err != nil {
		fmt.Println(err)
		goto Label0
	}
	fmt.Printf("userID:%v ,Claim user is %v \n", userID, claimUser)
	if userID != claimUser.UserId {
		verifyTokenResponse.Status = "error"
		verifyTokenResponse.Msg = "userId is not correct"
		goto Label1
	}
	if time.Now().After(time.Unix(claimUser.ExpireDateTime, 0)) {
		verifyTokenResponse.Status = "error"
		verifyTokenResponse.Msg = "the token is expired"
		goto Label1
	}
	verifyTokenResponse.Status = "success"
	verifyTokenResponse.Msg = fmt.Sprintf("token is verified")

Label0:
	if verifyTokenResponse.Status != "success" {
		verifyTokenResponse.Status = "error"
		verifyTokenResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(verifyTokenResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
