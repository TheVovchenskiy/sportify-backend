package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sort"
	"sync"
	"time"
)

//
//две ручки /events - список, /event/id - подробные
//ручка PUT запись на событие /event/sub/id {“sub” true или false, “user_id”:”uuid”}

type EventStorage interface {
	AddEvent(event FullEvent) error
	GetEvents() ([]ShortEvent, error)
	GetEvent(id uuid.UUID) (*FullEvent, error)
	SubscribeEvent(id uuid.UUID, userID uuid.UUID, subscribe bool) (*ResponseSubscribeEvent, error)
}

type SimpleEventStorage struct {
	mu *sync.RWMutex
	m  map[uuid.UUID]FullEvent
}

func NewSimpleEventStorage() (*SimpleEventStorage, error) {
	s := SimpleEventStorage{
		mu: &sync.RWMutex{},
		m:  make(map[uuid.UUID]FullEvent),
	}

	return &s, s.FillSimpleEventStorage()
}

var events = []FullEvent{
	{
		ShortEvent: ShortEvent{
			ID:          uuid.New(),
			SportType:   TypeFootball,
			Address:     "г. Москва Госпитальный пер. 4/6 стр. 3",
			Date:        time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC),
			StartTime:   time.Date(2024, 10, 14, 20, 0, 0, 0, time.UTC),
			EndTime:     nil,
			Price:       Ref(0),
			IsFree:      true,
			GameLevel:   Ref(GameLevelMidMinus),
			Capacity:    nil,
			Busy:        0,
			Subscribers: nil,
			PreviewURL:  "http://127.0.0.1:8080/img/default_football.jpeg",
			PhotoURLs:   nil,
		},
		Description: Ref("Приходите все! Чисто игровая тренировка"),
		RawMessage:  nil,
	},
	{
		ShortEvent: ShortEvent{
			ID:          uuid.New(),
			SportType:   TypeFootball,
			Address:     "г. Москва Госпитальный пер. 4/6 стр. 3",
			Date:        time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC),
			StartTime:   time.Date(2024, 10, 15, 18, 0, 0, 0, time.UTC),
			EndTime:     Ref(time.Date(2024, 10, 15, 21, 0, 0, 0, time.UTC)),
			Price:       Ref(700),
			IsFree:      false,
			GameLevel:   Ref(GameLevelMidPlus),
			Capacity:    Ref(22),
			Busy:        0,
			Subscribers: nil,
			PreviewURL:  "http://127.0.0.1:8080/img/default_football.jpeg",
			PhotoURLs:   nil,
		},
		Description: Ref("Половину тренировки отрабатываем схему 4-4-2, вторая половина игровая"),
		RawMessage:  nil,
	},
	{
		ShortEvent: ShortEvent{
			ID:          uuid.New(),
			SportType:   TypeFootball,
			Address:     "г. Москва Госпитальный пер. 4/6 стр. 3",
			Date:        time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC),
			StartTime:   time.Date(2024, 10, 15, 20, 0, 0, 0, time.UTC),
			EndTime:     nil,
			Price:       Ref(1000),
			IsFree:      false,
			GameLevel:   Ref(GameLevelMid),
			Capacity:    nil,
			Busy:        0,
			Subscribers: nil,
			PreviewURL:  "http://127.0.0.1:8080/img/default_football.jpeg",
			PhotoURLs:   nil,
		},
		Description: Ref("Сегодня чисто игровая тренировка. Вход с улицы напротив школы. На проходной скажите, что на игру"),
		RawMessage:  nil,
	},
}

func (s *SimpleEventStorage) FillSimpleEventStorage() error {
	for _, event := range events {
		err := s.AddEvent(event)
		if err != nil {
			return err
		}
	}

	return nil
}

var ErrEventAlreadyExist = errors.New("event already exists")

func (s *SimpleEventStorage) AddEvent(event FullEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.m[event.ID]; ok {
		return ErrEventAlreadyExist
	}

	s.m[event.ID] = event

	return nil
}

func (s *SimpleEventStorage) GetEvents() ([]ShortEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]ShortEvent, 0)
	for _, v := range s.m {
		events = append(events, v.ShortEvent)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[j].StartTime.After(events[i].StartTime)
	})

	return events, nil
}

var ErrNotFoundEvent = errors.New("not found event")

func (s *SimpleEventStorage) GetEvent(id uuid.UUID) (*FullEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.m[id]
	if !ok {
		return nil, ErrNotFoundEvent
	}

	return &event, nil
}

func (s *SimpleEventStorage) SubscribeEvent(id uuid.UUID, userID uuid.UUID, subscribe bool) (*ResponseSubscribeEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.m[id]
	if !ok {
		return nil, ErrNotFoundEvent
	}

	if subscribe {
		subscribes, err := event.AddSubscriber(userID)
		if err != nil {
			return nil, fmt.Errorf("add subscriber: %w", err)
		}

		s.m[id] = event

		return &ResponseSubscribeEvent{ID: event.ID, Subscribers: subscribes, Capacity: event.Capacity, Busy: event.Busy}, nil
	}

	subscribes, err := event.RemoveSubscriber(userID)
	if err != nil {
		return nil, fmt.Errorf("remove subscriber: %w", err)
	}

	s.m[id] = event

	return &ResponseSubscribeEvent{ID: event.ID, Subscribers: subscribes, Capacity: event.Capacity, Busy: event.Busy}, nil
}
