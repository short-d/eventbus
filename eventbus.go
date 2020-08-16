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

func (e *EventBus) Subscribe(eventName string, ch DataChannel) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if subscribers, ok := e.events[eventName]; ok {
		e.events[eventName] = append(subscribers, ch)
		return
	}

	e.events[eventName] = []DataChannel{ch}
}

func (e *EventBus) UnSubscribe(eventName string, ch DataChannel) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	subscribers, ok := e.events[eventName]
	if !ok {
		return
	}

	for idx, subscriber := range subscribers {
		if subscriber == ch {
			// drain the channel
			defer func() {
				close(ch)
				for range ch {
				}
			}()

			// the order of subscribers does not matter
			swap(idx, len(subscribers)-1, subscribers)
			e.events[eventName] = subscribers[:len(subscribers)-1]

			return
		}
	}
}

func (e *EventBus) Publish(eventName string, data Data) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	subscribers, ok := e.events[eventName]
	if !ok {
		return
	}
	for _, ch := range subscribers {
		ch <- data
	}
}

func swap(i, j int, subscribers []DataChannel) {
	if i >= len(subscribers) || j >= len(subscribers) {
		return
	}

	subscribers[i], subscribers[j] = subscribers[j], subscribers[i]
}

func NewEventBus() EventBus {
	return EventBus{
		events: make(map[string][]DataChannel),
	}
}
