package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/syntaqx/cookie"
)

type RequestCookies struct {
	ApplicationID uuid.UUID `cookie:"Application-ID"`
	Theme         string    `cookie:"THEME"`
	Debug         bool      `cookie:"DEBUG"`
	AccessToken   string    `cookie:"Access-Token,signed"`
	UserID        int       `cookie:"User-ID,signed"`
	IsAdmin       bool      `cookie:"Is-Admin,signed"`
	Permissions   []string  `cookie:"Permissions,signed"`
	ExpiresAt     time.Time `cookie:"Expires-At,signed"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// If none of the cookies are set, we'll set them and refresh the page
	// so the rest of the demo functions.
	_, err := cookie.Get(r, "Application-ID")
	if err != nil {
		setDemoCookies(w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Populate struct from cookies
	var req RequestCookies
	err = cookie.PopulateFromCookies(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Dump the struct as a response
	fmt.Fprintf(w, "RequestCookies: %+v", req)
}

func setDemoCookies(w http.ResponseWriter) {
	// Set cookies
	cookie.Set(w, "Application-ID", uuid.Must(uuid.NewV7()).String(), nil)
	cookie.Set(w, "THEME", "default", nil)
	cookie.Set(w, "DEBUG", "true", nil)

	secureOptions := &cookie.Options{
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}

	// Set signed cookies
	cookie.SetSigned(w, "Access-Token", "some-access-token", secureOptions)
	cookie.SetSigned(w, "User-ID", "123", secureOptions)
	cookie.SetSigned(w, "Is-Admin", "true", secureOptions)
	cookie.SetSigned(w, "Permissions", "read,write,execute", secureOptions)
	cookie.SetSigned(w, "Expires-At", time.Now().Add(24*time.Hour).Format(time.RFC3339), secureOptions)
}

func main() {
	cookie.DefaultOptions = &cookie.Options{
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	http.HandleFunc("/", handler)

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
