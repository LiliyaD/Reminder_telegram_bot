package api

import (
	"context"
	"fmt"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal/counter"
	commandPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/bot/command"
	dailyActPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity"
	multiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/cache/multithread"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(activity dailyActPkg.Interface) *implementation {
	return &implementation{
		activity: activity,
	}
}

type implementation struct {
	pb.UnimplementedAdminServer
	activity dailyActPkg.Interface
}

func (imp *implementation) ActivityCreate(ctx context.Context, in *pb.ActivityCreateRequest) (*pb.ActivityCreateResponse, error) {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return nil, err
	}
	chat := models.Chat{ChatID: chatID, UserName: in.GetUserName()}

	actName := in.GetName()
	if len(actName) == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "name should be filled")
		journal.LogError(err)
		return nil, err
	}

	beginDate, err := commandPkg.ParseDate(in.GetBeginDate())
	if err != nil {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, err.Error())
		journal.LogError(err)
		return nil, err
	}

	endDate, err := commandPkg.ParseDate(in.GetEndDate())
	if err != nil {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, err.Error())
		journal.LogError(err)
		return nil, err
	}

	activity, err := imp.activity.Add(ctx, actName, models.DailyActivity{
		BeginDate:       beginDate,
		EndDate:         endDate,
		TimesPerDay:     uint8(in.GetTimesPerDay()),
		QuantityPerTime: in.GetQuantityPerTime(),
	}, chat)
	if err != nil {
		counter.ErrorRequests.Increase()
		code := codes.InvalidArgument
		if err == multiPkg.ErrExist {
			code = codes.AlreadyExists
		} else if err == multiPkg.ErrTimeOut {
			code = codes.Internal
		}
		err := status.Error(code, err.Error())
		journal.LogError(err)
		return nil, err
	}

	d := activity.BeginDate
	beginDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
	d = activity.EndDate
	endDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())

	counter.SuccessRequests.Increase()
	return &pb.ActivityCreateResponse{
		Name:            actName,
		BeginDate:       beginDateR,
		EndDate:         endDateR,
		TimesPerDay:     uint32(activity.TimesPerDay),
		QuantityPerTime: activity.QuantityPerTime,
	}, nil
}

func (imp *implementation) ActivityList(ctx context.Context, in *pb.ActivityListRequest) (*pb.ActivityListResponse, error) {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return nil, err
	}

	pagination := models.Pagination{Limit: in.GetLimit(), Offset: in.GetOffset(), Order: in.GetOrder()}

	activities, err := imp.activity.List(ctx, chatID, pagination)
	if err != nil {
		counter.ErrorRequests.Increase()
		code := codes.InvalidArgument
		if err == multiPkg.ErrEmptyStorage {
			code = codes.NotFound
		} else if err == multiPkg.ErrTimeOut {
			code = codes.Internal
		}
		err := status.Error(code, err.Error())
		journal.LogError(err)
		return &pb.ActivityListResponse{}, err
	}

	result := make([]*pb.ActivityListResponse_Activity, 0, len(activities))
	for key, act := range activities {
		d := act.BeginDate
		beginDate := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
		d = act.EndDate
		endDate := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())

		result = append(result, &pb.ActivityListResponse_Activity{
			Name:            key,
			BeginDate:       beginDate,
			EndDate:         endDate,
			TimesPerDay:     uint32(act.TimesPerDay),
			QuantityPerTime: act.QuantityPerTime,
		})
	}

	counter.SuccessRequests.Increase()
	return &pb.ActivityListResponse{
		Activities: result,
	}, nil
}

func (imp *implementation) ActivityListStream(in *pb.ActivityListStreamRequest, stream pb.Admin_ActivityListStreamServer) error {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return err
	}

	ch := make(chan models.DailyActivityRec)
	chErr := make(chan error, 1)
	ctx, _ := context.WithCancel(stream.Context())
	go imp.activity.ListStream(ctx, chatID, ch, chErr)

	act := models.DailyActivityRec{}
	for {
		select {
		case act = <-ch:
			d := act.BeginDate
			beginDate := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
			d = act.EndDate
			endDate := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())

			if err := stream.Send(&pb.ActivityListStreamResponse{
				Name:            act.Name,
				BeginDate:       beginDate,
				EndDate:         endDate,
				TimesPerDay:     uint32(act.TimesPerDay),
				QuantityPerTime: act.QuantityPerTime,
			}); err != nil {
				counter.ErrorRequests.Increase()
				journal.LogError(err)
				return err
			}
		case err := <-chErr:
			counter.ErrorRequests.Increase()
			code := codes.InvalidArgument
			if err == multiPkg.ErrEmptyStorage {
				code = codes.NotFound
			} else if err == multiPkg.ErrTimeOut {
				code = codes.Internal
			} else if err == multiPkg.ErrOk {
				counter.SuccessRequests.Increase()
				return nil
			}
			err = status.Error(code, err.Error())
			journal.LogError(err)
			return err
		case <-ctx.Done():
			counter.SuccessRequests.Increase()
			return nil
		}
	}
}

