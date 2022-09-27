package main

import (
	"os"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	dailyActPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity"
	cachePkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache"
	redis "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/database/cache"
	postgresPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/database/postgresql"
	localPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/local"
)

var (
	activity  dailyActPkg.Interface
	cacheAnsw redis.InterfaceRedis
)

func main() {
	journal.New("server", false)
	journal.LogInfo("RUN Server")

	var dataStorage cachePkg.Interface
	if len(os.Args) > 1 && os.Args[1] == "--local" {
		dataStorage = localPkg.NewDataStorage()
	} else {
		dataStorage = postgresPkg.NewDataStorage()
		cacheAnsw = postgresPkg.GetCache()
		defer dataStorage.(postgresPkg.InterfaceStorage).CloseDataStorage()
	}

	activity = dailyActPkg.New(dataStorage)

	go runCount()

	go consume()

	go runREST()

	runGRPCServer(activity)
}
