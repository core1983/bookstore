package order

import "time"

type Order struct {
	ID         string        `json:"id"`
	CreatedAt  time.Time     `json:"created_at"`
	TotalPrice float64       `json:"total_price"`
	UserID     string        `json:"user_id"`
	Books      []OrderedBook `json:"books"`
}

type OrderedBook struct {
	ID       string  `json:"id"`
	ISBN     string  `json:"isbn"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Price    float64 `json:"price"`
	Quantity uint32  `json:"quantity"`
}

type User struct {
	ID        string    `json:"id"`
	UserName  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
type Book struct {
	ID     string  `json:"id"`
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}
