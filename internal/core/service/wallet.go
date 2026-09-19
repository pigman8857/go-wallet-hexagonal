package service

import (
	"context"
	"fmt"
	"hexagonal_intro/internal/core/model"
	"hexagonal_intro/internal/core/port"
	"sync"
)

type WalletService struct {
	repo port.WalletRepository
	mu   sync.RWMutex
}

//Contructor

func NewWalletService(repo port.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

// Member methods of WalletService
func (s *WalletService) getWallet(ctx context.Context, id string) (*model.Wallet, error) {
	wallet, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get wallet: %w\n", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("wallet %s not found\n", id)
	}

	return wallet, nil

}

func (ws *WalletService) CreateWallet(ctx context.Context, userID string) (*model.Wallet, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	existing, err := ws.repo.GetByUserID(ctx, userID)

	if err != nil {
		return nil, fmt.Errorf("checking existing wallet: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user %s already has a wallet\n", userID)
	}

	wallet := &model.Wallet{
		UserID:  userID,
		Balance: 0,
	}

	if err := ws.repo.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("Create wallet %w", err)
	}

	return wallet, nil

}

func (ws *WalletService) GetWallet(ctx context.Context, id string) (*model.Wallet, error) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	return ws.getWallet(ctx, id)
}

func (ws *WalletService) Deposit(ctx context.Context, id string, amount float64) (*model.Wallet, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()

	wallet, err := ws.getWallet(ctx, id)
	if err != nil {
		return nil, err
	}
	newBalance := wallet.Balance + amount
	if err := ws.repo.UpdateBalance(ctx, id, newBalance); err != nil {

		return nil, fmt.Errorf("Update balance: %w", err)

	}
	wallet.Balance = newBalance
	return wallet, nil
}

func (ws *WalletService) WithDraw(ctx context.Context, id string, amount float64) (*model.Wallet, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	wallet, err := ws.getWallet(ctx, id)
	if err != nil {
		return nil, err
	}
	newBalance := wallet.Balance - amount
	if wallet.Balance < amount {
		return nil, fmt.Errorf("insufficient balance: have %.2f, need %.2f", wallet.Balance, amount)
	}
	wallet.Balance = newBalance
	return wallet, nil
}

func (ws *WalletService) DeleteWallet(ctx context.Context, id string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if err := ws.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete wallet: %w", err)
	}

	return nil
}
