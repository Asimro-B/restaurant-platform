package orderworkflow

import (
	"fmt"
	"time"

	"restaurant-platform/internal/models"

	"github.com/shopspring/decimal"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueue            = "order-task-queue"
	SignalKitchenStarted = "kitchen-started"
	SignalKitchenDone    = "kitchen-done"
	SignalOrderServed    = "order-served"
)

func CreateOrderWorkflow(ctx workflow.Context, input models.CreateOrderInput) (*models.CreateOrderResult, error) {
	opts := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	})

	var a *OrderActivities

	// Step 1: Validate items + calculate prices
	var validateResult models.ValidateItemsResult
	err := workflow.ExecuteActivity(opts, a.ValidateAndPriceItems, input).Get(opts, &validateResult)
	if err != nil {
		return nil, err
	}

	// Step 2: Persist order
	var order models.Order
	err = workflow.ExecuteActivity(opts, a.PersistOrder, input, &validateResult).Get(opts, &order)
	if err != nil {
		return nil, err
	}

	// Step 3: Confirm + notify kitchen
	err = workflow.ExecuteActivity(opts, a.ConfirmOrder, order.ID, input.TenantID).Get(opts, nil)
	if err != nil {
		return nil, err
	}

	// Step 4: Wait for kitchen to start preparing
	workflow.GetSignalChannel(ctx, SignalKitchenStarted).Receive(ctx, nil)
	err = workflow.ExecuteActivity(opts, a.MarkPreparing, order.ID, input.TenantID).Get(opts, nil)
	if err != nil {
		return nil, err
	}

	// Step 5: Wait for kitchen done
	workflow.GetSignalChannel(ctx, SignalKitchenDone).Receive(ctx, nil)
	err = workflow.ExecuteActivity(opts, a.MarkReady, order.ID, input.TenantID).Get(opts, nil)
	if err != nil {
		return nil, err
	}

	// Step 6: Wait for served
	workflow.GetSignalChannel(ctx, SignalOrderServed).Receive(ctx, nil)
	err = workflow.ExecuteActivity(opts, a.MarkServed, order.ID, input.TenantID).Get(opts, nil)
	if err != nil {
		return nil, err
	}

	total, err := decimal.NewFromString(validateResult.Total)
	if err != nil {
		return nil, fmt.Errorf("invalid total: %w", err)
	}

	return &models.CreateOrderResult{
		OrderID:     order.ID,
		TotalAmount: total.InexactFloat64(),
		Status:      string(models.OrderStatusServed),
	}, nil
}
