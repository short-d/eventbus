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
			// the order of subscribers does not matter
			subscribers[idx], subscribers[len(subscribers)-1] = subscribers[len(subscribers)-1], nil
			e.events[eventName] = subscribers[:len(subscribers)-1]

			// drain the channel
			for {
				select {
				case <-ch:
				default:
					return
				}
			}
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

	for _, subscriber := range subscribers {
		subscriber <- data
	}
}

func NewEventBus() EventBus {
	return EventBus{
		events: make(map[string][]DataChannel),
	}
}
