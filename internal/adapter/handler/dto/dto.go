package dto

type CreateWalletRequestDto struct {
	UserID string `json:"user_id" binding:"required"`
}

type WalletResponseDto struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type TransactionRequestDto struct {
	Amount float64 `json:"amount" binding:"required"`
}

type ErrorResponseDto struct {
	Error string `json:"error"`
}
