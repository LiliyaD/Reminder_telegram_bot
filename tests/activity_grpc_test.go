//go:build integration
// +build integration

package tests

import (
	"context"
	"fmt"
	"io"
	"testing"

	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was created", func(t *testing.T) {
			//arrange in init_test.go

			//act
			resp, err := ActClient.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			//assert
			require.NoError(t, err)
			assert.Equal(t, resp.Name, actName)
			assert.Equal(t, resp.BeginDate, beginDateS)
			assert.Equal(t, resp.EndDate, endDateS)
			assert.Equal(t, resp.TimesPerDay, uint32(act.TimesPerDay))
			assert.Equal(t, resp.QuantityPerTime, act.QuantityPerTime)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("InvalidArgument: chatID should be filled", func(t *testing.T) {
			//arrange in init_test.go

			//act
			resp, err := ActClient.ActivityGet(context.Background(), &pb.ActivityGetRequest{Name: actName})

			//assert
			require.EqualError(t, err, "rpc error: code = InvalidArgument desc = chatID should be filled")
			assert.Nil(t, resp)
		})

		t.Run("InvalidArgument: name should be filled", func(t *testing.T) {
			// arrange in init_test.go

			// act
			resp, err := ActClient.ActivityGet(context.Background(), &pb.ActivityGetRequest{ChatID: chat.ChatID})

			//assert
			require.EqualError(t, err, "rpc error: code = InvalidArgument desc = name should be filled")
			assert.Nil(t, resp)
		})

		t.Run("AlreadyExists: activity exists", func(t *testing.T) {
			// arrange in init_test.go

			// act
			resp, err := ActClient.ActivityCreate(context.Background(),
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
			assert.Nil(t, resp)
		})
	})
}

func TestReadActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was found", func(t *testing.T) {
			//arrange in init_test.go
			ActClient.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			//act
			resp, err := ActClient.ActivityGet(context.Background(), &pb.ActivityGetRequest{ChatID: chat.ChatID, Name: actName})
			//assert
			require.NoError(t, err)
			assert.Equal(t, resp.Name, actName)
			assert.Equal(t, resp.BeginDate, beginDateS)
			assert.Equal(t, resp.EndDate, endDateS)
			assert.Equal(t, resp.TimesPerDay, uint32(act.TimesPerDay))
			assert.Equal(t, resp.QuantityPerTime, act.QuantityPerTime)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("InvalidArgument: chatID should be filled", func(t *testing.T) {
			//arrange in init_test.go
			//act
			resp, err := ActClient.ActivityGet(context.Background(), &pb.ActivityGetRequest{Name: actName})
			//assert
			require.EqualError(t, err, "rpc error: code = InvalidArgument desc = chatID should be filled")
			assert.Nil(t, resp)
		})

		t.Run("InvalidArgument: name should be filled", func(t *testing.T) {
			//arrange in init_test.go
			//act
			resp, err := ActClient.ActivityGet(context.Background(), &pb.ActivityGetRequest{ChatID: chat.ChatID})
			//assert
			require.EqualError(t, err, "rpc error: code = InvalidArgument desc = name should be filled")
			assert.Nil(t, resp)
		})

		t.Run("ErrNotExist: activity doesn't exist", func(t *testing.T) {
			// arrange in init_test.go

			// act
			resp, err := ActClient.ActivityGet(context.Background(), &pb.ActivityGetRequest{ChatID: 5, Name: "pushup"})

			// assert
			msg := "rpc error: code = NotFound desc = " + multiPkg.ErrNotExist.Error()
			require.EqualError(t, err, msg)
			assert.Nil(t, resp)
		})
	})
}

func TestListActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was found", func(t *testing.T) {
			//arrange in init_test.go
			ActClient.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			//act
			resp, err := ActClient.ActivityList(context.Background(), &pb.ActivityListRequest{ChatID: chat.ChatID})

			//assert
			require.NoError(t, err)
			require.Equal(t, len(resp.Activities), len(activitiesMap))
			for _, v := range resp.Activities {
				m, ok := activitiesMap[v.Name]
				if ok {
					d := m.BeginDate
					beginDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
					d = m.EndDate
					endDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
					assert.Equal(t, v.BeginDate, beginDateR)
					assert.Equal(t, v.EndDate, endDateR)
					assert.Equal(t, v.TimesPerDay, uint32(m.TimesPerDay))
					assert.Equal(t, v.QuantityPerTime, m.QuantityPerTime)
				} else {
					t.Fail()
				}
			}

		})
	})
}

func TestListStreamActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was found", func(t *testing.T) {
			//arrange in init_test.go
			ActClient.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			//act
			act, err := ActClient.ActivityListStream(context.Background(), &pb.ActivityListStreamRequest{ChatID: chat.ChatID})

			//assert
			require.NoError(t, err)

			for {
				v, err := act.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					require.NoError(t, err)
				}

				m, ok := activitiesMap[v.GetName()]
				if ok {
					d := m.BeginDate
					beginDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
					d = m.EndDate
					endDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
					assert.Equal(t, v.GetBeginDate(), beginDateR)
					assert.Equal(t, v.GetEndDate(), endDateR)
					assert.Equal(t, v.GetTimesPerDay(), uint32(m.TimesPerDay))
					assert.Equal(t, v.GetQuantityPerTime(), m.QuantityPerTime)
				} else {
					t.Fail()
				}
			}
		})
	})
}

func TestUpdateActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was updated", func(t *testing.T) {
			//arrange in init_test.go
			ActClient.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			//act
			resp, err := ActClient.ActivityUpdate(context.Background(),
				&pb.ActivityUpdateRequest{
					ChatID:          chat.ChatID,
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(3),
					QuantityPerTime: act.QuantityPerTime})

			//assert
			require.NoError(t, err)
			assert.Equal(t, resp.Name, actName)
			assert.Equal(t, resp.BeginDate, beginDateS)
			assert.Equal(t, resp.EndDate, endDateS)
			assert.Equal(t, resp.TimesPerDay, uint32(3))
			assert.Equal(t, resp.QuantityPerTime, act.QuantityPerTime)
		})
	})
}

func TestDeleteActivity(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Activity was deleted", func(t *testing.T) {
			//arrange in init_test.go
			ActClient.ActivityCreate(context.Background(),
				&pb.ActivityCreateRequest{
					Name:            actName,
					BeginDate:       beginDateS,
					EndDate:         endDateS,
					TimesPerDay:     uint32(act.TimesPerDay),
					QuantityPerTime: act.QuantityPerTime,
					ChatID:          chat.ChatID,
					UserName:        chat.UserName})

			//act
			_, err := ActClient.ActivityDelete(context.Background(), &pb.ActivityDeleteRequest{ChatID: chat.ChatID, Name: actName})

			//assert
			require.NoError(t, err)
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("ErrNotExist: activity doesn't exist", func(t *testing.T) {
			//arrange in init_test.go

			//act
			resp, err := ActClient.ActivityDelete(context.Background(), &pb.ActivityDeleteRequest{ChatID: chat.ChatID, Name: actName})

			//assert
			msg := "rpc error: code = NotFound desc = " + multiPkg.ErrNotExist.Error()
			require.EqualError(t, err, msg)
			assert.Nil(t, resp)
		})
	})
}
