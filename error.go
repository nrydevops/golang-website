package main

type authorizationError struct {

	Message string
}

type securityError struct {

	Message string
}

func (err authorizationError) Error() string {

	return err.Message
}

func (err securityError) Error() string {

	return err.Message
}