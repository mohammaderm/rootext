package params

import "rootext/entity"

type CreatePostReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreatePostRes struct {
	Post entity.Post `json:"post"`
}

type UpdatePostReq struct {
	Id      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type UpdatePostRes struct {
	Post entity.Post `json:"post"`
}

type GetAllPostRes struct {
	Posts []entity.Post `json:"posts"`
}

type GetByIdRes struct {
	Post entity.Post `json:"post"`
}

type VotePostReq struct {
	Id   uint `json:"id"`
	Vote int  `json:"vote"`
}

type GetSortedPostReq struct {
	Interval string `json:"interval"`
	SortBy   string `json:"sortBy"`
}

type GetSortedPostRes struct {
	Posts []entity.Post `json:"posts"`
}
