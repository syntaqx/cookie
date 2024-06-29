package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/syntaqx/cookie"
)

type AccessTokenRequest struct {
	ApplicationID uuid.UUID `cookie:"Application-ID"`
	AccessToken   string    `cookie:"Access-Token"`
	UserID        int       `cookie:"User-ID"`
	IsAdmin       bool      `cookie:"Is-Admin"`
	Permissions   []string  `cookie:"Permissions"`
	ExpiresAt     time.Time `cookie:"Expires-At"`
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
	var req AccessTokenRequest
	err = cookie.PopulateFromCookies(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Dump the struct as a response
	fmt.Fprintf(w, "AccessTokenRequest: %+v", req)
}

func setDemoCookies(w http.ResponseWriter) {
	// Set cookies
	cookie.Set(w, "Application-ID", uuid.Must(uuid.NewV7()).String(), &http.Cookie{
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	cookie.Set(w, "Access-Token", "some-access-token", &http.Cookie{
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	cookie.Set(w, "User-ID", "123", &http.Cookie{
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	cookie.Set(w, "Is-Admin", "true", &http.Cookie{
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	cookie.Set(w, "Permissions", "read,write,execute", &http.Cookie{
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	cookie.Set(w, "Expires-At", time.Now().Add(24*time.Hour).Format(time.RFC3339), &http.Cookie{
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
