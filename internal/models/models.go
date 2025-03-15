package models

import "time"

type User struct {
	ID         int
	Name       string
	Balance    float64
	CreatedAt time.Time
}

type TransactionHistory struct {
	ID int
	UserId int
	Amount float64
	TransactionType string
	CreatedAt time.Time
}