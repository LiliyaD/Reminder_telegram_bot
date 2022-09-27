//go:build integration
// +build integration

package tests

import (
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	postgresPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/database/postgresql"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"github.com/LiliyaD/Reminder_telegram_bot/tests/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ActClient pb.AdminClient
	Storage   postgresPkg.InterfaceStorage
)

var (
	beginDateS         = "10.07.2022"
	endDateS           = "15.09.2023"
	actName            = "push-up"
	act                models.DailyActivity
	chat               models.Chat
	beginDate, endDate time.Time
	activitiesMap      map[string]models.DailyActivity
	pagination         models.Pagination
)

func init() {
	// Logger
	journal.New("test", false)

	// Config
	cfg := config.FromEnv()

	// Test data base
	testDb := config.Connect()
	testDb.Truncate()
	Storage = postgresPkg.NewDataStoragePrep(testDb.Pool)

	// Test grpc server
	conns, err := grpc.Dial(cfg.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		journal.LogFatal(err)
	}
	//defer conns.Close()
	ActClient = pb.NewAdminClient(conns)

	prepareData()
}

func prepareData() {
	var err error
	beginDate, err = time.Parse("02.01.2006", beginDateS)
	if err != nil {
		journal.LogFatal(err)
	}

	endDate, err = time.Parse("02.01.2006", endDateS)
	if err != nil {
		journal.LogFatal(err)
	}

	act = models.DailyActivity{
		BeginDate:       beginDate,
		EndDate:         endDate,
		TimesPerDay:     2,
		QuantityPerTime: 30}

	chat = models.Chat{
		ChatID:   777,
		UserName: "agent007"}

	pagination = models.Pagination{
		Limit:  10,
		Offset: 0,
		Order:  "name",
	}

	activitiesMap = make(map[string]models.DailyActivity)
	activitiesMap[actName] = act
}
