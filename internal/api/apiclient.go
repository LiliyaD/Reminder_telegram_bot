package api

import (
	"context"
	"io"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal/counter"
	commandPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/bot/command"
	commonPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/common/commands"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewClient(client pb.AdminClient) pb.AdminServer {
	return &implementationClient{
		client: client,
	}
}

type implementationClient struct {
	pb.UnimplementedAdminServer
	client pb.AdminClient
}

func (imp *implementationClient) ActivityCreate(ctx context.Context, in *pb.ActivityCreateRequest) (*pb.ActivityCreateResponse, error) {
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

	response, err := commonPkg.Create(ctx, actName,
		models.DailyActivity{
			BeginDate:       beginDate,
			EndDate:         endDate,
			TimesPerDay:     uint8(in.GetTimesPerDay()),
			QuantityPerTime: in.GetQuantityPerTime()},
		models.Chat{ChatID: chatID, UserName: in.GetUserName()}, imp.client)
	if err != nil {
		counter.ErrorRequests.Increase()
		journal.LogError(err)
		return nil, err
	}

	return response, nil
}

func (imp *implementationClient) ActivityList(ctx context.Context, in *pb.ActivityListRequest) (*pb.ActivityListResponse, error) {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return nil, err
	}

	response, err := commonPkg.List(ctx, chatID, models.Pagination{Limit: in.GetLimit(), Offset: in.GetOffset(), Order: in.GetOrder()}, imp.client)
	if err != nil {
		counter.ErrorRequests.Increase()
		journal.LogError(err)
		return nil, err
	}

	return response, nil
}

func (imp *implementationClient) ActivityListStream(in *pb.ActivityListStreamRequest, stream pb.Admin_ActivityListStreamServer) error {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return err
	}

	act, err := commonPkg.ListStream(stream.Context(), chatID, imp.client)
	if err != nil {
		counter.ErrorRequests.Increase()
		journal.LogError(err)
		return err
	}

	for {
		v, err := act.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			counter.ErrorRequests.Increase()
			return err
		}
		if err := stream.Send(&pb.ActivityListStreamResponse{
			Name:            v.GetName(),
			BeginDate:       v.GetBeginDate(),
			EndDate:         v.GetEndDate(),
			TimesPerDay:     v.GetTimesPerDay(),
			QuantityPerTime: v.GetQuantityPerTime(),
		}); err != nil {
			journal.LogError(err)
			return err
		}
	}

	return nil
}

func (imp *implementationClient) ActivityToday(ctx context.Context, in *pb.ActivityTodayRequest) (*pb.ActivityTodayResponse, error) {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return nil, err
	}

	response, err := commonPkg.Today(ctx, chatID, imp.client)
	if err != nil {
		counter.ErrorRequests.Increase()
		journal.LogError(err)
		return nil, err
	}

	return response, nil
}

func (imp *implementationClient) ActivityGet(ctx context.Context, in *pb.ActivityGetRequest) (*pb.ActivityGetResponse, error) {
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

	response, err := commonPkg.Get(ctx, chatID, actName, imp.client)
	if err != nil {
		counter.ErrorRequests.Increase()
		journal.LogError(err)
		return nil, err
	}

	return response, nil
}

func (imp *implementationClient) ActivityUpdate(ctx context.Context, in *pb.ActivityUpdateRequest) (*pb.ActivityUpdateResponse, error) {
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

	response, err := commonPkg.Update(ctx, chatID, actName,
		models.DailyActivity{
			BeginDate:       beginDate,
			EndDate:         endDate,
			TimesPerDay:     uint8(in.GetTimesPerDay()),
			QuantityPerTime: in.GetQuantityPerTime()}, imp.client)
	if err != nil {
		counter.ErrorRequests.Increase()
		journal.LogError(err)
		return nil, err
	}

	return response, nil
}

func (imp *implementationClient) ActivityDelete(ctx context.Context, in *pb.ActivityDeleteRequest) (*pb.ActivityDeleteResponse, error) {
	counter.InputRequests.Increase()

	chatID := in.GetChatID()
	if chatID == 0 {
		counter.ErrorRequests.Increase()
		err := status.Error(codes.InvalidArgument, "chatID should be filled")
		journal.LogError(err)
		return nil, err
	}

	err := commonPkg.Delete(ctx, chatID, in.GetName(), imp.client)
	if err != nil {
		counter.ErrorRequests.Increase()
		journal.LogError(err)
		return nil, err
	}

	return &pb.ActivityDeleteResponse{}, nil
}
