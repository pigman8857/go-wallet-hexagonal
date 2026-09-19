package model

import "time"

type Wallet struct {
	ID       string
	UserID   string
	Balance  float64
	CreateAt time.Time
	UpdateAt time.Time
}
