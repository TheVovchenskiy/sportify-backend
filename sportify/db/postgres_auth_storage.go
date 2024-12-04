package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/common"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *PostgresStorage) GetUserFullByID(ctx context.Context, id uuid.UUID) (*models.UserFull, error) {
	sqlSelect := `
		SELECT id, tg_id, username, password, created_at, updated_at, 
       		first_name, second_name, sport_types, photo_url, description
		FROM "public".user WHERE id = $1;`

	row := p.pool.QueryRow(ctx, sqlSelect, id)

	var (
		user          models.UserFull
		rawSportTypes pgtype.Array[string]
	)

	err := row.Scan(
		&user.ID, &user.TgID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		&user.FirstName, &user.SecondName, &rawSportTypes, &user.PhotoURL, &user.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %d", ErrUserNotFound, id)
		}

		return nil, fmt.Errorf("to scan user: %w", err)
	}

	user.SportTypes = common.Map(func(item string) models.SportType {
		return models.SportType(item)
	}, rawSportTypes.Elements)

	return &user, nil
}

func (p *PostgresStorage) GetUserFullByTgID(ctx context.Context, tgID int64) (*models.UserFull, error) {
	sqlSelect := `
	SELECT id, tg_id, username, password, created_at, updated_at,
		first_name, second_name, sport_types, photo_url, description
	FROM "public".user WHERE tg_id = $1;`

	row := p.pool.QueryRow(ctx, sqlSelect, tgID)

	var (
		user          models.UserFull
		rawSportTypes pgtype.Array[string]
	)

	err := row.Scan(
		&user.ID, &user.TgID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		&user.FirstName, &user.SecondName, &rawSportTypes, &user.PhotoURL, &user.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %d", ErrUserNotFound, tgID)
		}

		return nil, fmt.Errorf("to scan user: %w", err)
	}

	user.SportTypes = common.Map(func(item string) models.SportType {
		return models.SportType(item)
	}, rawSportTypes.Elements)

	return &user, nil
}

func (p *PostgresStorage) GetUserFullByUsername(ctx context.Context, username string) (*models.UserFull, error) {
	sqlSelect := `
	SELECT id, tg_id, username, password, created_at, updated_at,
		first_name, second_name, sport_types, photo_url, description 
	FROM "public".user WHERE username = $1;`

	row := p.pool.QueryRow(ctx, sqlSelect, username)

	var (
		user          models.UserFull
		rawSportTypes pgtype.Array[string]
	)

	err := row.Scan(
		&user.ID, &user.TgID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		&user.FirstName, &user.SecondName, &rawSportTypes, &user.PhotoURL, &user.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", ErrUserNotFound, username)
		}

		return nil, fmt.Errorf("to scan user: %w", err)
	}

	user.SportTypes = common.Map(func(item string) models.SportType {
		return models.SportType(item)
	}, rawSportTypes.Elements)

	return &user, nil
}

func (p *PostgresStorage) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	sqlSelect := `SELECT username FROM "public".user WHERE username = $1;`

	row := p.pool.QueryRow(ctx, sqlSelect, username)

	var usernameFromDB string

	err := row.Scan(&usernameFromDB)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

var ErrUserNotFound = errors.New("не найден пользователь")

func (p *PostgresStorage) GetPasswordByUsername(ctx context.Context, username string) (*string, error) {
	sqlSelect := `SELECT password FROM "public".user WHERE username = $1;`
	row := p.pool.QueryRow(ctx, sqlSelect, username)

	var password *string

	err := row.Scan(&password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", ErrUserNotFound, username)
		}

		return nil, fmt.Errorf("to scan password: %w", err)
	}

	return password, nil
}

func (p *PostgresStorage) CreateUser(
	ctx context.Context,
	id uuid.UUID,
	username string,
	password *string,
	tgUserID *int64,
) (models.ResponseSuccessLogin, error) {
	sqlInsert := `INSERT INTO "public".user(id, tg_id, username, password) VALUES($1, $2, $3, $4)`

	_, err := p.pool.Exec(ctx, sqlInsert, id, tgUserID, username, password)
	if err != nil {
		return models.ResponseSuccessLogin{}, fmt.Errorf("to exec: %w", err)
	}

	return models.ResponseSuccessLogin{UserID: id, Username: username, TgUserID: tgUserID}, nil
}

func (p *PostgresStorage) UpdateProfile(ctx context.Context, userID uuid.UUID, reqUpdate models.RequestUpdateProfile) error {
	sqlUpdate := `UPDATE "public".user SET 
	first_name = $1, second_name = $2, photo_url = $3, description = $4, sport_types = $5
	WHERE id = $6;`

	rawSportTypes := common.Map(func(item models.SportType) string {
		return string(item)
	}, reqUpdate.SportTypes)

	_, err := p.pool.Exec(ctx, sqlUpdate,
		reqUpdate.FirstName, reqUpdate.SecondName, reqUpdate.PhotoURL, reqUpdate.Description, rawSportTypes,
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}
