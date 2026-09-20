package main

import (
	"context"
	"hexagonal_intro/internal/adapter/handler"
	"hexagonal_intro/internal/adapter/handler/middleware"
	"hexagonal_intro/internal/adapter/repository"
	"hexagonal_intro/internal/config"
	"hexagonal_intro/internal/core/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Connection pool tuning goes here, right after the pool is created and
	// before it's used — not at shutdown, where it'd have no effect.
	//
	// How to size MaxOpenConns per pod:
	//   max_open_conns_per_pod = (postgres_max_connections - reserved) / max_pod_replicas
	//   - postgres_max_connections: `SHOW max_connections;` on the DB (Postgres default 100)
	//   - reserved: headroom for migrations, psql sessions, monitoring (~10-20% of total)
	//   - max_pod_replicas: your autoscaler's MAX replica count, not the current count —
	//     otherwise scaling up can push total connections past Postgres's limit and
	//     start failing with "too many connections"
	//   Example: 100 max_connections, reserve 10, max 5 replicas -> (100-10)/5 = 18
	//
	// sqlDB, _ := db.DB()
	//
	// SetMaxOpenConns: hard cap on connections this pod can hold open to Postgres at once.
	// Requests beyond this queue and wait instead of opening more connections.
	// sqlDB.SetMaxOpenConns(25)
	//
	// SetMaxIdleConns: connections kept open (not closed) between requests when idle.
	// Too low relative to MaxOpenConns causes constant open/close churn under steady traffic.
	// sqlDB.SetMaxIdleConns(10)
	//
	// SetConnMaxLifetime: force a connection to close and be replaced after this long,
	// even if healthy. Protects against a DB failover/restart leaving stale connections
	// in the pool that look fine but no longer work.
	// sqlDB.SetConnMaxLifetime(30 * time.Minute)
	//
	// SetConnMaxIdleTime: close a connection if it's been idle this long, freeing it
	// back to Postgres instead of holding it open unused (helps when traffic is bursty).
	// sqlDB.SetConnMaxIdleTime(30 * time.Minute)

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

	//this is graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v\n", err)
		}
	}()

	// This is not graceful shutdown
	// if err := r.Run(":" + port); err != nil {
	// 	log.Fatalf("Server failed: %v\n", err)
	// }

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("received signal %v - shutting down gracefully...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Force shuwdown: %v\n", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	log.Println("server stopped cleanly")
}
