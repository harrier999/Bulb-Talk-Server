package redis

import (
	"context"
	"encoding/json"
	"server/internal/models/message"
	"server/internal/repository"
	"strconv"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	maxStreamLength = 5000
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

	streamKey := "stream:room:" + roomID + ":messages"

	var msgID uuid.UUID

	if baseMsg, ok := msg.(*message.BaseMessage); ok {
		msgID = baseMsg.Id
	}

	if msgID == uuid.Nil {
		msgID, _ = uuid.NewV7()
	}

	_, err = r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		ID:     msgID.String(),
		Values: map[string]interface{}{
			"message": string(msgJSON),
		},
	}).Result()

	if err != nil {
		return err
	}

	r.client.XTrimMaxLen(ctx, streamKey, maxStreamLength)

	return nil
}

func (r *RedisMessageRepository) GetMessages(ctx context.Context, roomID string, lastMessageID int64) ([]message.Message, error) {
	streamKey := "stream:room:" + roomID + ":messages"

	var start string
	if lastMessageID <= 0 {
		start = "-"
	} else {
		start = "(" + strconv.FormatInt(lastMessageID, 10)
	}

	streams, err := r.client.XRange(ctx, streamKey, start, "+").Result()
	if err != nil {
		return nil, err
	}

	messages, err := r.redisStreamToMessageList(streams)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *RedisMessageRepository) redisStreamToMessageList(streams []redis.XMessage) ([]message.Message, error) {
	messages := make([]message.Message, 0, len(streams))

	for _, stream := range streams {
		msgStr, ok := stream.Values["message"].(string)
		if !ok {
			continue
		}

		var baseMsg message.BaseMessage
		err := json.Unmarshal([]byte(msgStr), &baseMsg)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &baseMsg)
	}

	return messages, nil
}

func (r *RedisMessageRepository) GetMessagesByUUID(ctx context.Context, roomID string, lastMessageUUID uuid.UUID) ([]message.Message, error) {
	streamKey := "stream:room:" + roomID + ":messages"

	var start string
	if lastMessageUUID == uuid.Nil {
		start = "-"
	} else {
		start = "(" + lastMessageUUID.String()
	}

	streams, err := r.client.XRange(ctx, streamKey, start, "+").Result()
	if err != nil {
		return nil, err
	}

	messages, err := r.redisStreamToMessageList(streams)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
