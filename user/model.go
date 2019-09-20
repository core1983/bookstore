package user

import "time"

type User struct {
	ID        string    `json:"id"`
	UserName  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
