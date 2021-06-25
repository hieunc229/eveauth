# eveauth

Auth middleware and request handler wrapper for GO. It used [boltDB](https://github.com/boltdb/bolt) to store user data

What's in this guide:

1. Install eveauth
2. Auth wrapper and Middleware
3. Register handler
4. Login handler (return jwt token)
5. How to use JWT token
6. How to verify a request (*http.Request)


### 1. Install

Install from terminal:

```ssh
$ go get github.com/hieunc229/eveauth
```

After installed, when using `eveauth`, it usually imports the module automatically. Otherwise, you can add package `github.com/hieunc229/eveauth` to the import command

### 2. Auth Handle Wrapper and Middleware

Use `eveauth.AuthMiddleware` or `eveauth.AuthHandler` to authorize users.

```go
import (
    github.com/hieunc229/eveauth
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

```json
{
    "ok": false,
    "error": "invalid access"
}
```

### 3. Register handler

Use `eveauth.Register` handler to handle create account

```go
router.HandlerFunc("/auth/register", eveauth.Register)
```

The body json data must be:
```json
{
    "data": {
        "username": "xxxxxx",
        "password": "xxxxxx"
    }
}
```

Success return:
```json
{
    "ok": true
}
```

Error return:
```json
{
    "error": "error message",
    "ok": false
}
```

### 4. Login handler

Use `eveauth.Login` handler to handle login

```go
router.HandlerFunc("/auth/login", eveauth.Login)
```

The body json data must be:
```json
{
    "data": {
        "username": "xxxxxx",
        "password": "xxxxxx"
    }
}
```

Success return:
```json
{
    "data": {
        "token": "jwt token str" // use as bearer token
    },
    "ok": true
}
```

Error return:
```json
{
    "error": "error message",
    "ok": false
}
```

### 5. How to use JWT token

After send a login request and receive a sucess response, you'll be given a `token`. This token is meant to use as [Bearer](https://swagger.io/docs/specification/authentication/bearer-authentication/) token.

Whenever you make a request and want it to be authorized, added `Authorization: Bearer <token>` to the request header

Here is an example with [fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch) in JS
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

### 6. How to verify a request contains a JWT token

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
