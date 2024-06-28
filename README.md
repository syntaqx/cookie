# cookie

[![Go Reference](https://pkg.go.dev/badge/github.com/syntaqx/cookie.svg)](https://pkg.go.dev/github.com/syntaqx/cookie)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/cookie)](https://goreportcard.com/report/github.com/syntaqx/cookie)
[![codecov](https://codecov.io/gh/syntaqx/cookie/graph/badge.svg?token=2YEeUinfQe)](https://codecov.io/gh/syntaqx/cookie)

![Social Preview](./.github/repository-open-graph-template.png)

Cookies, but with structs, for happiness.

## Usage

```go
import (
  "github.com/syntaqx/cookie"
)

...

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

## Helper Methods

### Get

For when you just want the value of the cookie:

```go
userID, err := cookie.Get(r, "user_id")
if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
}
```

### Set

While it's very easy to set Cookies in Go, often times you'll be setting
multiple cookies with the same options:

```go
options := &cookie.Options{
  Domain: "example.com",
  Expires: time.Now().Add(24 * time.Hour),
  MaxAge: 86400,
  Secure: true,
  HttpOnly: true,
  SameSite: http.SameSiteStrictMode,
}

cookie.Set(w, "user_id", "123", options)
cookie.Set(w, "user_name", "syntaqx", options)
```

### Remove

```go
cookie.Remove(w, "user_id")
```
