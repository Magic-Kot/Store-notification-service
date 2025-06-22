package subjects

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const subject = "notification.topic"

func SubscribeNotification(ns *nats.Conn) error {
	_, err := ns.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Printf("[%s] Получено: %s\n", time.Now().Format(time.RFC3339), string(msg.Data))
	})
	if err != nil {
		log.Fatalf("nats.Subscribe: %v", err)
	}

	return nil
}
