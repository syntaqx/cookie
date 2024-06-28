package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/syntaqx/cookie"
)

type User struct {
	ID   uuid.UUID `cookie:"user_id"`
	Name string    `cookie:"user_name"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie.Set(w, "user_id", "406f82a8-59f7-4bb8-8328-54d1ce9fa709", nil)
		cookie.Set(w, "user_name", "Alice", nil)

		var user User
		err := cookie.PopulateFromCookies(r, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Hello, %s (%s)!", user.Name, user.ID)
	})

	log.Printf("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
