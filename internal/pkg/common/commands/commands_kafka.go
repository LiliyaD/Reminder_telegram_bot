package common

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	commonPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/common"
	redis "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/database/cache"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var clientRedis redis.InterfaceRedis

func init() {
	clientRedis = redis.Connect()
}

func createKafka(ctx context.Context, actName string, act models.DailyActivity, chat models.Chat, client pb.AdminClient) (*pb.ActivityCreateResponse, error) {
	actReq := models.DailyActivityCreationReq{
		Name:          actName,
		DailyActivity: act,
		Chat:          chat,
	}

	actReqBytes, err := json.Marshal(actReq)
	if err != nil {
		journal.LogError(errors.Wrap(err, "Marshal error"))
		return nil, err
	}

	uid := (uuid.New()).String()
	err = commonPkg.SendToKafka(uid, "Create", actReqBytes)
	if err != nil {
		return nil, err
	}

	if clientRedis != nil {
		for {
			var answ *models.DailyActivityAnsw
			err := clientRedis.Get(uid, &answ)
			if err != nil {
				journal.LogWarn(errors.Wrap(err, "Redis error, will try again in 10 millisecond"))
				time.Sleep(10 * time.Millisecond)
			} else {
				if len(answ.Error) != 0 {
					journal.LogError("Redis error" + answ.Error)
					return nil, errors.New(answ.Error)
				}

				d := answ.BeginDate
				beginDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
				d = answ.EndDate
				endDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())

				return &pb.ActivityCreateResponse{
					Name:            answ.Name,
					BeginDate:       beginDateR,
					EndDate:         endDateR,
					TimesPerDay:     uint32(answ.TimesPerDay),
					QuantityPerTime: answ.QuantityPerTime,
				}, nil
			}
		}
	}

	return nil, nil
}

func updateKafka(ctx context.Context, chatID int64, actName string, act models.DailyActivity, client pb.AdminClient) (*pb.ActivityUpdateResponse, error) {
	actReq := models.DailyActivityUpdateReq{
		ChatID:        chatID,
		Name:          actName,
		DailyActivity: act,
	}

	actReqBytes, err := json.Marshal(actReq)
	if err != nil {
		journal.LogError("Marshal error")
		return nil, err
	}

	uid := (uuid.New()).String()
	err = commonPkg.SendToKafka(uid, "Update", actReqBytes)
	if err != nil {
		return nil, err
	}

	if clientRedis != nil {
		for {
			var answ *models.DailyActivityAnsw
			err := clientRedis.Get(uid, &answ)
			if err != nil {
				journal.LogWarn(errors.Wrap(err, "Redis error, will try again in 10 millisecond"))
				time.Sleep(10 * time.Millisecond)
			} else {
				if len(answ.Error) != 0 {
					journal.LogError("Redis error" + answ.Error)
					return nil, errors.New(answ.Error)
				}

				d := answ.BeginDate
				beginDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
				d = answ.EndDate
				endDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())

				return &pb.ActivityUpdateResponse{
					Name:            answ.Name,
					BeginDate:       beginDateR,
					EndDate:         endDateR,
					TimesPerDay:     uint32(answ.TimesPerDay),
					QuantityPerTime: answ.QuantityPerTime,
				}, nil
			}
		}
	}

	return nil, nil
}

func deleteKafka(ctx context.Context, chatID int64, actName string, client pb.AdminClient) error {
	actReq := models.DailyActivityDeletionReq{
		ChatID: chatID,
		Name:   actName,
	}

	actReqBytes, errM := json.Marshal(actReq)
	if errM != nil {
		journal.LogError("Marshal error")
		return errM
	}

	uid := (uuid.New()).String()

	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		var answ *models.DailyActivityDelAnsw
		clientRedis.Subscribe(uid, &answ)
		if len(answ.Error) != 0 {
			journal.LogError("Redis error" + answ.Error)
			err = errors.New(answ.Error)
		}
		wg.Done()
	}()

	err = commonPkg.SendToKafka(uid, "Delete", actReqBytes)
	if err != nil {
		return err
	}

	wg.Wait()

	return err
}
