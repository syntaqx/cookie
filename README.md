# cookie

[![Go](https://github.com/syntaqx/cookie/actions/workflows/go.yml/badge.svg)](https://github.com/syntaqx/cookie/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/syntaqx/cookie.svg)](https://pkg.go.dev/github.com/syntaqx/cookie)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/cookie)](https://goreportcard.com/report/github.com/syntaqx/cookie)
[![codecov](https://codecov.io/gh/syntaqx/cookie/graph/badge.svg?token=2YEeUinfQe)](https://codecov.io/gh/syntaqx/cookie)

![Social Preview](./.github/repository-open-graph-template.png)

Cookies, but with structs, for happiness.

## Overview

`cookie` is a Go package designed to make handling HTTP cookies simple and robust by using structs. It supports standard data types, custom data types, and signed cookies to ensure data integrity.

## Features

- **Struct-based cookie management**: Easily set and get cookies using Go structs.
- **Signed cookies**: Ensure the integrity of your cookies with HMAC signatures.
- **Custom type support**: Extend cookie handling with your own data types.
- **Easy to use**: Simple API for managing cookies in your web applications.

## Installation

```bash
go get github.com/syntaqx/cookie
```

## Usage

By default, you can simply plug and play the `cookie` package into your existing
applications by using the `DefaultManager`:

You can retrieve values manually:

```go
func main() {
  cookie.Get(r, "DEBUG")
  cookie.GetSigned(r, "Access-Token")
  cookie.Set(w, "DEBUG", "true", cookie.Options{})
  cookie.Set(w, "Access-Token", "token_value", cookie.Options{Signed: true})
  cookie.SetSigned(w, "Access-Token", "token_value")
}
```

Or Populate a struct:

```go
type RequestCookies struct {
  Theme       string    `cookie:"THEME"`
  Debug       bool      `cookie:"DEBUG,unsigned"`
  AccessToken string    `cookie:"Access-Token,signed"`
}

var c RequestCookies
cookie.PopulateFromCookies(r, &c)
```

In order to sign cookies, you must provide a signing key:

```go
signingKey := []byte("super-secret-key")
cookie.DefaultManager = cookie.NewManager(
  cookie.WithSigningKey(signingKey),
)
```

## Manager

For more advanced usage, you can create a `Manager` to handle your cookies,
rather than relying on the `DefaultManager`:

```go
manager := cookie.NewManager()
```

You can optionally provide a signing key for signed cookies:

```go
signingKey := []byte("super-secret-key")
manager := cookie.NewManager(
  cookie.WithSigningKey(signingKey),
)
```

> [!TIP]
> Cookies are stored in plaintext by default (unsigned). A signed cookie is used
> to ensure the cookie value has not been tampered with. This is done by
> creating a [HMAC][] signature of the cookie value using a secret key. Then,
> when the cookie is read, the signature is verified to ensure the cookie value
> has not been modified.
>
> It is still recommended that sensitive data not be stored in cookies, and that
> HTTPS be used to prevent cookie [replay attacks][].

[HMAC]: https://en.wikipedia.org/wiki/HMAC
[replay attacks]: https://en.wikipedia.org/wiki/Replay_attack

### Setting Cookies

Use the `Set` method to set cookies. You can specify options such as path,
domain, expiration, and whether the cookie should be signed.

```go
err := manager.Set(w, "DEBUG", "true", cookie.Options{})
err := manager.Set(w, "Access-Token", "token_value", cookie.Options{Signed: true})
```

### Getting Cookies

Use the Get method to retrieve unsigned cookies and GetSigned for signed cookies.

```go
value, err := manager.Get(r, "DEBUG")
value, err := manager.GetSigned(r, "Access-Token")
```

### Populating Structs from Cookies

Use `PopulateFromCookies` to populate a struct with cookie values. The struct
fields should be tagged with the cookie names.

```go
type RequestCookies struct {
  Theme       string    `cookie:"THEME"`
  Debug       bool      `cookie:"DEBUG,unsigned"`
  AccessToken string    `cookie:"Access-Token,signed"`
}

var c RequestCookies
err := manager.PopulateFromCookies(r, &c)
```

### Supporting Custom Types

To support custom types, register a custom handler with the Manager.

```go
import (
  "reflect"
  "github.com/gofrs/uuid/v5"
  "github.com/syntaqx/manager"
)

...

manager := cookie.NewManager(
  cookie.WithSigningKey(signingKey),
  cookie.WithCustomHandler(reflect.TypeOf(uuid.UUID{}), func(value string) (interface{}, error) {
    return uuid.FromString(value)
  }),
)
```
