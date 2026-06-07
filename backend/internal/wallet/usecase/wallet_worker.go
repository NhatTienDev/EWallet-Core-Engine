package usecase

import (
	"context"
	"log"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func WalletWorker(userCreatedChan <-chan int64, walletUseCase domain.WalletUseCase) {
	log.Println("Background Worker: Wallet Creator is running...")

	for userID := range userCreatedChan {
		log.Println("[Worker] Receive the wallet creation event for User ID: %d\n", userID)

		// Auto create default VND wallet
		// Use context.Background() because this is background process
		_, err := walletUseCase.CreateUserWallet(context.Background(), userID, "VND")
		if err != nil {
			log.Printf("[Error Worker] Cannot create VND wallet for User ID %d\n", userID, err)
			continue
		}

		log.Printf("[Success Worker] Successfully created VND wallet for User ID %d\n", userID)
	}
}