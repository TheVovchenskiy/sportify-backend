package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"

	"github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

var ErrEventAlreadyExist = errors.New("Событие уже существует")

func (p *PostgresStorage) CreateEvent(ctx context.Context, event *models.FullEvent) error {
	sqlInsertEvent := `
	INSERT INTO "public".event (
    id, creator_id, subscriber_ids, sport_type, address, date_start, start_time, end_time,
    price, game_level, description, raw_message, capacity, busy, creation_type,
    url_message, url_author, url_preview, url_photos, tg_chat_id, tg_message_id
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, 
          $9, $10, $11, $12, $13, $14, $15,
          $16, $17, $18, $19, $20, $21);`

	preparedGameLevel := pq.Array(event.GameLevels)

	_, err := p.pool.Exec(ctx, sqlInsertEvent,
		event.ID, event.CreatorID, event.Subscribers, event.SportType, event.Address,
		event.DateAndTime.Date, event.DateAndTime.StartTime, event.DateAndTime.EndTime, event.Price, preparedGameLevel,
		event.Description, event.RawMessage, event.Capacity, event.Busy, event.CreationType,
		event.URLMessage, event.URLAuthor, event.URLPreview, event.URLPhotos, event.TgChatID, event.TgMessageID)
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
		url_author = $13, url_preview = $14, url_photos = $15,
		coordinates = ST_Point($16, $17, 4326)::geography
		WHERE id = $18 AND deleted_at IS NULL;`

	preparedGameLevels := pq.Array(event.GameLevels)

	_, err := p.pool.Exec(ctx, sqlUpdateEvent,
		event.CreatorID, event.SportType, event.Address,
		event.DateAndTime.Date, event.DateAndTime.StartTime, event.DateAndTime.EndTime, event.Price, preparedGameLevels,
		event.Description, event.Capacity, event.CreationType, event.URLMessage,
		event.URLAuthor, event.URLPreview, event.URLPhotos, event.Latitude, event.Longitude, event.ID)
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

var ErrNotFoundEvent = errors.New("Не найдено событие")

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

func (p *PostgresStorage) GetEventByTgChatAndMessageIDs(ctx context.Context, tgChatID, tgMessageID int64) (*models.FullEvent, error) {
	sqlSelectEvent := `
	SELECT id, creator_id, subscriber_ids, sport_type, address, date_start, start_time, end_time,
       price, game_level, description, raw_message, capacity, busy, creation_type,
       url_author, url_message, 
       url_preview, url_photos,
       ST_X(coordinates::geometry) as latitude, ST_Y(coordinates::geometry) as longitude,
	   tg_chat_id, tg_message_id
	FROM "public".event WHERE tg_chat_id = $1 AND $2 = tg_message_id AND deleted_at IS NULL;`

	rawRow := p.pool.QueryRow(ctx, sqlSelectEvent, tgChatID, tgMessageID)

	var (
		event            models.FullEvent
		rawSubscriberIDs pgtype.Array[uuid.UUID]
		rawURLPhotos     pgtype.Array[string]
		rawGameLevels    pgtype.Array[*string]
	)

	err := rawRow.Scan(&event.ID, &event.CreatorID, &rawSubscriberIDs, &event.SportType, &event.Address,
		&event.DateAndTime.Date, &event.DateAndTime.StartTime, &event.DateAndTime.EndTime, &event.Price, &rawGameLevels,
		&event.Description, &event.RawMessage, &event.Capacity, &event.Busy, &event.CreationType,
		&event.URLAuthor, &event.URLMessage, &event.URLPreview, &rawURLPhotos, &event.Latitude, &event.Longitude,
		&event.TgChatID, &event.TgMessageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFoundEvent
		}

		return nil, fmt.Errorf("to scan event: %w", err)
	}

	event.Subscribers = rawSubscriberIDs.Elements
	event.URLPhotos = rawURLPhotos.Elements
	event.IsFree = *event.Price == 0
	event.GameLevels = models.GameLevelFromRawNullable(rawGameLevels.Elements)

	return &event, nil
}

func (p *PostgresStorage) GetEvent(ctx context.Context, eventID uuid.UUID) (*models.FullEvent, error) {
	sqlSelectEvent := `
	SELECT creator_id, subscriber_ids, sport_type, address, date_start, start_time, end_time,
       price, game_level, description, raw_message, capacity, busy, creation_type,
       url_author, url_message, 
       url_preview, url_photos,
       ST_X(coordinates::geometry) as latitude, ST_Y(coordinates::geometry) as longitude,
	   tg_chat_id, tg_message_id
	FROM "public".event WHERE id = $1 AND deleted_at IS NULL;`

	rawRow := p.pool.QueryRow(ctx, sqlSelectEvent, eventID)

	var (
		event            models.FullEvent
		rawSubscriberIDs pgtype.Array[uuid.UUID]
		rawURLPhotos     pgtype.Array[string]
		rawGameLevels    pgtype.Array[*string]
	)

	err := rawRow.Scan(&event.CreatorID, &rawSubscriberIDs, &event.SportType, &event.Address,
		&event.DateAndTime.Date, &event.DateAndTime.StartTime, &event.DateAndTime.EndTime, &event.Price, &rawGameLevels,
		&event.Description, &event.RawMessage, &event.Capacity, &event.Busy, &event.CreationType,
		&event.URLAuthor, &event.URLMessage, &event.URLPreview, &rawURLPhotos, &event.Latitude, &event.Longitude,
		&event.TgChatID, &event.TgMessageID)
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

func NewPostgresStorage(ctx context.Context, urlDataBase string) (*PostgresStorage, *pgxpool.Pool, error) {
	pool, err := repository.NewPgxPool(ctx, urlDataBase)
	if err != nil {
		return nil, nil, err
	}

	return &PostgresStorage{pool: pool}, pool, nil
}

func getSQLEvents(rawRows pgx.Rows) ([]models.ShortEvent, error) {
	var (
		curEvent      models.ShortEvent
		result        []models.ShortEvent
		photoURLs     pgtype.Array[string]
		rawGameLevels pgtype.Array[*string]
	)

	_, err := pgx.ForEachRow(
		rawRows,
		[]any{
			&curEvent.ID, &curEvent.CreatorID, &curEvent.SportType, &curEvent.Address, &curEvent.DateAndTime.Date,
			&curEvent.DateAndTime.StartTime, &curEvent.DateAndTime.EndTime, &curEvent.Price, &rawGameLevels,
			&curEvent.Capacity, &curEvent.Busy, &curEvent.Subscribers,
			&curEvent.URLPreview, &photoURLs, &curEvent.Latitude, &curEvent.Longitude,
		},
		func() error {
			result = append(
				result, models.ShortEvent{
					ID:        curEvent.ID,
					CreatorID: curEvent.CreatorID,
					SportType: curEvent.SportType,
					Address:   curEvent.Address,
					DateAndTime: models.DateAndTime{
						Date:      curEvent.DateAndTime.Date,
						StartTime: curEvent.DateAndTime.StartTime,
						EndTime:   curEvent.DateAndTime.EndTime,
					},
					Price:       curEvent.Price,
					IsFree:      *curEvent.Price == 0,
					GameLevels:  models.GameLevelFromRawNullable(rawGameLevels.Elements),
					Capacity:    curEvent.Capacity,
					Busy:        curEvent.Busy,
					Subscribers: curEvent.Subscribers,
					URLPreview:  curEvent.URLPreview,
					URLPhotos:   photoURLs.Elements,
					Latitude:    curEvent.Latitude,
					Longitude:   curEvent.Longitude,
				})

			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}

	if result == nil {
		return []models.ShortEvent{}, nil
	}

	return result, nil
}

func (p *PostgresStorage) GetEvents(ctx context.Context) ([]models.ShortEvent, error) {
	sqlSelect := `
	SELECT id, creator_id, sport_type, address, date_start, start_time,
       end_time, price, game_level, capacity, busy,
       subscriber_ids, url_preview, url_photos,
       ST_X(coordinates::geometry) as latitude, ST_Y(coordinates::geometry) as longitude
	FROM "public".event WHERE deleted_at IS NULL AND start_time > NOW() - INTERVAL '1 day'
	ORDER BY start_time;`

	rawRows, err := p.pool.Query(ctx, sqlSelect)
	if err != nil {
		return nil, fmt.Errorf("select events: %w", err)
	}

	return getSQLEvents(rawRows)
}

//nolint:lll
func (p *PostgresStorage) FindEvents(ctx context.Context, filterParams *models.FilterParams) ([]models.ShortEvent, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, fmt.Errorf("get logger: %w", err)
	}

	logger.WithCtx(ctx).Infow("Getting events", zap.Any("filter_params", filterParams))

	query := squirrel.Select(`id, creator_id, sport_type, address, date_start, start_time,
		end_time, price, game_level, capacity, busy,
		subscriber_ids, url_preview, url_photos,
		ST_X(coordinates::geometry) as latitude, ST_Y(coordinates::geometry) as longitude`).
		From(`"public".event`).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"deleted_at": nil})
	// Where(squirrel.Gt{"start_time": time.Now().Add(-24 * time.Hour)}) // TODO: add later

	if filterParams.CreatorID != nil {
		query = query.Where(squirrel.Eq{"creator_id": filterParams.CreatorID})
	}

	if len(filterParams.SportTypes) > 0 {
		query = query.Where(squirrel.Eq{"sport_type": filterParams.SportTypes})
	}

	if len(filterParams.GameLevels) > 0 {
		query = query.Where("game_level && ?", pq.Array(filterParams.GameLevels))
	}

	if len(filterParams.DateStarts) > 0 {
		query = query.Where(squirrel.Eq{"date_start": filterParams.DateStarts})
	}

	if filterParams.DateExpression != nil {
		query = query.Where(filterParams.DateExpression)
	}

	if len(filterParams.SubscriberIDs) > 0 {
		query = query.Where("subscriber_ids && ?", pq.Array(filterParams.SubscriberIDs))
	}

	if filterParams.PriceMin != nil {
		query = query.Where(squirrel.GtOrEq{"price": *filterParams.PriceMin})
	}

	if filterParams.PriceMax != nil {
		query = query.Where(squirrel.LtOrEq{"price": *filterParams.PriceMax})
	}

	if filterParams.FreePlaces != nil {
		query = query.Where(squirrel.Expr("capacity - busy >= ?", *filterParams.FreePlaces))
	}

	if filterParams.Address != "" {
		if filterParams.AddressLatitude != nil && filterParams.AddressLongitude != nil {
			query = query.Where(fmt.Sprintf("ST_DWithin(ST_POINT(%s,%s, 4326)::geography, coordinates, 5000.0)",
				*filterParams.AddressLatitude, *filterParams.AddressLongitude))
		}
		// TODO add find address or by reqexp or another text search
	}

	query = query.OrderBy(filterParams.OrderBy + " " + filterParams.SortOrder)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query to sql: %w", err)
	}

	logger.WithCtx(ctx).Infow("SQL query", zap.String("query", sql), zap.Any("args", args))

	rawRows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("select events: %w", err)
	}

	return getSQLEvents(rawRows)
}

func (p *PostgresStorage) AddUserPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	sqlUpdate := `
	UPDATE public.event SET user_paid_ids = ARRAY_APPEND(user_paid_ids, $1) WHERE id = $2;`

	_, err := p.pool.Exec(ctx, sqlUpdate, userID, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) SetCoordinates(ctx context.Context, latitude, longitude string, id uuid.UUID) error {
	sqlUpdate := `UPDATE public.event SET coordinates = ST_Point( $1, $2, 4326)::geography WHERE id = $3`

	_, err := p.pool.Exec(ctx, sqlUpdate, latitude, longitude, id)
	if err != nil {
		return err
	}

	return nil
}
