package middleware

import (
	"encoding/json"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"log"
	"net/http"
	"paotui.sg/app/handler/auth"
	"paotui.sg/app/jwt"
	"strings"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		fmt.Printf("request URI:%v,request method:%v\n", r.RequestURI, r.Method)
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
		} else {
			if !strings.Contains(r.RequestURI, "auth") && !strings.Contains(r.RequestURI, "users/user") {
				token := r.Header.Get("Authorization")
				if strings.TrimSpace(token) != "" {
					var object *jose.JSONWebEncryption
					var decrypted []byte
					var userClaim auth.UserClaim
					var err error
					fmt.Printf("token->%v\n", token)
					privateKey := jwt.SignKey
					object, err = jose.ParseEncrypted(token)
					if err != nil {
						log.Println(err)
						goto Label0
					}
					decrypted, err = object.Decrypt(privateKey)
					if err != nil {
						log.Println(err)
						goto Label0
					}
					err = json.Unmarshal(decrypted, &userClaim)
					if err != nil {
						log.Println(err)
						goto Label0
					}
					fmt.Printf("user Claim is %v \n", userClaim)
					if time.Now().After(time.Unix(userClaim.ExpireDateTime, 0)) {
						goto Label0
					}
					next.ServeHTTP(w, r)
					return
				Label0:
					http.Error(w, "Forbidden", http.StatusForbidden)
				}
			} else {
				next.ServeHTTP(w, r)
			}
		}

	})
}
