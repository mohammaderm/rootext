package postService

import (
	"context"
	"fmt"

	"github.com/mohammaderm/rootext/entity"
	"github.com/mohammaderm/rootext/params"
	postrepository "github.com/mohammaderm/rootext/repository/postgres/postRepository"
	postrepositoryredis "github.com/mohammaderm/rootext/repository/redis/postRepositoryRedis"
)

type PostService struct {
	repo  postrepository.PostRepo
	cache postrepositoryredis.PostRepoCache
}

func New(postRepo postrepository.PostRepo, cache postrepositoryredis.PostRepoCache) PostService {
	return PostService{
		repo:  postRepo,
		cache: cache,
	}
}

// Operations related to user's own posts

func (s PostService) Create(ctx context.Context, req params.CreatePostReq, userID uint) (params.CreatePostRes, error) {
	post := entity.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}
	createdPost, err := s.repo.Create(ctx, post)
	if err != nil {
		return params.CreatePostRes{}, err
	}
	return params.CreatePostRes{
		Post: createdPost,
	}, nil
}

func (s PostService) Delete(ctx context.Context, id, userID uint) error {
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return err
	}
	return nil
}

func (s PostService) Update(ctx context.Context, req params.UpdatePostReq, userID uint) (params.UpdatePostRes, error) {
	post := entity.Post{
		ID:      req.Id,
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	_, err := s.cache.UpdatePost(ctx, post)
	if err != nil {
		return params.UpdatePostRes{}, err
	}

	updatedPost, err := s.repo.Update(ctx, post)
	if err != nil {
		return params.UpdatePostRes{}, err
	}
	return params.UpdatePostRes{
		Post: updatedPost,
	}, nil
}

func (s PostService) GetAll(ctx context.Context, userID uint) (params.GetAllPostRes, error) {
	posts, err := s.repo.GetAll(ctx, userID)
	if err != nil {
		return params.GetAllPostRes{}, err
	}
	return params.GetAllPostRes{
		Posts: posts,
	}, nil
}

func (s PostService) GetById(ctx context.Context, id, userID uint) (params.GetByIdRes, error) {
	post, err := s.repo.GetById(ctx, id, userID)
	if err != nil {
		return params.GetByIdRes{}, err
	}
	return params.GetByIdRes{
		Post: post,
	}, nil
}

// Operations related to all posts

func (s PostService) VotePost(ctx context.Context, req params.VotePostReq, userID uint) error {
	score, err := s.repo.VotePost(ctx, userID, req.Id, req.Vote)
	if err != nil {
		return err
	}
	fmt.Println(score)
	if err := s.cache.UpdateVote(ctx, req.Id, score); err != nil {
		fmt.Println(err)
	}

	return nil
}

func (s PostService) GetSortedPost(ctx context.Context, req params.GetSortedPostReq) (params.GetSortedPostRes, error) {

	posts, err := s.cache.GetTopPost(ctx, req.Interval, req.SortBy)
	if err != nil {
		return params.GetSortedPostRes{}, err
	}
	if posts == nil {
		posts, err := s.repo.GetSorted(ctx, req.Interval, req.SortBy)
		if err != nil {
			return params.GetSortedPostRes{}, err
		}
		return params.GetSortedPostRes{
			Posts: posts,
		}, nil
	}
	return params.GetSortedPostRes{
		Posts: posts,
	}, nil

}
