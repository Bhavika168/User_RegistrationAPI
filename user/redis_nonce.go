package user

import (
	"context"
	"math/rand"
	"time"

	"github.com/pquerna/otp/totp"
)

func SetRedisNonce(username string) string {

	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randStr := make([]byte, 20)
	for i := range randStr {
		randStr[i] = charset[seededRand.Intn(len(charset))]
	}
	nonce := username + "_" + string(randStr)

	var user User
	err1 := db.Where("username = ?", username).First(&user).Error
	if err1 != nil {
		panic(err1)
	}
	otp, _ := totp.GenerateCode(user.Secret, time.Now())

	err := rdb.Set(context.Background(), nonce, otp, time.Minute).Err()
	if err != nil {
		panic(err)
	}
	return string(randStr)
}

func GetRedisNonce(nonce string) string {

	otp, err := rdb.Get(context.Background(), nonce).Result()
	if err != nil {
		panic(err)
	}
	return otp
}
