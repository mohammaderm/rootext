package postrepositoryredis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mohammaderm/rootext/entity"
	"github.com/mohammaderm/rootext/repository/redis"

	go_redis "github.com/redis/go-redis/v9"
)

type Repository struct {
	conn *redis.Redis
}

type PostRepoCache interface {
	SavePosts(ctx context.Context, posts []entity.Post, interval, sortBy string) error
	GetTopPost(ctx context.Context, interval, sortBy string) ([]entity.Post, error)
	UpdatePost(ctx context.Context, post entity.Post) (entity.Post, error)
	UpdateVote(ctx context.Context, id uint, score int) error
}

const limit = 5

func New(conn *redis.Redis) PostRepoCache {
	return &Repository{
		conn: conn,
	}
}

func (r Repository) UpdatePost(ctx context.Context, post entity.Post) (entity.Post, error) {
	pipe := r.conn.Conn().Pipeline()
	now := time.Now()
	pipe.HSet(ctx, fmt.Sprintf("post:%d", post.ID), "title", post.Title)
	pipe.HSet(ctx, fmt.Sprintf("post:%d", post.ID), "content", post.Content)
	pipe.HSet(ctx, fmt.Sprintf("post:%d", post.ID), "score", post.Score)
	pipe.HSet(ctx, fmt.Sprintf("post:%d", post.ID), "updated_at", now)
	post.UpdatedAt = now

	_, err := pipe.Exec(ctx)
	if err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (r Repository) GetTopPost(ctx context.Context, interval, sortBy string) ([]entity.Post, error) {

	postIDs, err := r.conn.Conn().ZRevRange(ctx, fmt.Sprintf("top:%s:%s", interval, sortBy), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	if len(postIDs) == 0 {
		return nil, nil
	}
	var posts []entity.Post
	for _, postID := range postIDs {
		result, err := r.conn.Conn().HGetAll(ctx, postID).Result()
		if err != nil {
			return nil, err
		}

		intID, _ := strconv.Atoi(result["id"])
		intUserID, _ := strconv.Atoi(result["user_id"])
		intScore, _ := strconv.Atoi(result["score"])
		stringCreatedAt, _ := time.Parse("2006-01-02T15:04:05Z", result["created_at"])
		stringUpdatedAt, _ := time.Parse("2006-01-02T15:04:05Z", result["updated_at"])

		post := entity.Post{
			ID:        uint(intID),
			UserID:    uint(intUserID),
			Title:     result["title"],
			Content:   result["content"],
			Score:     intScore,
			CreatedAt: stringCreatedAt,
			UpdatedAt: stringUpdatedAt,
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (r Repository) SavePosts(ctx context.Context, posts []entity.Post, interval, sortBy string) error {

	pipe := r.conn.Conn().Pipeline()
	for _, post := range posts {
		pipe.HSet(ctx, fmt.Sprintf("post:%d", post.ID), map[string]interface{}{
			"id":         post.ID,
			"user_id":    post.UserID,
			"title":      post.Title,
			"content":    post.Content,
			"score":      post.Score,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		})

		score := float64(post.Score)
		if sortBy == "created_at" {
			score = float64(post.CreatedAt.Unix())
		}

		pipe.ZAdd(ctx, fmt.Sprintf("top:%s:%s", interval, sortBy), go_redis.Z{
			Score:  score,
			Member: fmt.Sprintf("post:%d", post.ID),
		})
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateVote(ctx context.Context, id uint, score int) error {

	intervals := []string{"1 day", "7 days", "30 days"}
	const sortBy = "score"

	pipe := r.conn.Conn().Pipeline()

	pipe.HSet(ctx, fmt.Sprintf("post:%d", id), "score", score)

	for _, interval := range intervals {
		key := fmt.Sprintf("top:%s:%s", interval, sortBy)
		exists, err := r.conn.Conn().Exists(ctx, key).Result()
		if err != nil {
			return err
		}

		if exists == 0 {
			fmt.Printf("Sorted set %s does not exist, skipping update.\n", key)
			continue
		}

		pipe.ZAdd(ctx, key, go_redis.Z{
			Score:  float64(score),
			Member: fmt.Sprintf("post:%d", id),
		})
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
