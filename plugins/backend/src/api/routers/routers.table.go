package routers

import (
	"github.com/gin-gonic/gin"
)

type TRouterMethod string

const (
	GET  TRouterMethod = "GET"
	POST               = "POST"
)

type TRouterProps struct {
	Method   TRouterMethod
	Path     string
	Function gin.HandlerFunc
	Auth     bool
}

func GenGetRouter(
	_path string,
	_func gin.HandlerFunc,
	_auth bool,
) TRouterProps {
	return genRouterProps(GET, _path, _func, _auth)
}

func GenPostRouter(
	_path string,
	_func gin.HandlerFunc,
	_auth bool,
) TRouterProps {
	return genRouterProps(POST, _path, _func, _auth)
}

func genRouterProps(
	_method TRouterMethod,
	_path string,
	_func gin.HandlerFunc,
	_auth bool,
) TRouterProps {
	return TRouterProps{
		Method:   _method,
		Path:     _path,
		Function: _func,
		Auth:     _auth,
	}
}

type TRouterTable struct {
	Table []TRouterProps
}

func (rt *TRouterTable) AddRouter(
	_routerProps TRouterProps,
) {
	rt.Table = append(rt.Table, _routerProps)
}

func (rt *TRouterTable) AddGetRouter(
	_path string,
	_func gin.HandlerFunc,
	_auth bool,
) {
	rt.AddRouter(GenGetRouter(_path, _func, _auth))
}

func (rt *TRouterTable) AddPostRouter(
	_path string,
	_func gin.HandlerFunc,
	_auth bool,
) {
	rt.AddRouter(GenPostRouter(_path, _func, _auth))
}

var (
	// router
	ROUTER_PING = GenGetRouter("/management/health/ping", Ping, false)
)
