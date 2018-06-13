package middleware

import (
	"net/http"
)

func (middleware *Middleware) Pass(writer http.ResponseWriter, request *http.Request) error {

	//params, _ := url.ParseQuery(request.URL.RawQuery)

	return nil
}
