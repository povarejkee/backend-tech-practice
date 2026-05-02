package core_http_server

import (
	"fmt"
	"net/http"
)

type APIVersion string

var (
	ApiVersion1 = APIVersion("v1")
	ApiVersion2 = APIVersion("v2")
	ApiVersion3 = APIVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
}

func NewAPIVersionRouter(apiVer APIVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVer,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		// GET /tasks
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.Handler)
	}
}
