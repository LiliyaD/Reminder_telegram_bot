package models

import "time"

type DailyActivity struct {
	BeginDate       time.Time `db:"begin_date"`
	EndDate         time.Time `db:"end_date"`
	TimesPerDay     uint8     `db:"times_per_day"`
	QuantityPerTime float32   `db:"quantity_per_time"`
}

type DailyActivityRec struct {
	Name string `db:"act_name"`
	DailyActivity
}

type Chat struct {
	ChatID   int64  `db:"chatID"`
	UserName string `db:"user_name"`
}

type Pagination struct {
	Limit  uint64
	Offset uint64
	Order  string
}

type DailyActivityCreationReq struct {
	Name string `db:"act_name"`
	DailyActivity
	Chat
}

type DailyActivityAnsw struct {
	Error  string
	ChatID int64  `db:"chatID"`
	Name   string `db:"act_name"`
	DailyActivity
}

type DailyActivityDelAnsw struct {
	Error string
}

type DailyActivityUpdateReq struct {
	ChatID int64  `db:"chatID"`
	Name   string `db:"act_name"`
	DailyActivity
}

type DailyActivityDeletionReq struct {
	ChatID int64  `db:"chatID"`
	Name   string `db:"act_name"`
}

type CacheAnswerValue struct {
	Error error
	Name  string `db:"act_name"`
	DailyActivity
}
