package providers

import (
	workers "chat/pkg/rabbitmq"
	"github.com/goravel/framework/contracts/foundation"
	"log"
	"time"
)

type AppServiceProvider struct {
}

func (receiver *AppServiceProvider) Register(app foundation.Application) {

}

func (receiver *AppServiceProvider) Boot(app foundation.Application) {
	// Create a new worker in go routine
	go func() {
		// Wait for 10 seconds before connecting to the worker
		time.Sleep(10 * time.Second)

		var worker *workers.MessageWorker
		var err error
		for i := 0; i < 3; i++ {
			worker, err = workers.NewMessageWorker("amqp://guest:guest@instabug-rabbitmq:5672/", "message_queue")
			if err == nil {
				break
			}
			log.Printf("Failed to create worker (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second) // Wait before retrying
		}
		if err != nil {
			log.Fatalf("Failed to create worker after 3 attempts: %v", err)
		}
		defer worker.Close()

		// Start the worker
		if err := worker.Start(); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()
}
