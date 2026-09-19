package repository

import (
	"context"
	"hexagonal_intro/internal/core/model"
	"time"

	"gorm.io/gorm"
)

type WalletGORM struct {
	ID       string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID   string    `gorm:"column:user_id;not null;uniqueIndex"`
	Balance  float64   `gorm:"column:balance;type:decimal(18,2);not null;default:0"`
	CreateAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdateAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}

func (WalletGORM) TableName() string {
	return "wallets"
}

// From gorm to domain model, so we can use within applcation
func (w WalletGORM) toDomain() model.Wallet {
	return model.Wallet{
		ID:       w.ID,
		UserID:   w.UserID,
		Balance:  w.Balance,
		CreateAt: w.CreateAt,
		UpdateAt: w.UpdateAt,
	}
}

// From domain to gorm model
func formDomain(w model.Wallet) WalletGORM {
	return WalletGORM{
		ID:       w.ID,
		UserID:   w.UserID,
		Balance:  w.Balance,
		CreateAt: w.CreateAt,
		UpdateAt: w.UpdateAt,
	}
}

type PostgresWalletRespository struct {
	db *gorm.DB
}

// Contructor
func NewPostgresWalletRespository(db *gorm.DB) *PostgresWalletRespository {
	return &PostgresWalletRespository{db: db}
}

func (repo *PostgresWalletRespository) Create(ctx context.Context, wallet *model.Wallet) error {
	gormWallet := formDomain(*wallet)
	result := repo.db.WithContext(ctx).Create(&gormWallet)
	if result.Error != nil {
		return result.Error
	}

	*wallet = gormWallet.toDomain()
	return nil
}

func (repo *PostgresWalletRespository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	var gormWallet WalletGORM
	result := repo.db.WithContext(ctx).First(&gormWallet, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	domainWallet := gormWallet.toDomain()
	return &domainWallet, nil
}

func (repo *PostgresWalletRespository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	var gormWallet WalletGORM
	result := repo.db.WithContext(ctx).First(&gormWallet, "user_id = ?", userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	domainWallet := gormWallet.toDomain()
	return &domainWallet, nil
}

func (repo *PostgresWalletRespository) UpdateBalance(ctx context.Context, id string, balance float64) error {
	result := repo.db.WithContext(ctx).Model(&WalletGORM{}).Where("id = ?", id).Update("balance", balance)
	return result.Error
}

func (repo *PostgresWalletRespository) Delete(ctx context.Context, id string) error {
	result := repo.db.WithContext(ctx).Delete(&WalletGORM{}, "id = ?", id)
	return result.Error
}
