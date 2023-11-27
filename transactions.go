package main

import "time"

type Transaction struct {
	ID     string    `json:"id"`
	Order  *Order    `json:"order"`
	Date   time.Time `json:"date"`
	Amount string    `json:"amount"`
}
