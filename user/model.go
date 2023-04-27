package user

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client

type User struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
	Status   string `json:"status"`
}

type SignupDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginOtpDto struct {
	Username string `json:"username"`
	Nonce    string `json:"nonce"`
	UserOTP  string `json:"userotp"`
}

type SignupOtpDto struct {
	Username string `json:"username"`
	Nonce    string `json:"nonce"`
	UserOTP  string `json:"userotp"`
}

type Data struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Nonce   string `json:"nonce"`
}

func InitialiseDb() *gorm.DB {
	dsn := "host=172.23.208.1 user=postgres password=123456 dbname=first_db port=5432 sslmode=disable"
	db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fmt.Println(dsn)
	return db
}

func InitialiseRedis() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
