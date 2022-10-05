package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"os"
)

type Message struct {
	Sender     string `json:"sender"`
	Receiver   string `json:"receiver"`
	Message    string `json:"message"`
	DateCrated int64  `json:"date_crated"`
}

type Redis struct {
	ctx context.Context
	db  *redis.Client
}

func New() *Redis {
	r := Redis{}

	r.ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	r.db = rdb
	return &r
}

func (r *Redis) DelUser(sellyId string) {
	r.db.Del(r.ctx, sellyId)
}

func (r *Redis) GetMessages(sellyId string) []Message {
	res, _ := r.db.LRange(r.ctx, sellyId, 0, -1).Result()

	messages := []Message{}
	var message Message

	for _, m := range res {
		json.Unmarshal([]byte(m), &message)

		messages = append(messages, message)
	}

	return messages
}
