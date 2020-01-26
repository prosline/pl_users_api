package errors

import "net/http"

type RestErr struct {
	Message string `json:"message"`
	Code int `json:"code"`
	Error string `json:"error"`
}

func BadRequestError(message string) *RestErr{
	return &RestErr{
		Message: message,
		Code:    http.StatusBadRequest,
		Error:   "Bad Request",
	}
}
func CreateUserError(message string) *RestErr{
	return &RestErr{
		Message: message,
		Code:    http.StatusInternalServerError,
		Error:   "Can't create user",
	}
}
