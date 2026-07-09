package pusherclient

import (
	"fmt"
	"os"

	"github.com/pusher/pusher-http-go/v5"
)

var client *pusher.Client

func Init() {
	client = &pusher.Client{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: os.Getenv("PUSHER_CLUSTER"),
		Secure:  true,
	}
}

func Publish(channel, event string, data interface{}) error {
	if client == nil {
		return fmt.Errorf("pusher client not initialized")
	}
	err := client.Trigger(channel, event, data)
	return err
}

func KitchenChannel(tenantID int64) string {
	return fmt.Sprintf("kitchen-%d", tenantID)
}

func FloorChannel(tenantID int64) string {
	return fmt.Sprintf("floor-%d", tenantID)
}

func ManageChannel(tenantID int64) string {
	return fmt.Sprintf("manager-%d", tenantID)
}
