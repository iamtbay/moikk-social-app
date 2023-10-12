package mytypes

import "time"

type GetFollows struct {
	ID         string    `json:"id"`
	FollowerID string    `json:"follower_user_id"`
	FollowedID string    `json:"followed_user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetFollowers struct {
	Username string `json:"username"`
}
