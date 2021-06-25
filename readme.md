# EveAuth

Auth middleware and request handler wrapper for GO. It uses [boltDB](https://github.com/boltdb/bolt) to store user data

Included in this guide:

1. [Usage](#1-getting-started)
- [Install eveauth](#install)
- [Auth wrapper and Middleware](#auth-handle-wrapper-and-middleware)
- [Register handler](#register-handler)
- [Login handler](#login-handler)
- [Change password handler](#change-password-handler)
- [How to use JWT token](#how-to-use-jwt-token)
- [How to verify a request](#how-to-verify-a-request-contains-a-jwt-token)
- [Setup enviroment variables](#setup-enviroment-variables)
2. [Changelog](#2-changelog)
3. [Feedback and Contribute](#3-feedback-and-contribute)
4. [Licenses](#4-licenses)

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

// auth options
authOptions := eveauth.AuthHandlerOptions{
    // allow only member
    // other values can be eveauth.RoleAdmin, eveauth.RoleAnonymous
    Role: eveauth.RoleMember,

	// Set to true to forbid access from users with different RoleLevel.
	// Or set to false (or nil) to forbid only users with lower RoleLevel
    RoleExact: true
}

// Use as middleware
router.Use("/user_only", eveauth.AuthMiddleware(&authOptions))

// Or wrap around a handler
router.HandlerFunc("/user_only", eveauth.AuthHandler(yourHandler, &authOptions))
```

When there is an anonymous request to these paths, it return the following:

```js
{
    "ok": false,
    "error": "invalid access"
}
```

### Register handler

Use `eveauth.RegisterHandler` handler to handle create account. Registed users will have `eveauth.RoleMember` role

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


### Change password handler

Use `eveauth.ChangePasswordhandler` to handle change password request. Note that **the request must be authorized with Bearer token** mention above. _If you don't have a bearer token, login to get a bearer token first_

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
Success response:
```js
{
    "data" {
        // If `change_token` or `clear_tokens` is true, you will need to use this new token
        // Otherwise, this value will be an empty string ("")
        "new_token": "" 
    },
    "ok": true
}
```

Error response:
```js
{
    "error": "error message",
    "ok": false
}
```

Here is a change password example using fetch in JavaScript
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

Use `eveauth.VerifyRequest(*http.Request, *eveauth.AuthHandlerOptions) (*JWTPayload, err)` to verify your http request. 

Verify a http.Request by (1) get bearer token, (2) verify if the token is a valid jwt token, (3) get userData then check if token is still active (4) then check if the user has the proper role if authOption != nil

Here is an example:
```go
func yourHandler(w http.ResponseWriter, r *http.Request) {

    payload, err := eveauth.VerifyRequest(r, &eveauth.AuthHandlerOptions{})

    if err != nil {
        // bearer token is not valid or expired
        return;
    }

    // no err, looking good

    username = payload.Username
}
```

### Setup enviroment variables

There are a few enviroment variables that you should update when using the product:

- `EVEAUTH_JWT_SECRET` (default `eveauth`): a secret string to create jwt token
- `EVEAUTH_PATH` (default `auth`): path to your auth database. You can use absolute or relative path. If you use relative path, the path root is where you run the command)

There are 2 ways to set these enviroment variables:

1. Using flags in command (_only used after you build the application (aka binary file)_). For example:
```sh
$ ./coolapp -EVEAUTH_JWT_SECRET=randomstringnoonecanguess -EVEAUTH_PATH=/ect/safe-area/coolappAuth.db
```

2. Use [godotenv](https://github.com/joho/godotenv) (or any dotenv alternative). First, install `godotenv` (`go get https://github.com/joho/godotenv`), then create a `.env` at your root directory with the following content:

```
# .env
EVEAUTH_JWT_SECRET=randomstringnoonecanguess
EVEAUTH_PATH=/ect/safe-area/coolappAuth.db
```

3. Load the `.env` file. Read the manual from `dotenv` package you use. For example, for `godotenv`:

```go
package main

import (
    ...
    "github.com/joho/godotenv"
)

func main() {

    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

}
```

## 2. Changelog

- 25 Jun 2021: added roles
- 24 Jun 2021: initiate project

## 3. Feedback and Contribute

Always welcome. Please [open a new thread](https://github.com/hieunc229/eveauth/issues/new)

## 4. Licenses

- eveauth MIT
- BoltDB MIT
