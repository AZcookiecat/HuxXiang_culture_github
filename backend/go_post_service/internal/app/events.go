package app

import "context"

type Event struct {
	Topic string
	Key   string
}

type EventBus struct {
	ch    chan Event
	cache Cache
}

func NewEventBus(cache Cache) *EventBus {
	return &EventBus{
		ch:    make(chan Event, 256),
		cache: cache,
	}
}

func (b *EventBus) Publish(event Event) {
	select {
	case b.ch <- event:
	default:
	}
}

func (b *EventBus) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-b.ch:
			switch event.Topic {
			case "community.cache.invalidate":
				b.cache.DeletePrefix("community:")
			}
		}
	}
}
