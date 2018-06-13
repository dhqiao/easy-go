package routers

import (
	"easy-go/appserver"
	"easy-go/controllers"
	"easy-go/library"
)

var indexController = &controllers.IndexController{}
var memberController = &controllers.MemberController{}

var logger = library.Logger

// 路由器
func Route(mux *appserver.AppServeMux) {
	mux.Middleware("LoginCheck").AppHandleFunc("/test", indexController.Test)
	mux.AppHandleFunc("/validator", indexController.ValidatorTest)
	mux.Middleware("Pass").AppHandleFunc("/", memberController.GetList)
	mux.AppHandleFunc("/redis", indexController.Redis)

}
