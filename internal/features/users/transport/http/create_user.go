package users_transport_http

import (
	"encoding/json"
	"net/http"

	core_logger "github.com/povarejkee/backend-tech-practice/internal/core/logger"
	core_http_request "github.com/povarejkee/backend-tech-practice/internal/core/transport/http/request"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	log.Debug("invoce CreateUser handler")

	var req CreateUserRequest
	if err := core_http_request.DecodeAndValidate(r, &req); err != nil {
		// todo
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// todo: handle error
	}

	rw.WriteHeader(http.StatusOK)
}
