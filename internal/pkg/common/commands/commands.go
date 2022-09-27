package common

import (
	"context"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"github.com/pkg/errors"
)

var (
	errLessDate        = errors.New("End date is less than begin date")
	errEmptyName       = errors.New("Name can't be empty")
	errEmptyDate       = errors.New("Date can't be empty")
	errZeroTimesPerDay = errors.New("Times per day can't be 0")
	errQuantityPerTime = errors.New("Quantity per time can't be 0")
)

func Create(ctx context.Context, actName string, act models.DailyActivity, chat models.Chat, client pb.AdminClient) (*pb.ActivityCreateResponse, error) {
	// checks
	if actName == "" {
		journal.LogError(errEmptyName)
		return nil, errEmptyName
	}

	if act.TimesPerDay == 0 {
		journal.LogError(errZeroTimesPerDay)
		return nil, errZeroTimesPerDay
	}

	if act.QuantityPerTime == 0 {
		journal.LogError(errQuantityPerTime)
		return nil, errQuantityPerTime
	}

	timeEmpt := time.Time{}
	if act.BeginDate == timeEmpt || act.EndDate == timeEmpt {
		journal.LogError(errEmptyDate)
		return nil, errEmptyDate
	}

	if act.EndDate.Before(act.BeginDate) {
		journal.LogError(errLessDate)
		return nil, errLessDate
	}

	return createKafka(ctx, actName, act, chat, client)
}

func List(ctx context.Context, chatID int64, pagination models.Pagination, client pb.AdminClient) (*pb.ActivityListResponse, error) {
	response, err := client.ActivityList(ctx, &pb.ActivityListRequest{
		ChatID: chatID,
		Limit:  &(pagination.Limit),
		Offset: &(pagination.Offset),
		Order:  &(pagination.Order)})
	if err != nil {
		journal.LogError(err)
		return nil, err
	}

	return response, nil
}

func ListStream(ctx context.Context, chatID int64, client pb.AdminClient) (pb.Admin_ActivityListStreamClient, error) {
	stream, err := client.ActivityListStream(ctx, &pb.ActivityListStreamRequest{ChatID: chatID})
	if err != nil {
		journal.LogError(err)
		return nil, err
	}

	return stream, nil
}

func Today(ctx context.Context, chatID int64, client pb.AdminClient) (*pb.ActivityTodayResponse, error) {
	stream, err := client.ActivityToday(ctx, &pb.ActivityTodayRequest{ChatID: chatID})
	if err != nil {
		journal.LogError(err)
		return nil, err
	}

	return stream, nil
}

func Get(ctx context.Context, chatID int64, actName string, client pb.AdminClient) (*pb.ActivityGetResponse, error) {
	response, err := client.ActivityGet(ctx, &pb.ActivityGetRequest{ChatID: chatID, Name: actName})
	if err != nil {
		journal.LogError(err)
		return nil, err
	}

	return response, nil
}

func Delete(ctx context.Context, chatID int64, actName string, client pb.AdminClient) error {
	return deleteKafka(ctx, chatID, actName, client)
}

func Update(ctx context.Context, chatID int64, actName string, act models.DailyActivity, client pb.AdminClient) (*pb.ActivityUpdateResponse, error) {
	return updateKafka(ctx, chatID, actName, act, client)
}
