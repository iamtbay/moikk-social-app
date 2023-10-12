package mytypes

import "time"

type NewPost struct {
	UserID  string   `json:"user_id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Files   []string `json:"files"`
}

type UpdatePost struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Files   []string `json:"files"`
}

type GetPost struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Files     []string  `json:"files"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetPostWithUsername struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Files     []string  `json:"files"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
