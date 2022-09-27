package redis

import (
	"encoding/json"
	"time"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal/counter"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type InterfaceRedis interface {
	Get(key string, resp interface{}) error
	Set(key string, val interface{}, exp time.Duration)
	Del(key string)
	Publish(channel string, msg interface{})
	Subscribe(channel string, resp interface{})
}

type cache struct {
	client *redis.Client
}

func Connect() InterfaceRedis {
	c := &cache{client: redis.NewClient(&redis.Options{Addr: "localhost:6379"})}

	if c == nil {
		journal.LogError("Redis Server is unavailable")
	}

	return c
}

func (c *cache) Get(key string, resp interface{}) error {
	journal.LogInfo("Redis Get by key", key)
	getResp := c.client.Get(key)
	if err := getResp.Err(); err != nil {
		counter.MissCacheRequests.Increase()
		journal.LogError(errors.Wrap(err, "On Redis Get"))
		journal.LogErrorf("On Redis Key %v String value %v", key, getResp.Val())
		return err
	} else {
		counter.HitCacheRequests.Increase()
	}

	respBytes, err := getResp.Bytes()
	if err != nil {
		journal.LogError(errors.Wrap(err, "On Redis Bytes"))
		journal.LogErrorf("On Redis Key %v String value %v", key, getResp.Val())
		return err
	}

	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		journal.LogError(errors.Wrap(err, "Redis Unmarshal"))
		return err
	}

	return nil
}

func (c *cache) Set(key string, val interface{}, exp time.Duration) {
	journal.LogInfo("Redis Set by key", key)

	valBytes, err := json.Marshal(val)
	if err != nil {
		journal.LogError(errors.Wrap(err, "Redis Marshal"))
	}

	resp := c.client.Set(key, valBytes, exp)
	if err := resp.Err(); err != nil {
		journal.LogError(errors.Wrap(err, "On Set Redis"))
	}
}

func (c *cache) Del(key string) {
	journal.LogInfo("Redis Del by key", key)

	resp := c.client.Del(key)
	if err := resp.Err(); err != nil {
		journal.LogError(errors.Wrap(err, "On Del Redis"))
	}
}

func (c *cache) Publish(channel string, msg interface{}) {
	journal.LogInfo("Redis Publish channel", channel)

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		journal.LogError(errors.Wrap(err, "Redis Marshal"))
	}

	if err := c.client.Publish(channel, msgBytes).Err(); err != nil {
		journal.LogError(errors.Wrap(err, "On Publish Redis"))
	}
}

func (c *cache) Subscribe(channel string, resp interface{}) {
	journal.LogInfo("Redis Subscribe channel", channel)

	subscriber := c.client.Subscribe(channel)

	for {
		msg, err := subscriber.ReceiveMessage()
		if err != nil {
			journal.LogError(errors.Wrap(err, "On Subscribe Redis"))
			continue
		}

		err = json.Unmarshal([]byte(msg.Payload), resp)
		if err != nil {
			journal.LogError(errors.Wrap(err, "Redis Unmarshal"))
		}
		return
	}

}
