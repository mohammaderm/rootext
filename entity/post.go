package entity

import "time"

type (
	Post struct {
		ID        uint      `db:"id"`
		UserID    uint      `db:"user_id"`
		Title     string    `db:"title"`
		Content   string    `db:"content"`
		Score     int       `db:"score"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	Vote struct {
		ID        int       `db:"id"`
		UserID    int       `db:"user_id"`
		PostID    int       `db:"post_id"`
		VoteValue int       `db:"vote_value"`
		CreatedAt time.Time `db:"created_at"`
	}
)
