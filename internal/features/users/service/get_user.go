package users_service

import (
	"context"
	"fmt"

	"github.com/povarejkee/backend-tech-practice/internal/core/domain"
	core_errors "github.com/povarejkee/backend-tech-practice/internal/core/errors"
)

func (s *UsersService) GetUser(ctx context.Context, id int) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repo: %w", core_errors.ErrInvalidArgument)
	}

	return user, nil
}
