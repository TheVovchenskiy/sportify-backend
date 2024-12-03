package app

import (
	"context"

	"github.com/TheVovchenskiy/sportify-backend/models"

	"github.com/google/uuid"
)

func (a *App) GetUserFullByUserID(ctx context.Context, userID uuid.UUID) (*models.UserFull, error) {
	userFull, err := a.authStorage.GetUserFullByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return userFull, nil
}
