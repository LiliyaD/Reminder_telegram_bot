package common

import (
	"log"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal/counter"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

var syncProducer sarama.SyncProducer

func init() {
	brokers := []string{"localhost:9095", "localhost:9096", "localhost:9097"}
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	var err error
	syncProducer, err = sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Sync Kafka"))
	}
}

func SendToKafka(uid, topic string, value []byte) error {
	counter.OutputRequests.Increase()
	_, _, err := syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(uid),
		Value: sarama.ByteEncoder(value),
	})
	if err != nil {
		counter.FailedRequests.Increase()
		journal.LogError(errors.Wrap(err, "error after sendToKafka in producer.go"))
		return err
	} else {
		counter.SuccessRequests.Increase()
		return nil
	}
}
