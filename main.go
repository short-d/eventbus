package main

import (
	"fmt"
	"time"
)

func main() {
	bus := NewEventBus()

	notificationChannel := make(DataChannel)
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
