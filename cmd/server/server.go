package main

import (
	"context"
	"net"
	"net/http"

	apiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/api"
	configPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/config"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	dailyActPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runGRPCServer(act dailyActPkg.Interface) {
	journal.LogInfo("RUN Grpc Server")

	listener, err := net.Listen("tcp", configPkg.ServerAddress)
	if err != nil {
		journal.LogFatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.New(act))

	if err = grpcServer.Serve(listener); err != nil {
		journal.LogFatal(err)
	}
}

func runREST() {
	journal.LogInfo("RUN REST")

	ctx := context.Background()

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAdminHandlerFromEndpoint(ctx, grpcMux, configPkg.ServerAddress, opts); err != nil {
		journal.LogFatal(err)
	}

	// for swagger
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	fs := http.FileServer(http.Dir("./pkg/api/doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	if err := http.ListenAndServe(":8070", mux); err != nil {
		journal.LogFatal(err)
	}

}

func runCount() {
	//localhost:8060/debug/vars
	if err := http.ListenAndServe(":8060", nil); err != nil {
		journal.LogFatal(err)
	}
}
