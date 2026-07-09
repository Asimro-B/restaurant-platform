package worker

import (
	"log"
	orderworkflow "restaurant-platform/internal/workflows/order"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func StartWorker(temporalClient client.Client, activities *orderworkflow.OrderActivities) {
	w := worker.New(temporalClient, orderworkflow.TaskQueue, worker.Options{})

	w.RegisterWorkflow(orderworkflow.CreateOrderWorkflow)
	w.RegisterActivity(activities)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("Worker failed: %v", err)
	}
}
