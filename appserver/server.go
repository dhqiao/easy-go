/*
simple middleware, not standard style
*/
package appserver

import (
	"easy-go/appserver/middleware"
	"easy-go/controllers"
	"easy-go/library"
	"fmt"
	. "net/http"
	"reflect"
)

// 日志
var logger = library.Logger

// middleware实体
var middlewareObj = &middleware.Middleware{}

// 返回数据
var responseObj = &controllers.BaseController{}

type GenericMiddleware func(ResponseWriter, *Request, NextMiddlewareFunc)

type NextMiddlewareFunc func(ResponseWriter, *Request)

type AppServer struct {
	Server
}

var middlewares []string
var nullMiddleware []string

type AppServeMux struct {
	ServeMux
	appMiddleware map[string][]string
	before        bool
	path          string
}

func NewAppServeMux() *AppServeMux {
	return new(AppServeMux)
}

// Middleware adds the specified middleware tot he router and returns the router.
func (mux *AppServeMux) Middleware(filterNames ...interface{}) *AppServeMux {
	for _, filterName := range filterNames {
		middlewares = append(middlewares, reflect.ValueOf(filterName).String())
	}
	return mux
}

// extend base function HandleFunc
func (mux *AppServeMux) AppHandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	mux.before = true
	mux.path = pattern
	if len(mux.appMiddleware) == 0 {
		mux.appMiddleware = make(map[string][]string)
	}
	mux.appMiddleware[pattern] = middlewares
	middlewares = nullMiddleware
	mux.HandleFunc(pattern, handler)
}

// override base class's function
func (mux *AppServeMux) ServeHTTP(w ResponseWriter, r *Request) {

	if mux.before {
		err := mux.executeMiddleware(r.URL.Path, w, r)
		if err != nil {
			fmt.Fprintf(w, responseObj.ErrorData(err.Error()))
			return
		}
	}
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
	mux.before = false
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

// execute middleware
func (mux *AppServeMux) executeMiddleware(path string, w ResponseWriter, r *Request) error {
	middlewareLists := mux.findMiddleware(path)
	middlewareLen := len(middlewareLists)

	if middlewareLen == 0 {
		return nil
	}
	for i := 0; i < middlewareLen; i++ {
		outValue := library.InvokeObjectMethod(middlewareObj, middlewareLists[i], w, r)
		outValLen := len(outValue)
		if outValLen >= 1 && !outValue[outValLen-1].IsNil() {
			err := outValue[outValLen-1].Elem().Interface().(error)
			return err
		}
	}
	return nil
}

// 从middleware容器中匹配当前路径配置的middleware
// 可能会有bug，此处是匹配当前的path，如果匹配不到就只想根路径，和实际上go的匹配可能会有差异
func (mux *AppServeMux) findMiddleware(path string) []string {
	var routeMiddleware []string

	if len(mux.appMiddleware) == 0 {
		return routeMiddleware
	}
	for k, v := range mux.appMiddleware {
		if k == "/" {
			routeMiddleware = v
		}
		if k == path {
			return v
		}
	}
	return routeMiddleware
}
