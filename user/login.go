package user

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user LoginDto
	json.NewDecoder(r.Body).Decode(&user)

	if isValidUser(user.Username, user.Password) {

		nonce := SetRedisNonce(user.Username)
		message := Message{Status: "Successful", Message: "Enter Your OTP", Nonce: nonce}
		jsonStr, _ := json.Marshal(message)
		w.Write(jsonStr)

	} else {
		w.WriteHeader(http.StatusNotFound)
		data := Data{Status: "Unsuccessful", Message: "User not found."}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)
		return
	}
}

func isValidUser(username, password string) bool {

	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return false
	} else if user.Status != "Verified" {
		return false
	} else {
		storedHash := user.Password
		plaintextPassword := password
		err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(plaintextPassword))
		if err != nil {
			return false
		}
	}
	return true
}
