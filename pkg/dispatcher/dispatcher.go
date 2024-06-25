package dispatcher

import (
	"context"
	"cyclic/pkg/colonel"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var Dispatcher *redis.Client

var ctx = context.Background()

var (
	Verify = "verify"
	Notify = "notify"
)

type Message struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

func Init() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", colonel.Writ.Redis.Host, colonel.Writ.Redis.Port),
		Password: colonel.Writ.Redis.Password,
		DB:       colonel.Writ.Redis.DB,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return err
	}

	Dispatcher = rdb

	return nil
}

func Enqueue(message *Message) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = Dispatcher.LPush(ctx, "cyclic-mailer", b).Result()

	if err != nil {
		return err
	}

	return nil
}

func Dequeue() (*Message, error) {
	result, err := Dispatcher.BLPop(ctx, 0, "cyclic-mailer").Result()
	if err != nil {
		return nil, err
	}

	var message Message

	if err := json.Unmarshal([]byte(result[1]), &message); err != nil {
		return nil, err
	}

	return &message, nil
}
