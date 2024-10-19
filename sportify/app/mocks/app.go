// Code generated by MockGen. DO NOT EDIT.
// Source: app.go
//
// Generated by this command:
//
//	mockgen -source=app.go -destination=mocks/app.go -package=mocks EventStorage
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/TheVovchenskiy/sportify-backend/models"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockEventStorage is a mock of EventStorage interface.
type MockEventStorage struct {
	ctrl     *gomock.Controller
	recorder *MockEventStorageMockRecorder
	isgomock struct{}
}

// MockEventStorageMockRecorder is the mock recorder for MockEventStorage.
type MockEventStorageMockRecorder struct {
	mock *MockEventStorage
}

// NewMockEventStorage creates a new mock instance.
func NewMockEventStorage(ctrl *gomock.Controller) *MockEventStorage {
	mock := &MockEventStorage{ctrl: ctrl}
	mock.recorder = &MockEventStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventStorage) EXPECT() *MockEventStorageMockRecorder {
	return m.recorder
}

// AddEvent mocks base method.
func (m *MockEventStorage) AddEvent(event models.FullEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEvent", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEvent indicates an expected call of AddEvent.
func (mr *MockEventStorageMockRecorder) AddEvent(event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockEventStorage)(nil).AddEvent), event)
}

// GetEvent mocks base method.
func (m *MockEventStorage) GetEvent(id uuid.UUID) (*models.FullEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", id)
	ret0, _ := ret[0].(*models.FullEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvent indicates an expected call of GetEvent.
func (mr *MockEventStorageMockRecorder) GetEvent(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockEventStorage)(nil).GetEvent), id)
}

// GetEvents mocks base method.
func (m *MockEventStorage) GetEvents() ([]models.ShortEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents")
	ret0, _ := ret[0].([]models.ShortEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents.
func (mr *MockEventStorageMockRecorder) GetEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockEventStorage)(nil).GetEvents))
}

// SubscribeEvent mocks base method.
func (m *MockEventStorage) SubscribeEvent(id, userID uuid.UUID, subscribe bool) (*models.ResponseSubscribeEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeEvent", id, userID, subscribe)
	ret0, _ := ret[0].(*models.ResponseSubscribeEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeEvent indicates an expected call of SubscribeEvent.
func (mr *MockEventStorageMockRecorder) SubscribeEvent(id, userID, subscribe any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeEvent", reflect.TypeOf((*MockEventStorage)(nil).SubscribeEvent), id, userID, subscribe)
}
