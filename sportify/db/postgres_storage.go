package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/TheVovchenskiy/sportify-backend/models"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

var ErrEventAlreadyExist = errors.New("событие уже существует")

func (p *PostgresStorage) CreateEvent(ctx context.Context, event *models.FullEvent) error {
	sqlInsertEvent := `
	INSERT INTO "public".event (
    id, creator_id, subscriber_ids, sport_type, address, date_start, start_time, end_time,
    price, game_level, description, raw_message, capacity, busy, creation_type,
    url_message, url_author, url_preview, url_photos
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, 
          $9, $10, $11, $12, $13, $14, $15,
          $16, $17, $18, $19);`

	preparedGameLevel := pgtype.Array[string]{Elements: models.RawGameLevel(event.GameLevels)}

	_, err := p.pool.Exec(ctx, sqlInsertEvent,
		event.ID, event.CreatorID, event.Subscribers, event.SportType, event.Address,
		event.Date, event.StartTime, event.EndTime, event.Price, preparedGameLevel,
		event.Description, event.RawMessage, event.Capacity, event.Busy, event.CreationType,
		event.URLMessage, event.URLAuthor, event.URLPreview, event.URLPhotos)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) EditEvent(ctx context.Context, event *models.FullEvent) error {
	sqlUpdateEvent := `
	UPDATE "public".event SET creator_id = $1, sport_type = $2, address = $3, 
		date_start = $4, start_time = $5, end_time = $6, price = $7, game_level = $8,
		description = $9, capacity = $10, creation_type = $11, url_message = $12, 
		url_author = $13, url_preview = $14, url_photos = $15
		WHERE id = $16 AND deleted_at IS NULL;`

	preparedGameLevels := pgtype.Array[string]{Elements: models.RawGameLevel(event.GameLevels)}

	_, err := p.pool.Exec(ctx, sqlUpdateEvent,
		event.CreatorID, event.SportType, event.Address,
		event.Date, event.StartTime, event.EndTime, event.Price, preparedGameLevels,
		event.Description, event.Capacity, event.CreationType, event.URLMessage,
		event.URLAuthor, event.URLPreview, event.URLPhotos, event.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	sqlDelete := `UPDATE "public".event SET deleted_at = NOW() WHERE id = $1 AND creator_id = $2`

	_, err := p.pool.Exec(ctx, sqlDelete, eventID, userID)
	if err != nil {
		return err
	}

	return nil
}

var ErrNotFoundEvent = errors.New("не найдено событие")

func (p *PostgresStorage) GetCreatorID(ctx context.Context, eventID uuid.UUID) (uuid.UUID, error) {
	sqlSelectEvent := `
		SELECT creator_id FROM "public".event WHERE id = $1 AND deleted_at IS NULL;`

	rawRow := p.pool.QueryRow(ctx, sqlSelectEvent, eventID)

	var creatorID uuid.UUID

	err := rawRow.Scan(&creatorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, ErrNotFoundEvent
		}

		return uuid.UUID{}, fmt.Errorf("to scan event: %w", err)
	}

	return creatorID, nil
}

func (p *PostgresStorage) GetEvent(ctx context.Context, eventID uuid.UUID) (*models.FullEvent, error) {
	sqlSelectEvent := `
	SELECT creator_id, subscriber_ids, sport_type, address, date_start, start_time, end_time,
       price, game_level, description, raw_message, capacity, busy, creation_type,
       url_author, url_message, 
       url_preview, url_photos FROM "public".event WHERE id = $1 AND deleted_at IS NULL;`

	rawRow := p.pool.QueryRow(ctx, sqlSelectEvent, eventID)

	var (
		event            models.FullEvent
		rawSubscriberIDs pgtype.Array[uuid.UUID]
		rawURLPhotos     pgtype.Array[string]
		rawGameLevels    pgtype.Array[*string]
	)

	err := rawRow.Scan(&event.CreatorID, &rawSubscriberIDs, &event.SportType, &event.Address,
		&event.Date, &event.StartTime, &event.EndTime, &event.Price, &rawGameLevels,
		&event.Description, &event.RawMessage, &event.Capacity, &event.Busy, &event.CreationType,
		&event.URLAuthor, &event.URLMessage, &event.URLPreview, &rawURLPhotos)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFoundEvent
		}

		return nil, fmt.Errorf("to scan event: %w", err)
	}

	event.ID = eventID
	event.Subscribers = rawSubscriberIDs.Elements
	event.URLPhotos = rawURLPhotos.Elements
	event.IsFree = *event.Price == 0
	event.GameLevels = models.GameLevelFromRawNullable(rawGameLevels.Elements)

	return &event, nil
}

