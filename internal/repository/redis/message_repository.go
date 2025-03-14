package redis

import (
	"context"
	"encoding/json"
	"server/internal/models/message"
	"server/internal/repository"

	"github.com/redis/go-redis/v9"
)

type RedisMessageRepository struct {
	client *redis.Client
}

func NewRedisMessageRepository(client *redis.Client) repository.MessageRepository {
	return &RedisMessageRepository{
		client: client,
	}
}

func (r *RedisMessageRepository) SaveMessage(ctx context.Context, roomID string, msg message.Message) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = r.client.RPush(ctx, "room:"+roomID+":messages", string(msgJSON)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisMessageRepository) GetMessages(ctx context.Context, roomID string, lastMessageID int64) ([]message.Message, error) {

	redisKey := "room:" + roomID + ":messages"

	var redisList []string
	var err error

	if lastMessageID <= 0 {
		redisList, err = r.client.LRange(ctx, redisKey, 0, -1).Result()
	} else {

		redisList, err = r.client.LRange(ctx, redisKey, lastMessageID, -1).Result()
	}

	if err != nil {
		return nil, err
	}

	messages, err := r.redisListToMessageList(redisList)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *RedisMessageRepository) redisListToMessageList(redisList []string) ([]message.Message, error) {
	messages := make([]message.Message, 0, len(redisList))

	for _, item := range redisList {
		var msg message.Message
		err := json.Unmarshal([]byte(item), &msg)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
