package users_service

import (
	"context"

	"github.com/povarejkee/backend-tech-practice/internal/core/domain"
)

func (s UsersService) GetUsers(ctx context.Context) ([]domain.User, error) {
	// todo s.usersRepository.GetUsers(ctx)

	return nil, nil
}
