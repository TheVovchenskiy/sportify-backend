package app

import (
	"context"
	"errors"
	"fmt"

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

var ErrValidationRequestUpdateProfile = errors.New("не правильные параметры")

func (a *App) UpdateProfile(ctx context.Context, userID uuid.UUID, reqUpdate models.RequestUpdateProfile) error {
	err := reqUpdate.Valid()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationRequestUpdateProfile, err)
	}

	err = a.authStorage.UpdateProfile(ctx, userID, reqUpdate)
	if err != nil {
		return fmt.Errorf("to update profile: %w", err)
	}

	return nil
}
