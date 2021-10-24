package models

type Feedback struct {
	Rating       float64 `json:"rating"`
	Memo         string  `json:"memo"`
	CustomerID   int64   `json:"custID"`
	CustomerName string  `json:"name"`
}
