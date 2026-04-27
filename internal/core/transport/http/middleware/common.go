package core_http_middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_logger "github.com/povarejkee/backend-tech-practice/internal/core/logger"
	core_http_response "github.com/povarejkee/backend-tech-practice/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			rw.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(rw, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", l)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicReponse(p, "during handle http req got unexpected panic")
				}
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
