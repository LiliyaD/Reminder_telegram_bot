package api

import (
	"context"
	"testing"
	"time"

	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivityCreate(t *testing.T) {
	beginDateS := "10.07.2022"
	beginDate, err := time.Parse("02.01.2006", beginDateS)
	if err != nil {
		t.Fatal(err)
	}
	endDateS := "15.09.2023"
	endDate, err := time.Parse("02.01.2006", endDateS)
	if err != nil {
		t.Fatal(err)
	}

	act := models.DailyActivity{
		BeginDate:       beginDate,
		EndDate:         endDate,
		TimesPerDay:     2,
		QuantityPerTime: 30}

	chat := models.Chat{
		ChatID:   777,
		UserName: "agent007"}

	actName := "push-up"

	t.Run("Success service creation new activity", func(t *testing.T) {
		// arrange
		f := connectMock(t)

		f.activity.EXPECT().Add(gomock.Any(), actName, act, chat).Return(act, nil)

		// act
		resp, err := f.service.ActivityCreate(context.Background(),
			&pb.ActivityCreateRequest{
				Name:            actName,
				BeginDate:       beginDateS,
				EndDate:         endDateS,
				TimesPerDay:     uint32(act.TimesPerDay),
				QuantityPerTime: act.QuantityPerTime,
				ChatID:          chat.ChatID,
				UserName:        chat.UserName})

		// assert
		require.NoError(t, err)
		assert.Equal(t, resp,
			&pb.ActivityCreateResponse{
				Name:            actName,
				BeginDate:       beginDateS,
				EndDate:         endDateS,
				TimesPerDay:     uint32(act.TimesPerDay),
				QuantityPerTime: act.QuantityPerTime})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("AlreadyExists: activity exists", func(t *testing.T) {
			// arrange
			f := connectMock(t)

			f.activity.EXPECT().Add(gomock.Any(), actName, act, chat).Return(models.DailyActivity{}, multiPkg.ErrExist)

			// act
			_, err := f.service.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			// assert
			msg := "rpc error: code = AlreadyExists desc = " + multiPkg.ErrExist.Error()
			assert.EqualError(t, err, msg)
		})

		t.Run("Internal: timeout", func(t *testing.T) {
			// arrange
			f := connectMock(t)

			f.activity.EXPECT().Add(gomock.Any(), actName, act, chat).Return(models.DailyActivity{}, multiPkg.ErrTimeOut)

			// act
			_, err := f.service.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			// assert
			msg := "rpc error: code = Internal desc = " + multiPkg.ErrTimeOut.Error()
			assert.EqualError(t, err, msg)
		})

		t.Run("InvalidArgument: chatID should be filled", func(t *testing.T) {
			// arrange
			f := connectMock(t)

			// act
			_, err := f.service.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					// ChatID:          chat.ChatID,
					UserName: chat.UserName})

			// assert
			assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = chatID should be filled")
		})

		t.Run("InvalidArgument: name should be filled", func(t *testing.T) {
			// arrange
			f := connectMock(t)

			// act
			_, err := f.service.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					// Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			// assert
			assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = name should be filled")
		})

		t.Run("InvalidArgument: BeginDate (without separator)", func(t *testing.T) {
			// arrange
			f := connectMock(t)

			// act
			_, err := f.service.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       "10122022",
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			// assert
			assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Date can't be parsed")
		})

		t.Run("InvalidArgument: EndDate (without separator)", func(t *testing.T) {
			// arrange
			f := connectMock(t)

			// act
			_, err := f.service.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         "10122022",
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			// assert
			assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Date can't be parsed")
		})

		t.Run("InvalidArgument: EndDate (empty)", func(t *testing.T) {
			// arrange
			f := connectMock(t)

			// act
			_, err := f.service.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         "",
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			// assert
			assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Date can't be parsed")
		})
	})
}
