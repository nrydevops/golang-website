package main

import (
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

func GetUserFromRequest(req *http.Request) *User {

	token, err := req.Cookie("Token")

	if err != nil {

		return nil
	}

	return GetUserByToken(token.Value)
}

func GetUserByToken(session_token string) *User {

	user, ok := user_cache[session_token]

	if ok {

		return user

	} else {

		return nil
	}
}

func RefreshSession(res *http.ResponseWriter, req *http.Request) error {

	cookie, err := req.Cookie("Token")

	if err != nil {

		return err;
	}

	user := GetUserByToken(cookie.Value)
	if user == nil {

		return authorizationError{"User not found"}
	}

	if user.JustSigned {

		GetUserByToken(cookie.Value).JustSigned = false

	} else {

		new_session_token, err := RefreshSessionToken(cookie.Value)
		if err != nil {

			return err;
		}

		http.SetCookie(*res, &http.Cookie{
			Name:    "Token",
			Value:   new_session_token,
			Path:    "/",
			Expires: time.Now().Add(300 * time.Second),
		})
	}

	return nil
}

func RefreshSessionToken(previous_token string) (string, error) {

	user, ok := user_cache[previous_token]
	if !ok {

		return "", authorizationError{"User not recognized"}
	}
	delete(user_cache, previous_token)

	newToken, err := GenerateSessionToken(user)
	if err != nil {

		return "", err
	}

	user_cache[newToken] = user

	return newToken, nil
}

func IsFakeToken(session_token string) bool {

	_, ok := user_cache[session_token];

	return !ok
}

func GenerateSessionToken(user *User) (string, error) {

	var id uuid.UUID
	var err error

	for {

		// Create a new random session token
		id, err = uuid.NewV4()
		if err != nil {

			return "", err
		}

		if _,contains := user_cache[id.String()]; !contains {

			break
		}
	}


	return id.String(), nil
}
