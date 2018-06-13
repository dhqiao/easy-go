package middleware

import (
	"errors"
	"net/http"
)

func (middleware *Middleware) LoginCheck(writer http.ResponseWriter, request *http.Request) error {

	//params, _ := url.ParseQuery(request.URL.RawQuery)

	return errors.New("new meddleware logincheck error")
}
