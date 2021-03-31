package eventbus

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
)

type (
	// EventBus implements the publish–subscribe messaging pattern
	EventBus interface {
		// Publish publishes arguments to the given topic subscribers
		Publish(topic string, args ...interface{})
		// Subscribe subscribes to the given topic
		Subscribe(topic string, fn interface{}, priority ...Priority) error
		// SubscribeAsync subscribes to the given topic asynchronously
		SubscribeAsync(topic string, fn interface{}) error
		// Unsubscribe unsubscribes handler from the given topic
		Unsubscribe(topic string, fn interface{}) error
		// Close unsubscribes all handlers from given topic
		Close(topic string) error
	}
	// Priority represents an event priority
	Priority int
)

const (
	Lowest  Priority = -200
	Low     Priority = -100
	Normal  Priority = 0
	High    Priority = 100
	Highest Priority = 200
)

type (
	eventBus struct {
		mutex    sync.RWMutex
		handlers map[string][]handler
	}
	handler struct {
		priority Priority
		async    bool
		fn       reflect.Value
	}
)

func (b *eventBus) Publish(topic string, args ...interface{}) {
	var reflectedArgs []reflect.Value
	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	handlers := b.copyHandlers(topic)
	for _, handler := range handlers {
		if handler.async {
			go handler.fn.Call(reflectedArgs)
		} else {
			handler.fn.Call(reflectedArgs)
		}
	}
}

func (b *eventBus) Subscribe(topic string, fn interface{}, priority ...Priority) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	prio := Normal
	if len(priority) > 0 {
		prio = priority[0]
	}

	b.subscribeHandler(topic, handler{
		priority: prio,
		async:    false,
		fn:       reflect.ValueOf(fn),
	})
	return nil
}

func (b *eventBus) SubscribeAsync(topic string, fn interface{}) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	b.subscribeHandler(topic, handler{
		priority: Normal,
		async:    true,
		fn:       reflect.ValueOf(fn),
	})
	return nil
}

func (b *eventBus) subscribeHandler(topic string, handler handler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	handlers := append(b.handlers[topic], handler)
	sort.SliceStable(handlers, func(i, j int) bool {
		return handlers[i].priority < handlers[j].priority
	})
	b.handlers[topic] = handlers
}

func (b *eventBus) Unsubscribe(topic string, fn interface{}) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	rv := reflect.ValueOf(fn)

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.handlers[topic]; ok {
		for i, handler := range b.handlers[topic] {
			if handler.fn == rv {
				if len(b.handlers[topic]) == 1 {
					delete(b.handlers, topic)
				} else {
					b.handlers[topic] = append(b.handlers[topic][:i], b.handlers[topic][i+1:]...)
				}
			}
		}

		return nil
	}

	return fmt.Errorf("topic %s doesn't exist", topic)
}

func (b *eventBus) Close(topic string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.handlers[topic]; ok {
		delete(b.handlers, topic)
		return nil
	}

	return fmt.Errorf("topic %s doesn't exist", topic)
}

func (b *eventBus) copyHandlers(topic string) []handler {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if handlers, ok := b.handlers[topic]; ok {
		return handlers
	}
	return []handler{}
}

func isValidHandler(fn interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}
	return nil
}

// New creates a new EventBus
func New() EventBus {
	return &eventBus{
		handlers: make(map[string][]handler),
	}
}
