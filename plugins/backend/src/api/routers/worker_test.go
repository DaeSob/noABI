package routers

import (
	"cia/api/preference"
	"testing"
)

func Test_(t *testing.T) {
	preference.Initialize("local")

	rt := TRouterTable{
		[]TRouterProps{
			TRouterProps{
				Method:   GET,
				Path:     "/ping",
				Function: Ping,
				Auth:     true,
			},
			TRouterProps{
				Method:   GET,
				Path:     "/management/health/ping",
				Function: Ping,
				Auth:     false,
			},
		},
	}

	Run(rt)
}
