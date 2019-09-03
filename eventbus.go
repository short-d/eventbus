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
	if subscribers, ok := e.events[eventName]; ok {
		e.events[eventName] = append(subscribers, ch)
		e.mutex.Unlock()
		return
	}

	e.events[eventName] = []DataChannel{ch}
	e.mutex.Unlock()
}

func (e EventBus) UnSubscribe(eventName string, ch DataChannel) {
	e.mutex.Lock()
	subscribers, ok := e.events[eventName]
	if !ok {
		e.mutex.Unlock()
		return
	}

	for idx, subscriber := range subscribers {
		if subscriber == ch {
			e.events[eventName] = append(subscribers[:idx], subscribers[idx+1:]...)
			e.mutex.Unlock()
			return
		}
	}

	e.mutex.Unlock()
}

func (e EventBus) Publish(eventName string, data Data) {
	e.mutex.RLock()
	subscribers, ok := e.events[eventName]
	if !ok {
		e.mutex.RUnlock()
		return
	}

	for _, subscriber := range subscribers {
		subscriber <- data
	}
	e.mutex.RUnlock()
}

func NewEventBus() EventBus {
	return EventBus{
		events: make(map[string][]DataChannel),
	}
}
