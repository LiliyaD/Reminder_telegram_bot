package multithread

import (
	cachePkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache"
	"github.com/pkg/errors"
)

var (
	ErrNotExist             = errors.New("This activity doesn't exist")
	ErrEmptyStorage         = errors.New("Storage is empty")
	ErrExist                = errors.New("This activity has already existed")
	ErrInvalidParameter     = errors.New("Invalid parameter")
	ErrNoActivitiesForToday = errors.New("No activities for today")
	ErrOk                   = errors.New("Ok")
)

type InterfaceStorage cachePkg.Interface
