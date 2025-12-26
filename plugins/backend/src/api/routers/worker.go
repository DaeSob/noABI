package routers

import (
	"cia/api/middleware"
	"cia/api/preference"
	"cia/common/errors"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(_router TRouterTable) error {
	return worker(_router)
}

func worker(_router TRouterTable) error {
	// gin mode
	var r *gin.Engine
	if preference.IsDebugMode() {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	// set middleware
	onSetMiddleware(r)

	// set router
	onSetRouters(r, _router)

	// start gin
	r.Run(preference.GetHost())

	return nil
}

func onSetMiddleware(_r *gin.Engine) {
	if !preference.IsDebugMode() {
		_r.Use(gin.Recovery())
	}

	_r.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"POST", "GET"},
				AllowHeaders:     []string{"*"},
				AllowCredentials: true,
			},
		),
	)
}

func onSetRouters(_r *gin.Engine, _router TRouterTable) {
	// auth
	// path prefix
	authPath := "/" + preference.PathPrefix()
	authRequired := _r.Group(authPath)
	if preference.IsEnableAuth() {
		authRequired.Use(middleware.AuthCheckDate)
		authRequired.Use(middleware.AuthCheckSignature)
	} else {
		authRequired.Use(middleware.NoAuth)
	}

	// no auth
	noAuth := _r.Group("/")

	for _, props := range _router.Table {
		if props.Auth {
			_setRouter(authRequired, props)
		} else {
			_setRouter(noAuth, props)
		}
	}
}

func _setRouter(_rg *gin.RouterGroup, _props TRouterProps) {
	switch _props.Method {
	case GET:
		_rg.GET(_props.Path, _props.Function)
	case POST:
		_rg.POST(_props.Path, _props.Function)
	default:
		panic(errors.TError{"", "unsupported router method"})
	}
}
