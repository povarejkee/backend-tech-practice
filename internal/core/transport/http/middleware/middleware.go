package core_http_middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func ChainMw(h http.Handler, mw ...Middleware) http.Handler {
	if len(mw) == 0 {
		return h
	}

	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}

	return h
}
