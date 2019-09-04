package main

import (
	"fmt"
	"time"

	"github.com/byliuyang/eventbus"
)

func main() {
	bus := eventbus.NewEventBus()

	notificationChannel := make(eventbus.DataChannel)
	notification := "notification"
	bus.Subscribe(notification, notificationChannel)

	go func() {
		for {
			select {
			case data := <-notificationChannel:
				fmt.Println(data)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	bus.Publish(notification, "Hello!")
	bus.UnSubscribe(notification, notificationChannel)
	time.Sleep(1 * time.Second)
}
