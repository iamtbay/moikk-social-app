package mytypes

import "time"

type GetLikes struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Username    string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
