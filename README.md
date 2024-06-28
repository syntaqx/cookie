# cookie

[![Go Reference](https://pkg.go.dev/badge/github.com/syntaqx/cookie.svg)](https://pkg.go.dev/github.com/syntaqx/cookie)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/cookie)](https://goreportcard.com/report/github.com/syntaqx/cookie)
[![codecov](https://codecov.io/gh/syntaqx/cookie/graph/badge.svg?token=2YEeUinfQe)](https://codecov.io/gh/syntaqx/cookie)

Cookies, but with structs, for happiness.

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
