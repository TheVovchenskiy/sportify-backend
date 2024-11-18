package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

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

func (p *PostgresStorage) GetPasswordByUsername(ctx context.Context, username string) (string, error) {
	sqlSelect := `SELECT password FROM "public".user WHERE username = $1;`
	row := p.pool.QueryRow(ctx, sqlSelect, username)

	var password string

	err := row.Scan(&password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%w: %s", ErrUserNotFound, username)
		}

		return "", fmt.Errorf("to scan password: %w", err)
	}

	return password, nil
}

func (p *PostgresStorage) CreateUser(
	ctx context.Context,
	id uuid.UUID,
	username,
	password string,
	tgUserID *int64,
) (models.ResponseSuccessLogin, error) {
	sqlInsert := `INSERT INTO "public".user(id, tg_id, username, password) VALUES($1, $2, $3, $4)`

	_, err := p.pool.Exec(ctx, sqlInsert, id, tgUserID, username, password)
	if err != nil {
		return models.ResponseSuccessLogin{}, fmt.Errorf("to exec: %w", err)
	}

	return models.ResponseSuccessLogin{UserID: id, Username: username}, nil
}