func (p *PostgresStorage) updateEventSubscribe(
	ctx context.Context,
	subEvent *models.ResponseSubscribeEvent,
) error {
	sqlUpdateSub := `
	UPDATE event SET subscriber_ids = $1, capacity = $2, busy = $3
		WHERE id = $4 AND deleted_at IS NULL`

	_, err := p.pool.Exec(ctx, sqlUpdateSub, subEvent.Subscribers, subEvent.Capacity, subEvent.Busy, subEvent.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) SubscribeEvent(
	ctx context.Context,
	eventID uuid.UUID,
	userID uuid.UUID,
	subscribe bool,
) (*models.ResponseSubscribeEvent, error) {
	// TODO add support of creator_id event notify
	sqlSelectEvent := `
	SELECT subscriber_ids, busy, capacity FROM "public".event 
	                                      WHERE id = $1 AND deleted_at IS NULL;`

	rawRow := p.pool.QueryRow(ctx, sqlSelectEvent, eventID)

	var (
		rawSubscriberIDs       pgtype.Array[uuid.UUID]
		responseSubscribeEvent models.ResponseSubscribeEvent
	)

	err := rawRow.Scan(&rawSubscriberIDs, &responseSubscribeEvent.Busy, &responseSubscribeEvent.Capacity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFoundEvent
		}

		return nil, fmt.Errorf("to scan event for subcribe: %w", err)
	}

	responseSubscribeEvent.Subscribers = rawSubscriberIDs.Elements

	if subscribe {
		err = responseSubscribeEvent.AddSubscriber(userID)
		if err != nil {
			return nil, fmt.Errorf("add subscriber: %w", err)
		}
	} else {
		err = responseSubscribeEvent.RemoveSubscriber(userID)
		if err != nil {
			return nil, fmt.Errorf("remove subscriber: %w", err)
		}
	}

	responseSubscribeEvent.ID = eventID

	err = p.updateEventSubscribe(ctx, &responseSubscribeEvent)
	if err != nil {
		return nil, fmt.Errorf("to update event subscribe: %w", err)
	}

	return &responseSubscribeEvent, nil
}

func NewPostgresStorage(ctx context.Context, urlDataBase string) (*PostgresStorage, error) {
	pool, err := repository.NewPgxPool(ctx, urlDataBase)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{pool: pool}, nil
}

func (p *PostgresStorage) GetEvents(ctx context.Context) ([]models.ShortEvent, error) {
	sqlSelect := `SELECT id, creator_id, sport_type, address, date_start, start_time,
		end_time, price, game_level, capacity, busy,
		subscriber_ids, url_preview, url_photos
		FROM "public".event WHERE deleted_at IS NULL;`

	rawRows, err := p.pool.Query(ctx, sqlSelect)
	if err != nil {
		return nil, fmt.Errorf("select events: %w", err)
	}

	var (
		curEvent      models.ShortEvent
		result        []models.ShortEvent
		photoURLs     pgtype.Array[string]
		rawGameLevels pgtype.Array[*string]
	)

	_, err = pgx.ForEachRow(
		rawRows,
		[]any{
			&curEvent.ID, &curEvent.CreatorID, &curEvent.SportType, &curEvent.Address, &curEvent.Date,
			&curEvent.StartTime, &curEvent.EndTime, &curEvent.Price, &rawGameLevels,
			&curEvent.Capacity, &curEvent.Busy, &curEvent.Subscribers,
			&curEvent.URLPreview, &photoURLs,
		},
		func() error {
			result = append(
				result, models.ShortEvent{
					ID:          curEvent.ID,
					CreatorID:   curEvent.CreatorID,
					SportType:   curEvent.SportType,
					Address:     curEvent.Address,
					Date:        curEvent.Date,
					StartTime:   curEvent.StartTime,
					EndTime:     curEvent.EndTime,
					Price:       curEvent.Price,
					IsFree:      *curEvent.Price == 0,
					GameLevels:  models.GameLevelFromRawNullable(rawGameLevels.Elements),
					Capacity:    curEvent.Capacity,
					Busy:        curEvent.Busy,
					Subscribers: curEvent.Subscribers,
					URLPreview:  curEvent.URLPreview,
					URLPhotos:   photoURLs.Elements,
				})

			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("to get events: %w", err)
	}

	return result, nil
}
