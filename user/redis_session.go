package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func SetSessionKey(name string) string {
	ctx := context.Background()
	sessionId := uuid.New().String()

	sessionKey := "session_" + name + "_:" + sessionId
	sessiondata := name

	err := rdb.Set(ctx, sessionKey, sessiondata, time.Hour).Err()
	if err != nil {
		panic(err)
	}
	return sessionKey
}

func GetSessionKey(sessionkey string) string {

	ctx := context.Background()
	name, err := rdb.Get(ctx, sessionkey).Result()
	if err != nil {
		panic(err)
	}
	return name
}
