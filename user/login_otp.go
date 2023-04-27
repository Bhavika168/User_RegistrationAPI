package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func rateLimit(key string, limit int64, duration time.Duration) bool {
	count, err := rdb.Incr(context.Background(), key).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	if count == 1 {
		err := rdb.Expire(context.Background(), key, duration).Err()
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	if count > limit {
		return false
	}
	return true
}

func LoginOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			data := Data{Status: "Unsuccessful", Message: "Wrong OTP."}
			jsonStr, _ := json.Marshal(data)
			w.Write(jsonStr)
		}
	}()

	var userlogin LoginOtpDto
	json.NewDecoder(r.Body).Decode(&userlogin)
	username := userlogin.Username

	if !rateLimit(username, 3, 1*time.Minute) {
		w.WriteHeader(http.StatusTooManyRequests)
		data := Data{Status: "Unsuccessful", Message: "Too many requests."}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)
		return
	}

	nonce := username + "_" + userlogin.Nonce
	otp := GetRedisNonce(nonce)

	if userlogin.UserOTP == otp {
		pattern := "session_" + username + "*"
		keys, _ := rdb.Keys(context.Background(), pattern).Result()

		if len(keys) > 0 {
			rdb.Del(context.Background(), keys...)
		}

		key := SetSessionKey(username)
		data1 := Data{Status: "Successful", Message: key}
		jsonStr, _ := json.Marshal(data1)
		w.Write(jsonStr)

	} else {
		panic("Wrong OTP.")
	}

}
