package core

type User struct {
	ID       string `json:"-"`
	UserName string `json:"user_name"`
}