func (imp *implementation) ActivityToday(ctx context.Context, in *pb.ActivityTodayRequest) (*pb.ActivityTodayResponse, error) {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return &pb.ActivityTodayResponse{}, err
	}

	activities, err := imp.activity.Today(ctx, chatID)
	if err != nil {
		counter.ErrorRequests.Increase()
		code := codes.InvalidArgument
		if err == multiPkg.ErrNoActivitiesForToday {
			code = codes.NotFound
		} else if err == multiPkg.ErrTimeOut {
			code = codes.Internal
		}
		err := status.Error(code, err.Error())
		journal.LogError(err)
		return &pb.ActivityTodayResponse{}, err
	}

	result := make([]*pb.ActivityTodayResponse_Activity, 0, len(activities))
	for key, act := range activities {
		d := act.BeginDate
		beginDate := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
		d = act.EndDate
		endDate := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())

		result = append(result, &pb.ActivityTodayResponse_Activity{
			Name:            key,
			BeginDate:       beginDate,
			EndDate:         endDate,
			TimesPerDay:     uint32(act.TimesPerDay),
			QuantityPerTime: act.QuantityPerTime,
		})
	}

	counter.SuccessRequests.Increase()
	return &pb.ActivityTodayResponse{
		Activities: result,
	}, nil
}

func (imp *implementation) ActivityGet(ctx context.Context, in *pb.ActivityGetRequest) (*pb.ActivityGetResponse, error) {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return nil, err
	}
	actName := in.GetName()
	if len(actName) == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "name should be filled")
		journal.LogError(err)
		return nil, err
	}
	activity, err := imp.activity.Get(ctx, actName, chatID)
	if err != nil {
		counter.ErrorRequests.Increase()
		code := codes.InvalidArgument
		if err == multiPkg.ErrNotExist {
			code = codes.NotFound
		} else if err == multiPkg.ErrTimeOut {
			code = codes.Internal
		}
		err := status.Error(code, err.Error())
		journal.LogError(err)
		return nil, err
	}

	d := activity.BeginDate
	beginDate := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
	d = activity.EndDate
	endDate := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())

	counter.SuccessRequests.Increase()
	return &pb.ActivityGetResponse{
		Name:            actName,
		BeginDate:       beginDate,
		EndDate:         endDate,
		TimesPerDay:     uint32(activity.TimesPerDay),
		QuantityPerTime: activity.QuantityPerTime,
	}, nil
}

func (imp *implementation) ActivityUpdate(ctx context.Context, in *pb.ActivityUpdateRequest) (*pb.ActivityUpdateResponse, error) {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return nil, err
	}
	actName := in.GetName()
	if len(actName) == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "name should be filled")
		journal.LogError(err)
		return nil, err
	}

	beginDate, err := commandPkg.ParseDate(in.GetBeginDate())
	if err != nil {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, err.Error())
		journal.LogError(err)
		return nil, err
	}

	endDate, err := commandPkg.ParseDate(in.GetEndDate())
	if err != nil {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, err.Error())
		journal.LogError(err)
		return nil, err
	}

	activity, err := imp.activity.Update(ctx, actName, models.DailyActivity{
		BeginDate:       beginDate,
		EndDate:         endDate,
		TimesPerDay:     uint8(in.GetTimesPerDay()),
		QuantityPerTime: in.GetQuantityPerTime(),
	}, chatID)
	if err != nil {
		counter.ErrorRequests.Increase()
		code := codes.InvalidArgument
		if err == multiPkg.ErrNotExist {
			code = codes.NotFound
		} else if err == multiPkg.ErrTimeOut {
			code = codes.Internal
		}
		err := status.Error(code, err.Error())
		journal.LogError(err)
		return nil, err
	}

	d := activity.BeginDate
	beginDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())
	d = activity.EndDate
	endDateR := fmt.Sprintf("%02d.%02d.%04d", d.Day(), d.Month(), d.Year())

	counter.SuccessRequests.Increase()
	return &pb.ActivityUpdateResponse{
		Name:            actName,
		BeginDate:       beginDateR,
		EndDate:         endDateR,
		TimesPerDay:     uint32(activity.TimesPerDay),
		QuantityPerTime: activity.QuantityPerTime,
	}, nil
}

func (imp *implementation) ActivityDelete(ctx context.Context, in *pb.ActivityDeleteRequest) (*pb.ActivityDeleteResponse, error) {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return nil, err
	}

	if err := imp.activity.Delete(ctx, in.GetName(), chatID); err != nil {
		counter.ErrorRequests.Increase()
		code := codes.InvalidArgument
		if err == multiPkg.ErrNotExist {
			code = codes.NotFound
		} else if err == multiPkg.ErrTimeOut {
			code = codes.Internal
		}
		err := status.Error(code, err.Error())
		journal.LogError(err)
		return nil, err
	}

	counter.SuccessRequests.Increase()
	return &pb.ActivityDeleteResponse{}, nil
}
