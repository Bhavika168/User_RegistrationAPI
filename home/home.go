package home

import (
	"Project18/user"
	"encoding/json"
	"net/http"
)

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				data := Message{Status: "Unsuccessful", Message: "Wrong Session Token."}
				jsonStr, _ := json.Marshal(data)
				w.Write(jsonStr)
			}
		}()

		username := user.GetSessionKey(r.Header.Get("X-Session"))

		if username != "" {

		} else {
			panic("Wrong Header")
		}
		next.ServeHTTP(w, r)
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := Message{Status: "Successful", Message: "Inside Home."}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Write(jsonStr)
}
