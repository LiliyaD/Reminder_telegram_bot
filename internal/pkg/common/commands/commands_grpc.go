package common

import (
	"context"
	"fmt"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
)

func createGrpc(ctx context.Context, actName string, act models.DailyActivity, chat models.Chat, client pb.AdminClient) (*pb.ActivityCreateResponse, error) {
	beginDate := fmt.Sprintf("%02d.%02d.%04d", act.BeginDate.Day(), act.BeginDate.Month(), act.BeginDate.Year())
	endDate := fmt.Sprintf("%02d.%02d.%04d", act.EndDate.Day(), act.EndDate.Month(), act.EndDate.Year())

	response, err := client.ActivityCreate(ctx, &pb.ActivityCreateRequest{
		Name:            actName,
		BeginDate:       beginDate,
		EndDate:         endDate,
		TimesPerDay:     uint32(act.TimesPerDay),
		QuantityPerTime: act.QuantityPerTime,
		ChatID:          chat.ChatID,
		UserName:        chat.UserName})
	if err != nil {
		journal.LogError(err)
		return nil, err
	}

	return response, nil
}

func updateGrpc(ctx context.Context, chatID int64, actName string, act models.DailyActivity, client pb.AdminClient) (*pb.ActivityUpdateResponse, error) {
	beginDate := fmt.Sprintf("%02d.%02d.%04d", act.BeginDate.Day(), act.BeginDate.Month(), act.BeginDate.Year())
	endDate := fmt.Sprintf("%02d.%02d.%04d", act.EndDate.Day(), act.EndDate.Month(), act.EndDate.Year())

	response, err := client.ActivityUpdate(ctx, &pb.ActivityUpdateRequest{
		Name:            actName,
		BeginDate:       beginDate,
		EndDate:         endDate,
		TimesPerDay:     uint32(act.TimesPerDay),
		QuantityPerTime: act.QuantityPerTime,
		ChatID:          chatID})
	if err != nil {
		journal.LogError(err)
		return nil, err
	}

	return response, nil
}

func deleteGrpc(ctx context.Context, chatID int64, actName string, client pb.AdminClient) error {
	_, err := client.ActivityDelete(ctx, &pb.ActivityDeleteRequest{ChatID: chatID, Name: actName})
	if err != nil {
		journal.LogError(err)
		return err
	}

	return nil
}
