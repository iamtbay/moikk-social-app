package mytypes

import (
	"database/sql"
	"time"
)

type NewComment struct {
	Content string `json:"content"`
}

type GetComment struct {
	ID             string         `json:"id"`
	PostID         string         `json:"post_id"`
	UpperCommentID sql.NullString `json:"upper_comment_id"`
	UserID         string         `json:"user_id"`
	Content        string         `json:"content"`
	CreatedAt      time.Time      `json:"created_at"`
}
type GetCommentWithUsername struct {
	ID             string         `json:"id"`
	PostID         string         `json:"post_id"`
	UpperCommentID sql.NullString `json:"upper_comment_id"`
	UserID         string         `json:"user_id"`
	Username       string         `json:"username"`
	Content        string         `json:"content"`
	CreatedAt      time.Time      `json:"created_at"`
}
