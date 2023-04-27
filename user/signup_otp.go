package user

import (
	"encoding/json"
	"net/http"
	"time"
)

func SignupOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			data := Data{Status: "Unsuccessful", Message: "Wrong OTP."}
			jsonStr, _ := json.Marshal(data)
			w.Write(jsonStr)
		}
	}()

	var user SignupOtpDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !rateLimit(user.Username, 3, 1*time.Minute) {
		w.WriteHeader(http.StatusTooManyRequests)
		data := Data{Status: "Unsuccessful", Message: "Too many requests."}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)
		return
	}

	nonce := user.Username + "_" + user.Nonce
	otp := GetRedisNonce(nonce)
	if user.UserOTP == otp {

		if err := db.Model(&User{}).Where("username = ?", user.Username).Update("status", "Verified").Error; err != nil {
			panic(err)
		}

		key := SetSessionKey(user.Username)
		data1 := Data{Status: "Successful", Message: key}
		jsonStr, _ := json.Marshal(data1)
		w.Write(jsonStr)
	} else {
		panic("Wrong OTP.")
	}
}
