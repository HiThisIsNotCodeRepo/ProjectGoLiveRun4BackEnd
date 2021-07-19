package avatar

import (
	"context"
	"encoding/json"
	"fmt"
	"paotui.sg/cloudinary"
	"paotui.sg/cloudinary/api/uploader"

	"log"
	"mime/multipart"
	"net/http"
	"os"
	"paotui.sg/app/db"
	"paotui.sg/app/handler/error_util"
	"strings"
)

type NewAvatarResponse struct {
	Status    string `json:"status"`
	Msg       string `json:"msg"`
	AvatarUrl string `json:"avatarUrl"`
}

func NewAvatar(w http.ResponseWriter, r *http.Request) {
	defer error_util.ErrorHandle(w)
	if r.Method == http.MethodOptions {
		return
	}
	var newAvatarResponse NewAvatarResponse
	var userID = r.URL.Query().Get("userID")
	var err error
	var uploadFile multipart.File
	var cld *cloudinary.Cloudinary
	var ctx context.Context
	var uploadResult *uploader.UploadResult
	APISecret := strings.TrimSpace(os.Getenv("APISECRET"))
	fmt.Printf("newAvatar->request URI:%v\n", r.RequestURI)
	encoder := json.NewEncoder(w)
	if strings.TrimSpace(userID) == "" {
		newAvatarResponse.Status = "error"
		newAvatarResponse.Msg = "no user id"
		goto Label1
	}
	uploadFile, _, err = r.FormFile("avatar")
	if err != nil {
		log.Println(err)
		newAvatarResponse.Status = "error"
		newAvatarResponse.Msg = "no avatar fail"
		goto Label1
	}
	cld, err = cloudinary.NewFromParams("paotui", "354942911769922", APISecret)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	ctx = context.Background()
	uploadResult, err = cld.Upload.Upload(
		ctx,
		uploadFile,
		uploader.UploadParams{PublicID: userID})
	if err != nil {
		log.Printf("upload err:%v\n", err)
		goto Label0
	}

	fmt.Printf("new avatar upload url:%v\n", uploadResult.SecureURL)
	_, err = db.Db.Exec("UPDATE user SET avatar_url = ? WHERE uid =?", uploadResult.SecureURL, userID)
	if err != nil {
		log.Println(err)
		goto Label0
	}
	newAvatarResponse.Status = "success"
	newAvatarResponse.Msg = "new avatar added success"
	newAvatarResponse.AvatarUrl = uploadResult.SecureURL

Label0:
	if newAvatarResponse.Status != "success" {
		newAvatarResponse.Status = "error"
		newAvatarResponse.Msg = "server error"
	}
Label1:
	encodeErr := encoder.Encode(newAvatarResponse)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
