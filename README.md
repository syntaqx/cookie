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

type MyCookies struct {
  Debug bool `cookie:"DEBUG"`
}
...

var cookies Cookies
err := cookie.PopulateFromCookies(r, &cookies)
if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
}

fmt.Println(cookies.Debug)
```

## Helper Methods

### Get

For when you just want the value of the cookie:

```go
debug, err := cookie.Get(r, "DEBUG")
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

cookie.Set(w, "debug", "true", options)
cookie.Set(w, "theme", "default", options)
```

### Remove

```go
cookie.Remove(w, "debug")
```

## Signed Cookies

By default, cookies are not signed. If you want to make sure that your cookies
are signed, you can pass the `signed` tag to the struct field:

```go
type User struct {
  ID uuid.UUID `cookie:"user_id,signed"`
}
```

### `SetSigned`

If you want to set a signed cookie, you can use the `SetSigned` helper method:

```go
cookie.SetSigned(w, "user_id", "123")
```

Alternatively, you can pass `Signed` to the options when setting a cookie:

```go
cookie.Set(w, "user_id", "123", &cookie.Options{
  Signed: true,
})
```

These are functionally identical.

### `GetSigned`

If you want to get a signed cookie, you can use the `GetSigned` helper method:

```go
userID, err := cookie.GetSigned(r, "user_id")
if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
}
```

### Signing Key

By default, the signing key is set to `[]byte(cookie.DefaultSigningKey)`. You
should change this signing key for your application by assigning the
`cookie.SigningKey` variable to a secret value of your own:

```go
cookie.SigningKey = []byte("my-secret-key")
```

## Default Options

You can set default options for all cookies by assigning the
`cookie.DefaultOptions` variable:

```go
cookie.DefaultOptions = &cookie.Options{
  Domain: "example.com",
  Expires: time.Now().Add(24 * time.Hour),
  MaxAge: 86400,
  Secure: true,
  HttpOnly: true,
  SameSite: http.SameSiteStrictMode,
}
```

These options will be used as the defaults for cookies that do not strictly
override them, allowing you to only set the values you care about.

### Signed by Default

If you want all cookies to be signed by default, you can set the `Signed` field
in the `cookie.DefaultOptions`:

```go
cookie.DefaultOptions = &cookie.Options{
  Signed: true,
}
```

Which makes all cookies signed by default.

If you have any unsigned cookies, you can still access their values by using the
`unsigned` tag in the struct field:

```go
type User struct {
  Debug bool `cookie:"user_id,unsigned"`
}
```

However you will need to explicitly override this value when setting the cookie,
as the default will be to sign the cookie:

```go
cookie.Set(w, "debug", "true", &cookie.Options{
  Signed: false,
})
```

Due to the default value now overriding the option.
