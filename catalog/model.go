package catalog

type Book struct {
	ID string `json:"id"`
	ISBN string `json:"isbn"`
	Title string `json:"title"`
	Author string `json:"author"`
	Price float64 `json:"price"`
}