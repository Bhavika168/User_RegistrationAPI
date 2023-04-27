package user

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"image/png"
	"net/http"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var newUser SignupDto
	json.NewDecoder(r.Body).Decode(&newUser)
	defer func() {
		if err := recover(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			data := Data{Status: "Unsuccessful", Message: "User already exist."}
			jsonStr, _ := json.Marshal(data)
			w.Write(jsonStr)
		}
	}()
	if isNewUser(newUser.Username) {
		totpCode, base64Str := GenerateTotpcodeAndQR(newUser.Username)

		user := CreateUser(newUser.Username, newUser.Password, totpCode)

		nonce := SetRedisNonce(user.Username)

		message := Message{Status: "Successful", Message: base64Str, Nonce: nonce}
		jsonStr, _ := json.Marshal(message)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)

	} else {
		panic("User already exist.")
	}
}

func isNewUser(username string) bool {

	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return true // user does not exist
	} else {
		return false //user exist
	}
}

func GenerateTotpcodeAndQR(username string) (string, string) {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      username,
		AccountName: "qr",
	})
	img, _ := key.Image(256, 256)

	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	return key.Secret(), base64Str
}

func CreateUser(username, password, totpcode string) User {
	hpassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	db.AutoMigrate(&User{})

	user := User{Username: username, Password: string(hpassword), Secret: totpcode, Status: "New User"}
	if err := db.WithContext(context.Background()).Create(&user).Error; err != nil {
		panic("failed to create user")
	}
	return user
}
