package eventbus

import (
	"sync"
)

type Data interface{}
type DataChannel chan Data

type EventBus struct {
	mutex  sync.RWMutex
	events map[string][]DataChannel
}

func (e EventBus) Subscribe(eventName string, ch DataChannel) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if subscribers, ok := e.events[eventName]; ok {
		e.events[eventName] = append(subscribers, ch)
		return
	}

	e.events[eventName] = []DataChannel{ch}
}

func (e EventBus) UnSubscribe(eventName string, ch DataChannel) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	subscribers, ok := e.events[eventName]
	if !ok {
		return
	}

	for idx, subscriber := range subscribers {
		if subscriber == ch {
			// drain the channel, not closing it based on The Channel Closing Principle
			go func() {
				defer close(ch)
				for range ch {
				}
			}()

			// the order of subscribers does not matter
			swapSubscribers(idx, len(subscribers)-1, subscribers)
			e.events[eventName] = subscribers[:len(subscribers)-1]

			return
		}
	}
}

func (e EventBus) Publish(eventName string, data Data) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	subscribers, ok := e.events[eventName]
	if !ok {
		return
	}

	go func(data Data, subscribers []DataChannel) {
		for _, ch := range subscribers {
			ch <- data
		}
	}(data, subscribers)
}

func swapSubscribers(i, j int, subscribers []DataChannel) {
	if len(subscribers)-1 < i || len(subscribers)-1 < j {
		return
	}

	subscribers[i], subscribers[j] = subscribers[j], subscribers[i]
}

func NewEventBus() EventBus {
	return EventBus{
		events: make(map[string][]DataChannel),
	}
}
