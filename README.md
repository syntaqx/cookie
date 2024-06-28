# cookie

Cookies, but with structs

```go
type User struct {
	ID   uuid.UUID `cookie:"user_id"`
	Name string    `cookie:"user_name"`
}

...

var user User
err := cookie.PopulateFromCookies(r, &user)
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}

fmt.Println(user.ID, user.Name)
```
