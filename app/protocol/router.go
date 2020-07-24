package protocol

import (
	"fmt"

	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/netcloth/netcloth-chain/utils"
)

type Router struct {
	routes map[string]sdk.Handler
}

var _ sdk.Router = NewRouter()

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]sdk.Handler),
	}
}

func (rtr *Router) AddRoute(path string, h sdk.Handler) sdk.Router {
	if !utils.IsAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if rtr.routes[path] != nil {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	rtr.routes[path] = h
	return rtr
}

func (rtr *Router) Route(_ sdk.Context, path string) sdk.Handler {
	return rtr.routes[path]
}
