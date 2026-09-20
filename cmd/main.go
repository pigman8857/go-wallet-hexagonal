package main

import (
	"hexagonal_intro/internal/adapter/handler"
	"hexagonal_intro/internal/adapter/handler/middleware"
	"hexagonal_intro/internal/adapter/repository"
	"hexagonal_intro/internal/config"
	"hexagonal_intro/internal/core/service"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	config.LoadEnv(".env")
	dsn := config.GetEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=walletdb port=5432 sslmode=disable TimeZone=Asia/Bangkok")
	port := config.GetEnv("PORT", "8080")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}

	if err := db.AutoMigrate(&repository.WalletGORM{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	walletRepo := repository.NewPostgresWalletRespository(db)
	walletSvc := service.NewWalletService(walletRepo)
	walletHdr := handler.NewWalletHandler(walletSvc)

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	wallets := r.Group("/wallets")
	{
		wallets.POST("", walletHdr.CreateWallet)
		wallets.GET("/:id", walletHdr.GetWallet)
		wallets.POST("/:id/deposit", walletHdr.Deposit)
		wallets.POST("/:id/withDraw", walletHdr.WithDraw)
		wallets.DELETE("/:id", walletHdr.DeleteWallet)
	}

	log.Printf("Server starting on: %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
