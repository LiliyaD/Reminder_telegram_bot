package main

import (
	"encoding/json"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal/counter"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/pkg/core/daily_activity/models"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type Consumer struct {
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	journal.LogInfo("Begin method ConsumeClaim")
	for {
		select {
		case <-session.Context().Done():
			journal.LogInfo("Done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				journal.LogWarn("Data channel closed")
				return nil
			}

			key := string(msg.Key)

			switch msg.Topic {
			case "Create":
				counter.InputRequests.Increase()
				var act *models.DailyActivityCreationReq
				err := json.Unmarshal(msg.Value, &act)
				if err != nil {
					counter.ErrorRequests.Increase()
					journal.LogError(errors.Wrap(err, "Consumer Unmarshal error"))
					continue
				}

				actC, err := activity.Add(context.Background(), act.Name, act.DailyActivity, act.Chat)
				if err != nil {
					counter.ErrorRequests.Increase()
					journal.LogError(err)

					if cacheAnsw != nil {
						cacheAnsw.Set(key, models.DailyActivityAnsw{Error: err.Error()}, 10*time.Second)
					}

				} else {
					counter.SuccessRequests.Increase()

					if cacheAnsw != nil {
						cacheAnsw.Set(key, models.DailyActivityAnsw{
							ChatID: act.Chat.ChatID,
							Name:   act.Name,
							DailyActivity: models.DailyActivity{
								BeginDate:       actC.BeginDate,
								EndDate:         actC.EndDate,
								TimesPerDay:     actC.TimesPerDay,
								QuantityPerTime: actC.QuantityPerTime}}, 10*time.Second)
					}
				}

			case "Update":
				counter.InputRequests.Increase()
				var act *models.DailyActivityUpdateReq
				err := json.Unmarshal(msg.Value, &act)
				if err != nil {
					counter.ErrorRequests.Increase()
					journal.LogError(errors.Wrap(err, "Consumer Unmarshal error"))
					continue
				}

				actU, err := activity.Update(context.Background(), act.Name, act.DailyActivity, act.ChatID)
				if err != nil {
					counter.ErrorRequests.Increase()
					journal.LogError(err)

					if cacheAnsw != nil {
						cacheAnsw.Set(key, models.DailyActivityAnsw{Error: err.Error()}, 10*time.Second)
					}
				} else {
					counter.SuccessRequests.Increase()

					if cacheAnsw != nil {
						cacheAnsw.Set(key, models.DailyActivityAnsw{
							ChatID: act.ChatID,
							Name:   act.Name,
							DailyActivity: models.DailyActivity{
								BeginDate:       actU.BeginDate,
								EndDate:         actU.EndDate,
								TimesPerDay:     actU.TimesPerDay,
								QuantityPerTime: actU.QuantityPerTime}}, 10*time.Second)
					}
				}

			case "Delete":
				counter.InputRequests.Increase()
				var act *models.DailyActivityDeletionReq
				err := json.Unmarshal(msg.Value, &act)
				if err != nil {
					counter.ErrorRequests.Increase()
					journal.LogError(errors.Wrap(err, "Consumer Unmarshal error"))
					continue
				}

				err = activity.Delete(context.Background(), act.Name, act.ChatID)
				if err != nil {
					counter.ErrorRequests.Increase()
					journal.LogError(err)

					if cacheAnsw != nil {
						cacheAnsw.Publish(key, models.DailyActivityDelAnsw{Error: err.Error()})
					}
				} else {
					counter.SuccessRequests.Increase()

					if cacheAnsw != nil {
						cacheAnsw.Publish(key, models.DailyActivityDelAnsw{})
					}
				}

			default:
				journal.LogError("Unknown topic")
				continue
			}
		}
	}
}

func consume() {
	brokers := []string{"localhost:9095", "localhost:9096", "localhost:9097"}
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Version = sarama.V0_10_2_0

	client, err := sarama.NewConsumerGroup(brokers, "startConsuming", cfg)
	if err != nil {
		journal.LogFatal(err)
	}

	ctx := context.Background()
	consumer := &Consumer{}
	for {
		if err := client.Consume(ctx, []string{"Create", "Update", "Delete"}, consumer); err != nil {
			journal.LogErrorf("on consume: %v", err)
			time.Sleep(time.Second * 10)
		}
	}
}
