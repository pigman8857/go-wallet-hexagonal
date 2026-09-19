package port

import (
	"context"
	"hexagonal_intro/internal/core/model"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *model.Wallet) error
	GetByID(ctx context.Context, id string) (*model.Wallet, error)
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	UpdateBalance(ctx context.Context, id string, balance float64) error
	Delete(ctx context.Context, id string) error
}
