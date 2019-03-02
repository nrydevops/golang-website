package main

var user_cache map[string]*User


func init() {

	user_cache = make(map[string]*User)
}

func SignIn(username, password string) (string, error) {

	user := GetUser(username, password)
	if user == nil {

		return "", authorizationError{"User not recognized"}
	}

	session_token, err := GenerateSessionToken(user)
	if err != nil {

		return "", authorizationError{"Can't create session token for user"}
	}

	user_cache[session_token] = user

	return session_token, nil;
}

func LogOut(session_token string) {

	delete(user_cache, session_token);
}
