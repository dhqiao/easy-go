package main

import (
	"time"
	"easy-go/appserver"
	"easy-go/conf"
	"easy-go/routers"
)

// door of this project
func main() {
	mux := appserver.NewAppServeMux()
	// load route config
	routers.Route(mux)
	// starting up the server
	config := conf.Config
	server := &appserver.AppServer{}
	server.Addr = config.Address
	server.Handler = mux
	server.ReadTimeout = time.Duration(config.ReadTimeout * int64(time.Second))
	server.WriteTimeout = time.Duration(config.WriteTimeout * int64(time.Second))
	server.MaxHeaderBytes = 1 << 20
	server.ListenAndServe()
}
