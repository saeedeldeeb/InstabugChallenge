package providers

import (
	workers "chat/pkg/rabbitmq"
	"github.com/goravel/framework/contracts/foundation"
	"log"
)

type AppServiceProvider struct {
}

func (receiver *AppServiceProvider) Register(app foundation.Application) {

}

func (receiver *AppServiceProvider) Boot(app foundation.Application) {
	// Create a new worker in go routine
	go func() {
		worker, err := workers.NewMessageWorker("amqp://guest:guest@localhost:5672/", "message_queue")
		if err != nil {
			log.Fatalf("Failed to create worker: %v", err)
		}
		defer worker.Close()

		// Start the worker
		if err := worker.Start(); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()
}
