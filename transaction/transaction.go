package transaction

import (
	"time"
)

type Transaction struct {
	ID   string `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	Description string `json:"description"`
	Amount float64 `json:"amount"`
	Type string `json:"type"`
	From string `json:"from"`
	To string `json:"to"`
}
