# eveauth

Auth middleware and request handler wrapper for GO. It uses [boltDB](https://github.com/boltdb/bolt) to store user data

Included in this guide:

1. [Usage](#1-getting-started)
- [Install eveauth](#install)
- [Auth wrapper and Middleware](#auth-handle-wrapper-and-middleware)
- [Register handler](#register-handler)
- [Login handler (return jwt token)](#login-handler)
- [How to use JWT token](#how-to-use-jwt-token)
- [How to verify a request (*http.Request)](#how-to-verify-a-request-contains-a-jwt-token)
- [How to change password](#how-to-change-password)
2. [Feedback and Contribute](#feedback-and-contribute)
3. [Licenses](#licenses)

## 1. Getting started

### Install

```ssh
$ go get github.com/hieunc229/eveauth
```

### Auth Handle Wrapper and Middleware

Use `eveauth.AuthMiddleware` or `eveauth.AuthHandler` to authorize users.

```go
import (
    "github.com/hieunc229/eveauth"
)

// It can be used for any router with standard handler
// ie. func (http.ResponseWriter, r *http.Request)
router := mux.NewRouter()

// Use as middleware
router.Use("/auth", eveauth.AuthMiddleware)

// Or wrap around a handler
router.HandlerFunc("/user_only", eveauth.AuthHandler(yourHandler))
```

When there is an anonymous request to these paths, it return the following:

```js
{
    "ok": false,
    "error": "invalid access"
}
```

### Register handler

Use `eveauth.RegisterHandler` handler to handle create account

```go
router.HandlerFunc("/auth/register", eveauth.RegisterHandler)
```

The body json data must be:
```js
{
    "data": {
        "username": "xxxxxx",
        "password": "xxxxxx"
    }
}
```

Success return:
```js
{
    "ok": true
}
```

Error return:
```js
{
    "error": "error message",
    "ok": false
}
```

### Login handler

Use `eveauth.LoginHandler` handler to handle login

```go
router.HandlerFunc("/auth/login", eveauth.LoginHandler)
```

The body json data must be:
```js
{
    "data": {
        "username": "xxxxxx",
        "password": "xxxxxx"
    }
}
```

Success return:
```js
{
    "data": {
        "token": "jwt token str" // use as bearer token
    },
    "ok": true
}
```

Error return:
```js
{
    "error": "error message",
    "ok": false
}
```

### How to use JWT token

After send a login request and receive a sucess response, you'll be given a `token`. This token is meant to use as [Bearer](https://swagger.io/docs/specification/authentication/bearer-authentication/) token.

Whenever you make a request and want it to be authorized, added `Authorization: Bearer <token>` to the request header

Here is an example with [fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch) in JavaScript
```js
fetch("/user_only/items/goodItemId", {
    method: "POST",
    headers: {
        'Authorization': 'Bearer <token>'
        // 'Content-Type': 'application/json'
        // ...
    }
    // body: ...
})
```

### How to verify a request contains a JWT token

Use `eveauth.VerifyRequest(*http.Request) (*JWTPayload, err)` to verify your http request. This handler will (1) get the bearer token, (2) check if it is valid, (3) return the data contains in the token or error

Here is an example:
```go
func yourHandler(w http.ResponseWriter, r *http.Request) {

    payload, err := eveauth.VerifyRequest(r)

    if err != nil {
        // bearer token is not valid or expired
        return;
    }

    // no err, looking good

    username = payload.Username
}
```

### How to change password

Use `eveauth.ChangePasswordhandler` to handle change password request. Note that the request must be authorized with Bearer token mention above

```go
router.HandlerFunc("/auth/change-password", eveauth.ChangePasswordhandler)
```

The body json payload must be:
```js
{
    "data": {
        "password": "oldPassword",
        "new_password": "xxxxxxx",

        //// set to `true` to replace the current token with a new one
        "change_token": false, 

        // set to `true` to remove all existing tokens, then add a new one 
        // (i.e useful for logout all other devices feature)
        "clear_tokens": false, 
    }
}
```

Here is an example using fetch in JavaScript
```js
fetch("/user_only/items/goodItemId", {
    method: "POST",
    headers: {
        'Authorization': 'Bearer <token>'
        // 'Content-Type': 'application/json'
        // ...
    }
    body: JSON.stringify({
        data: {
            password: "xxxxx",
            new_password: "xxxxxxxx"
        }
    })
})
```


## 2. Feedback and Contribute

Always welcome. Please [open a new thread](https://github.com/hieunc229/eveauth/issues/new)

## 3. Licenses

- eveauth MIT
- BoltDB MIT
