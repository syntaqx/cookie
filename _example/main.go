package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/syntaqx/cookie"
)

var defaultCookieOptions = cookie.Options{
	HttpOnly: true,
}

var signedCookieOptions = cookie.Options{
	HttpOnly: true,
	Secure:   true,
	Signed:   true,
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := cookie.Get(r, "DEBUG")
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
		NotExists   string    `cookie:"Does-Not-Exist,omitempty"`
	}

	var c RequestCookies
	if err := cookie.PopulateFromCookies(r, &c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "RequestCookies: %v\n", c)
}

func setDemoCookies(w http.ResponseWriter) {
	cookie.Set(w, "DEBUG", "true", defaultCookieOptions)
	cookie.Set(w, "THEME", "dark", defaultCookieOptions)
	cookie.Set(w, "Access-Token", "token_value", signedCookieOptions)
	cookie.Set(w, "User-ID", "12345", signedCookieOptions)
	cookie.Set(w, "Is-Admin", "true", signedCookieOptions)
	cookie.Set(w, "Permissions", "read,write,execute", signedCookieOptions)
	cookie.Set(w, "Friends", "1,2,3,4,5", defaultCookieOptions)
	cookie.Set(w, "Expires-At", time.Now().Add(24*time.Hour).Format(time.RFC3339), signedCookieOptions)
}

func main() {
	// Create a new cookie manager with a signing key.
	manager := cookie.NewManager(
		cookie.WithSigningKey([]byte("super-secret-key")),
	)

	// Set the default manager to the one we just created. This allows us to use
	// the default package functions without having to pass the manager.
	//
	// This is optional, as you can create a new manager and pass it through to
	// the functions that require it, potentially allowing you to have different
	// managers with different options.
	cookie.DefaultManager = manager

	http.HandleFunc("/", handler)

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
