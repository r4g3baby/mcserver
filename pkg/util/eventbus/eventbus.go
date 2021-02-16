package eventbus

import (
	"fmt"
	"reflect"
	"sync"
)

// EventBus implements the publish–subscribe messaging pattern
type EventBus interface {
	// Publish publishes arguments to the given topic subscribers
	Publish(topic string, args ...interface{})
	// Subscribe subscribes to the given topic
	Subscribe(topic string, fn interface{}) error
	// SubscribeAsync subscribes to the given topic asynchronously
	SubscribeAsync(topic string, fn interface{}) error
	// Unsubscribe unsubscribe handler from the given topic
	Unsubscribe(topic string, fn interface{}) error
	// Close unsubscribe all handlers from given topic
	Close(topic string) error
}

type (
	eventBus struct {
		mutex    sync.RWMutex
		handlers map[string][]handler
	}
	handler struct {
		async bool
		fn    reflect.Value
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

func (b *eventBus) Subscribe(topic string, fn interface{}) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	handler := handler{
		async: false,
		fn:    reflect.ValueOf(fn),
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.handlers[topic] = append(b.handlers[topic], handler)
	return nil
}

func (b *eventBus) SubscribeAsync(topic string, fn interface{}) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	handler := handler{
		async: true,
		fn:    reflect.ValueOf(fn),
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.handlers[topic] = append(b.handlers[topic], handler)
	return nil
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
