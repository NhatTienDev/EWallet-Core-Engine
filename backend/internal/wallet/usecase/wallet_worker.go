package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func WalletWorker(ctx context.Context, userCreatedChan <-chan int64, walletUseCase domain.WalletUseCase) {
	log.Println("Background Worker: Wallet Creator is running...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Background Worker: Shutdown signal recieved")
			return
		case userID, ok := <-userCreatedChan:
			if !ok {
				log.Println("Background Worker: Channel closed")
				return
			}

			log.Printf("[Worker] Receive the wallet creation event for User ID: %d", userID)

			walletWorkerCtx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)

			_, err := walletUseCase.CreateUserWallet(walletWorkerCtx, userID, "VND")

			cancel()

			if err != nil {
				log.Printf("[Worker] Cannot create wallet for User ID %d: %v", userID, err)
				continue
			}

			log.Printf("[Worker] Successfully created VND wallet for User ID %d", userID)
		}
	}
}