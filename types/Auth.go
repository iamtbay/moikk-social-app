package mytypes

type Register struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Location string `json:"location"`
	Username string `json:"username"`
}
type UserInfos struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Location     string `json:"location"`
	Username     string `json:"username"`
	ProfilePhoto string `json:"profile_photo"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	ProfilePhoto string `json:"profile_photo"`
}
