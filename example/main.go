package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/short-d/eventbus"
)

func main() {
	bus := eventbus.NewEventBus()

	notificationChannel := make(eventbus.DataChannel)
	notification := "notification"
	bus.Subscribe(notification, notificationChannel)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		for {
			select {
			case data := <-notificationChannel:
				if data != nil {
					fmt.Println(data)
					bus.UnSubscribe(notification, notificationChannel)
					waitGroup.Done()
				}
			default:
				break
			}
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		bus.Publish(notification, "Hello!")
	}()
	waitGroup.Wait()
}
