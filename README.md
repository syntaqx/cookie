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

By default, cookies are stored in plaintext.

Cookies can be signed to ensure their value has not been tampered with. This
works by creating a [HMAC](https://en.wikipedia.org/wiki/HMAC) of the value
(current cookie), and base64 encoding it. When the cookie gets read, it
recalculates the signature and makes sure that it matches the signature attached
to it.

It is still recommended that sensitive data not be stored in cookies, and that
HTTPS be used to prevent cookie
[replay attacks](https://en.wikipedia.org/wiki/Replay_attack).

If you want to sign your cookies, this can be accomplished by:

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

### Reading Signed Cookies

To read signed cookies into your struct, you can use the `signed` tag:

```go
type User struct {
  ID uuid.UUID `cookie:"user_id,signed"`
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

Which will now sign all cookies by default when using the `Set` method. You can
still override this by passing `Signed: false` to the options when setting a
cookie.

```go
cookie.Set(w, "debug", "true", &cookie.Options{
  Signed: false,
})
```

This will require the use of the `GetSigned` method to retrieve cookie values.

```go
debug, err := cookie.GetSigned(r, "debug")
if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
}
```

When defaulting to signed cookies, unsigned cookies can still be populated by
using the `unsigned` tag in the struct field:

```go
type MyCookies struct {
  Debug bool `cookie:"debug,unsigned"`
}
```

Or retrieved using the `Get` method, which always retrieves the plaintext value:

```go
debug, err := cookie.Get(r, "debug")
if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
}
```
