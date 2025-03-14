package postrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mohammaderm/rootext/entity"
	"github.com/mohammaderm/rootext/repository/postgres"
	postrepositoryredis "github.com/mohammaderm/rootext/repository/redis/postRepositoryRedis"
)

type PostRepo interface {
	Create(ctx context.Context, post entity.Post) (entity.Post, error)
	Delete(ctx context.Context, postID uint, userID uint) error
	Update(ctx context.Context, post entity.Post) (entity.Post, error)
	GetAll(ctx context.Context, userID uint) ([]entity.Post, error)
	GetById(ctx context.Context, postID uint, userID uint) (entity.Post, error)

	VotePost(ctx context.Context, userID, postID uint, vote int) (int, error)
	GetSorted(ctx context.Context, interval, sortBy string) ([]entity.Post, error)
}

const limit = 5

type Repository struct {
	conn  *postgres.PostgresDB
	cache postrepositoryredis.PostRepoCache
}

func New(conn *postgres.PostgresDB, cache postrepositoryredis.PostRepoCache) PostRepo {
	return &Repository{
		conn:  conn,
		cache: cache,
	}
}

func (r Repository) Create(ctx context.Context, post entity.Post) (entity.Post, error) {
	result, err := r.conn.Conn().ExecContext(ctx, "insert into posts (user_id, title, content) values ($1, $2, $3);", post.UserID, post.Title, post.Content)
	if err != nil {
		return entity.Post{}, err
	}
	id, _ := result.LastInsertId()
	post.ID = uint(id)
	return post, nil
}

func (r Repository) Delete(ctx context.Context, postID uint, userID uint) error {
	result, err := r.conn.Conn().ExecContext(ctx, "delete from posts where id = $1 and user_id = $2;", postID, userID)
	if err != nil {
		return err
	}
	if effected, err := result.RowsAffected(); effected == 0 || err != nil {
		return errors.New("no posts found with this ID")
	}
	return nil
}

func (r Repository) Update(ctx context.Context, post entity.Post) (entity.Post, error) {
	result, err := r.conn.Conn().ExecContext(ctx, "update posts set title = $1, content = $2, updated_at = now() where id = $3 and user_id = $4;", post.Title, post.Content, post.ID, post.UserID)
	if err != nil {
		return entity.Post{}, err
	}
	if rowEfected, err := result.RowsAffected(); rowEfected == 0 || err != nil {
		return entity.Post{}, errors.New("no posts found with this ID")
	}
	return post, nil
}

func (r Repository) GetAll(ctx context.Context, userID uint) ([]entity.Post, error) {
	var posts []entity.Post
	err := r.conn.Conn().SelectContext(ctx, &posts, "select * from posts where user_id = $1;", userID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r Repository) GetById(ctx context.Context, postID uint, userID uint) (entity.Post, error) {
	var post entity.Post
	err := r.conn.Conn().GetContext(ctx, &post, "select * from posts where id = $1 and user_id = $2;", postID, userID)
	if err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (r Repository) VotePost(ctx context.Context, userID, postID uint, vote int) (int, error) {
	_, err := r.conn.Conn().ExecContext(ctx, "INSERT INTO votes (user_id, post_id, vote_value) VALUES ($1, $2, $3) ON CONFLICT (user_id, post_id) DO UPDATE SET vote_value = EXCLUDED.vote_value", userID, postID, vote)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	var newScore int
	err = r.conn.Conn().QueryRowContext(ctx, "UPDATE posts SET score = (SELECT COALESCE(SUM(vote_value), 0) FROM votes WHERE post_id = $1) WHERE id = $1 RETURNING score", postID).Scan(&newScore)
	if err != nil {
		return 0, err
	}
	return newScore, nil
}

func (r Repository) GetSorted(ctx context.Context, interval, sortBy string) ([]entity.Post, error) {
	var posts []entity.Post

	query := fmt.Sprintf(`SELECT * FROM posts WHERE created_at >= NOW() - INTERVAL '%s' ORDER BY %s DESC LIMIT %d`, interval, sortBy, limit)
	if err := r.conn.Conn().SelectContext(ctx, &posts, query); err != nil {
		return nil, err
	}

	if err := r.cache.SavePosts(ctx, posts, interval, sortBy); err != nil {
		return nil, err
	}

	return posts, nil
}
