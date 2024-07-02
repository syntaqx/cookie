package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/syntaqx/cookie"
)

var signingKey = []byte("super-secret-key")

var defaultCookieOptions = cookie.Options{
	HttpOnly: true,
}

var signedCookieOptions = cookie.Options{
	HttpOnly: true,
	Secure:   true,
	Signed:   true,
}

var manager *cookie.Manager

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := manager.Get(r, "DEBUG")
	if err != nil {
		setDemoCookies(w)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "RequestCookies: %v\n", c)
}

func setDemoCookies(w http.ResponseWriter) {
	manager.Set(w, "DEBUG", "true", defaultCookieOptions)
	manager.Set(w, "THEME", "dark", defaultCookieOptions)
	manager.Set(w, "Access-Token", "token_value", signedCookieOptions)
	manager.Set(w, "User-ID", "12345", signedCookieOptions)
	manager.Set(w, "Is-Admin", "true", signedCookieOptions)
	manager.Set(w, "Permissions", "read,write,execute", signedCookieOptions)
	manager.Set(w, "Friends", "1,2,3,4,5", defaultCookieOptions)
	manager.Set(w, "Expires-At", time.Now().Add(24*time.Hour).Format(time.RFC3339), signedCookieOptions)
}

func main() {
	manager = cookie.NewManager(
		cookie.WithSigningKey(signingKey),
	)

	http.HandleFunc("/", handler)

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
