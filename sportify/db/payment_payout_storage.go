package db

import (
	"context"
	"errors"

	"github.com/TheVovchenskiy/sportify-backend/models"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPaymentPayoutStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresPaymentPayoutStorage(pool *pgxpool.Pool) *PostgresPaymentPayoutStorage {
	return &PostgresPaymentPayoutStorage{pool: pool}
}

func (p *PostgresPaymentPayoutStorage) CreatePayment(ctx context.Context, payment *models.Payment) error {
	sqlInsert := `
	INSERT INTO public.payment (id, user_id, event_id, confirmation_url, status, amount)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := p.pool.Exec(ctx,
		sqlInsert, payment.ID, payment.UserID, payment.EventID, payment.ConfirmationURL, payment.Status, payment.Amount,
	)
	if err != nil {
		// TODO check uniq violation constraint or do it before
		return err
	}

	return nil
}

var ErrNotFoundPayment = errors.New("платеж не найден")

func (p *PostgresPaymentPayoutStorage) GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	var result models.Payment

	sqlSelect := `SELECT id, user_id, event_id, confirmation_url, status, amount FROM public.payment WHERE id = $1`

	err := p.pool.QueryRow(ctx, sqlSelect, &id).Scan(
		&result.ID,
		&result.UserID,
		&result.EventID,
		&result.ConfirmationURL,
		&result.Status,
		&result.Amount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFoundPayment
		}

		return nil, err
	}

	return &result, nil
}

func (p *PostgresPaymentPayoutStorage) UpdateStatusPayment(
	ctx context.Context,
	id uuid.UUID,
	status models.PaymentStatus,
) error {
	sqlUpdate := `UPDATE public.payment SET status = $1 WHERE id = $2`

	_, err := p.pool.Exec(ctx, sqlUpdate, status, id)
	if err != nil {
		return err
	}

	return nil
}
