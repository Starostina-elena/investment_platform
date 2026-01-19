package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Starostina-elena/investment_platform/services/comment/core"
)

type Cache struct {
	client *redis.Client
	log    slog.Logger
	ttl    time.Duration
}

func NewCache(client *redis.Client, log slog.Logger) *Cache {
	return &Cache{
		client: client,
		log:    log,
		ttl:    10 * time.Minute,
	}
}

func (c *Cache) SetComment(ctx context.Context, comment *core.Comment) error {
	key := "comment:" + strconv.Itoa(comment.ID)
	data, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *Cache) GetComment(ctx context.Context, id int) (*core.Comment, error) {
	key := "comment:" + strconv.Itoa(id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var comment core.Comment
	if err := json.Unmarshal([]byte(val), &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *Cache) DeleteComment(ctx context.Context, id int) error {
	key := "comment:" + strconv.Itoa(id)
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) SetProjectComments(ctx context.Context, projectID, limit, offset int, comments []core.Comment) error {
	key := "project_comments:" + strconv.Itoa(projectID) + ":" + strconv.Itoa(limit) + ":" + strconv.Itoa(offset)
	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *Cache) GetProjectComments(ctx context.Context, projectID, limit, offset int) ([]core.Comment, error) {
	key := "project_comments:" + strconv.Itoa(projectID) + ":" + strconv.Itoa(limit) + ":" + strconv.Itoa(offset)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var comments []core.Comment
	if err := json.Unmarshal([]byte(val), &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Cache) InvalidateProjectComments(ctx context.Context, projectID int) error {
	pattern := "project_comments:" + strconv.Itoa(projectID) + ":*"
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		c.client.Del(ctx, iter.Val())
	}
	return iter.Err()
}
