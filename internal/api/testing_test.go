package api

import (
	"testing"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	mock_activity "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/mock"
	"github.com/golang/mock/gomock"
)

// for api.go
type serverFixture struct {
	activity *mock_activity.MockInterface
	service  *implementation
}

func connectMock(t *testing.T) serverFixture {
	journal.New("test", false)
	t.Parallel()
	f := serverFixture{}
	f.activity = mock_activity.NewMockInterface(gomock.NewController(t))
	f.service = new(f.activity)
	return f
}

func new(activity *mock_activity.MockInterface) *implementation {
	return &implementation{
		activity: activity,
	}
}
