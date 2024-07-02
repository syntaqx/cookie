package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/syntaqx/cookie"
)

func main() {
	signingKey := []byte("super-secret-key")

	manager := cookie.NewManager(
		cookie.WithSigningKey(signingKey),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := manager.Get(r, "DEBUG")
		if err != nil {
			manager.Set(w, "DEBUG", "true", cookie.Options{})
			manager.Set(w, "THEME", "dark", cookie.Options{})
			manager.Set(w, "Access-Token", "token_value", cookie.Options{Signed: true})
			manager.Set(w, "User-ID", "12345", cookie.Options{Signed: true})
			manager.Set(w, "Is-Admin", "true", cookie.Options{Signed: true})
			manager.Set(w, "Permissions", "read,write,execute", cookie.Options{Signed: true})
			manager.Set(w, "Friends", "1,2,3,4,5", cookie.Options{})
			manager.Set(w, "Expires-At", time.Now().Add(24*time.Hour).Format(time.RFC3339), cookie.Options{Signed: true})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		type RequestCookies struct {
			Theme       string    `cookie:"THEME"`
			Debug       bool      `cookie:"DEBUG,unsigned"`
			AccessToken string    `cookie:"Access-Token,signed"`
			UserID      int       `cookie:"User-ID,signed"`
			IsAdmin     bool      `cookie:"Is-Admin,signed"`
			Permissions []string  `cookie:"Permissions,signed"`
			Friends     []int     `cookie:"Friends,unsigned"`
			ExpiresAt   time.Time `cookie:"Expires-At,signed"`
		}

		var c RequestCookies
		if err := manager.PopulateFromCookies(r, &c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("Cookies: %+v", c)))
	})

	fmt.Printf("Listening on http://localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}
