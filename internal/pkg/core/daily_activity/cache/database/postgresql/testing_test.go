package postgres

import (
	"testing"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/edbmanniwood/pgxpoolmock"
	"github.com/golang/mock/gomock"
)

type postgresFixture struct {
	storage  InterfaceStorage
	mockPool *pgxpoolmock.MockPgxPool
	ctrl     *gomock.Controller
}

func connectMock(t *testing.T) postgresFixture {
	// Logger
	journal.New("test", false)

	// Fixture for postgresql
	var fixture postgresFixture
	fixture.ctrl = gomock.NewController(t)
	fixture.mockPool = pgxpoolmock.NewMockPgxPool(fixture.ctrl)
	fixture.storage = &cache{pool: fixture.mockPool}

	return fixture
}

func (f *postgresFixture) closeMockConnection() {
	f.ctrl.Finish()
}
