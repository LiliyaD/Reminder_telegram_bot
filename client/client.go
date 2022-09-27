package main

import (
	"net"
	"net/http"

	apiPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/api"
	configPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/config"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	botPkg "github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/bot"
	pb "github.com/LiliyaD/Reminder_telegram_bot/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	journal.New("client", false)
	journal.LogInfo("RUN Client")

	// Set up a connection to the server
	conns, err := grpc.Dial(configPkg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		journal.LogFatal(err)
	}
	defer conns.Close()

	client := pb.NewAdminClient(conns)

	go runCount()

	go runBotService(client)

	runTechService(client)

}

func runCount() {
	//localhost:8050/debug/vars
	if err := http.ListenAndServe(":8050", nil); err != nil {
		journal.LogFatal(err)
	}
}

func runBotService(client pb.AdminClient) {
	bot := botPkg.MustNew(client)
	bot.InitCommands()
	if err := bot.Run(); err != nil {
		journal.LogFatal(err)
	}
}

func runTechService(client pb.AdminClient) {
	listener, err := net.Listen("tcp", configPkg.ClientAddress)
	if err != nil {
		journal.LogFatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.NewClient(client))

	if err = grpcServer.Serve(listener); err != nil {
		journal.LogFatal(err)
	}
}
