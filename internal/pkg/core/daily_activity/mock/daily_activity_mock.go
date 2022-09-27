// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/core/daily_activity/daily_activity.go

// Package mock_daily_activity is a generated GoMock package.
package mock_daily_activity

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockInterface) Add(ctx context.Context, name string, act models.DailyActivity, chat models.Chat) (models.DailyActivity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, name, act, chat)
	ret0, _ := ret[0].(models.DailyActivity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockInterfaceMockRecorder) Add(ctx, name, act, chat interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockInterface)(nil).Add), ctx, name, act, chat)
}

// Delete mocks base method.
func (m *MockInterface) Delete(ctx context.Context, name string, chatID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, name, chatID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockInterfaceMockRecorder) Delete(ctx, name, chatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockInterface)(nil).Delete), ctx, name, chatID)
}

// Get mocks base method.
func (m *MockInterface) Get(ctx context.Context, name string, chatID int64) (models.DailyActivity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, name, chatID)
	ret0, _ := ret[0].(models.DailyActivity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockInterfaceMockRecorder) Get(ctx, name, chatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInterface)(nil).Get), ctx, name, chatID)
}

// List mocks base method.
func (m *MockInterface) List(ctx context.Context, chatID int64, pagination models.Pagination) (map[string]models.DailyActivity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, chatID, pagination)
	ret0, _ := ret[0].(map[string]models.DailyActivity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockInterfaceMockRecorder) List(ctx, chatID, pagination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockInterface)(nil).List), ctx, chatID, pagination)
}

// ListStream mocks base method.
func (m *MockInterface) ListStream(ctx context.Context, chatID int64, ch chan models.DailyActivityRec, chErr chan error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListStream", ctx, chatID, ch, chErr)
}

// ListStream indicates an expected call of ListStream.
func (mr *MockInterfaceMockRecorder) ListStream(ctx, chatID, ch, chErr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStream", reflect.TypeOf((*MockInterface)(nil).ListStream), ctx, chatID, ch, chErr)
}

// Today mocks base method.
func (m *MockInterface) Today(ctx context.Context, chatID int64) (map[string]models.DailyActivity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Today", ctx, chatID)
	ret0, _ := ret[0].(map[string]models.DailyActivity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Today indicates an expected call of Today.
func (mr *MockInterfaceMockRecorder) Today(ctx, chatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Today", reflect.TypeOf((*MockInterface)(nil).Today), ctx, chatID)
}

// Update mocks base method.
func (m *MockInterface) Update(ctx context.Context, name string, act models.DailyActivity, chatID int64) (models.DailyActivity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, name, act, chatID)
	ret0, _ := ret[0].(models.DailyActivity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockInterfaceMockRecorder) Update(ctx, name, act, chatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockInterface)(nil).Update), ctx, name, act, chatID)
}
