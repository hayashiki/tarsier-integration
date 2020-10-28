package handler

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

type Wrap func(w http.ResponseWriter, r *http.Request) error

func NewHTTPError(code int, msg interface{}) error {
	return HTTPError{Code: code, Message: msg}
}

type data map[string]interface{}

// Implementation of the error interface
func (e HTTPError) Error() string {
	return fmt.Sprintf("code: %d; message: %v", e.Code, e.Message)
}

func (h Wrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		switch err.(type) {
		case HTTPError:
			er := err.(HTTPError)
			jsonResponse(w, er.Code, data{"error": er.Message})
			return
		}
		jsonResponse(w, http.StatusInternalServerError, data{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
}
