package service_test

import (
	"context"
	"hexagonal_intro/internal/core/model"
	"hexagonal_intro/internal/core/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	wallets map[string]*model.Wallet
}

func (mockRepo *mockRepo) Create(_ context.Context, wallet *model.Wallet) error {
	wallet.ID = "mock-id-001"
	mockRepo.wallets[wallet.ID] = wallet
	return nil
}
func (mockRepo *mockRepo) GetByID(_ context.Context, id string) (*model.Wallet, error) {
	w, ok := mockRepo.wallets[id]

	if !ok {
		return nil, nil
	}

	return w, nil
}
func (mockRepo *mockRepo) GetByUserID(_ context.Context, userID string) (*model.Wallet, error) {
	for _, w := range mockRepo.wallets {
		if w.UserID == userID {
			return w, nil
		}
	}

	return nil, nil
}
func (mockRepo *mockRepo) UpdateBalance(_ context.Context, id string, balance float64) error {
	if w, ok := mockRepo.wallets[id]; ok {
		w.Balance = balance
	}

	return nil
}
func (mock *mockRepo) Delete(_ context.Context, id string) error {
	delete(mock.wallets, id)

	return nil
}

func newMockRepo() *mockRepo {
	return &mockRepo{wallets: make(map[string]*model.Wallet)}
}

func TestCreateWallet_Success(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, err := svc.CreateWallet(context.Background(), "user-1")
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, "user-1", wallet.UserID)
	assert.Equal(t, 0.0, wallet.Balance)

}

func TestCreateWallet_DuplicateUser(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	_, err := svc.CreateWallet(context.Background(), "user-1")
	assert.NoError(t, err)

	_, err2 := svc.CreateWallet(context.Background(), "user-1")
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "already has a wallet")
}

func TestDeposit_Success(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, _ := svc.CreateWallet(context.Background(), "user-1")
	updated, err := svc.Deposit(context.Background(), wallet.ID, 100.0)

	assert.NoError(t, err)
	assert.Equal(t, 100.0, updated.Balance)
}

func TestDeposit_NegativeAmount(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, _ := svc.CreateWallet(context.Background(), "user-1")
	_, err := svc.Deposit(context.Background(), wallet.ID, -50.0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")
}

func TestWithDraw_Sucess(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, _ := svc.CreateWallet(context.Background(), "user-1")
	svc.Deposit(context.Background(), wallet.ID, 200.0)
	updated, err := svc.WithDraw(context.Background(), wallet.ID, 50.0)
	assert.NoError(t, err)
	assert.Equal(t, 150.0, updated.Balance)
}

func TestWithDraw_InsufficentBalance(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	wallet, _ := svc.CreateWallet(context.Background(), "user-1")
	_, err := svc.WithDraw(context.Background(), wallet.ID, 10.0)
	assert.Error(t, err)
	assert.Equal(t, "insufficient balance: have 0.00, need 10.00", err.Error())
}

func TestGetWallet_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewWalletService(repo)

	_, err := svc.GetWallet(context.Background(), "non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
