package db

// TODO remove after this will be full not useful
//import (
//	"errors"
//	"fmt"
//	"sort"
//	"sync"
//	"time"
//
//	"github.com/TheVovchenskiy/sportify-backend/models"
//	"github.com/TheVovchenskiy/sportify-backend/pkg/common"
//
//	"github.com/google/uuid"
//)
//
//type SimpleEventStorage struct {
//	mu *sync.RWMutex
//	m  map[uuid.UUID]models.FullEvent
//}
//
//func NewSimpleEventStorage() (*SimpleEventStorage, error) {
//	s := SimpleEventStorage{
//		mu: &sync.RWMutex{},
//		m:  make(map[uuid.UUID]models.FullEvent),
//	}
//
//	return &s, s.FillSimpleEventStorage()
//}
//
//var events = []models.FullEvent{
//	{
//		ShortEvent: models.ShortEvent{
//			ID:          uuid.New(),
//			SportType:   models.SportTypeFootball,
//			Address:     "г. Москва Госпитальный пер. 4/6 стр. 3",
//			Date:        time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC),
//			StartTime:   time.Date(2024, 10, 14, 20, 0, 0, 0, time.UTC),
//			EndTime:     nil,
//			Price:       common.Ref(0),
//			IsFree:      true,
//			GameLevels:   common.Ref(models.GameLevelMidMinus),
//			Capacity:    nil,
//			Busy:        0,
//			Subscribers: nil,
//			URLPreview:  "http://127.0.0.1:8080/img/default_football.jpeg",
//			URLPhotos:   nil,
//		},
//		Description: common.Ref("Приходите все! Чисто игровая тренировка"),
//		RawMessage:  nil,
//	},
//	{
//		ShortEvent: models.ShortEvent{
//			ID:          uuid.New(),
//			SportType:   models.SportTypeFootball,
//			Address:     "г. Москва Госпитальный пер. 4/6 стр. 3",
//			Date:        time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC),
//			StartTime:   time.Date(2024, 10, 15, 18, 0, 0, 0, time.UTC),
//			EndTime:     common.Ref(time.Date(2024, 10, 15, 21, 0, 0, 0, time.UTC)),
//			Price:       common.Ref(700),
//			IsFree:      false,
//			GameLevels:   common.Ref(models.GameLevelMidPlus),
//			Capacity:    common.Ref(22),
//			Busy:        0,
//			Subscribers: nil,
//			URLPreview:  "http://127.0.0.1:8080/img/default_football.jpeg",
//			URLPhotos:   nil,
//		},
//		Description: common.Ref("Половину тренировки отрабатываем схему 4-4-2, вторая половина игровая"),
//		RawMessage:  nil,
//	},
//	{
//		ShortEvent: models.ShortEvent{
//			ID:          uuid.New(),
//			SportType:   models.SportTypeFootball,
//			Address:     "г. Москва Госпитальный пер. 4/6 стр. 3",
//			Date:        time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC),
//			StartTime:   time.Date(2024, 10, 15, 20, 0, 0, 0, time.UTC),
//			EndTime:     nil,
//			Price:       common.Ref(1000),
//			IsFree:      false,
//			GameLevels:   common.Ref(models.GameLevelMid),
//			Capacity:    nil,
//			Busy:        0,
//			Subscribers: nil,
//			URLPreview:  "http://127.0.0.1:8080/img/default_football.jpeg",
//			URLPhotos:   nil,
//		},
//		Description: common.Ref("Сегодня чисто игровая тренировка. " +
//			"Вход с улицы напротив школы. На проходной скажите, что на игру"),
//		RawMessage: nil,
//	},
//}
//
//func (s *SimpleEventStorage) FillSimpleEventStorage() error {
//	for _, event := range events {
//		err := s.AddEvent(event)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//

//
//func (s *SimpleEventStorage) AddEvent(event models.FullEvent) error {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	if _, ok := s.m[event.ID]; ok {
//		return ErrEventAlreadyExist
//	}
//
//	s.m[event.ID] = event
//
//	return nil
//}
//
//func (s *SimpleEventStorage) GetEvents() ([]models.ShortEvent, error) {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	events := make([]models.ShortEvent, 0)
//	for _, v := range s.m {
//		events = append(events, v.ShortEvent)
//	}
//
//	sort.Slice(events, func(i, j int) bool {
//		return events[j].StartTime.After(events[i].StartTime)
//	})
//
//	return events, nil
//}
//
//
//func (s *SimpleEventStorage) GetEvent(id uuid.UUID) (*models.FullEvent, error) {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	event, ok := s.m[id]
//	if !ok {
//		return nil, ErrNotFoundEvent
//	}
//
//	return &event, nil
//}
//
//func (s *SimpleEventStorage) SubscribeEvent(
//	id uuid.UUID,
//	userID uuid.UUID,
//	subscribe bool,
//) (*models.ResponseSubscribeEvent, error) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	event, ok := s.m[id]
//	if !ok {
//		return nil, ErrNotFoundEvent
//	}
//
//	if subscribe {
//		subscribes, err := event.AddSubscriber(userID)
//		if err != nil {
//			return nil, fmt.Errorf("add subscriber: %w", err)
//		}
//
//		s.m[id] = event
//
//		return &models.ResponseSubscribeEvent{
//			ID: event.ID, Subscribers: subscribes, Capacity: event.Capacity, Busy: event.Busy,
//		}, nil
//	}
//
//	subscribes, err := event.RemoveSubscriber(userID)
//	if err != nil {
//		return nil, fmt.Errorf("remove subscriber: %w", err)
//	}
//
//	s.m[id] = event
//
//	return &models.ResponseSubscribeEvent{
//		ID: event.ID, Subscribers: subscribes, Capacity: event.Capacity, Busy: event.Busy,
//	}, nil
//}
