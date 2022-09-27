package config

import (
	"os"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"

	"github.com/pkg/errors"
)

const ServerAddress = ":8081"
const ClientAddress = ":8082"

func GetApiKey() string {
	apiKey, exists := os.LookupEnv("TELEGRAM_BOT_API_KEY")
	if !exists {
		journal.LogFatal(errors.New("Environment variable TELEGRAM_BOT_API_KEY doesn't exist"))
	}

	return apiKey
}
