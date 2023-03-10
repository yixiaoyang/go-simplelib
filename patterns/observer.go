package patterns

import "container/list"

type (
	Event struct {
		Msg string
	}

	Notifier interface {
		Add(Observer)
		Remove(Observer)
		Notify(Event)
	}

	Observer interface {
		OnNotify(Event)
	}
)

type (
	ChatNotifier struct {
		observers map[Observer]struct{}
	}

	ChatObserver struct {
		Name      string
		EventList *list.List
	}

	ChatEvent struct {
		Content string
	}
)

func (notifier *ChatNotifier) Add(observer Observer) {
	notifier.observers[observer] = struct{}{}
}

func (notifier *ChatNotifier) Remove(observer Observer) {
	delete(notifier.observers, observer)
}

func (notifier *ChatNotifier) Notify(e Event) {
	for observer := range notifier.observers {
		observer.OnNotify(e)
	}
}

func (observer *ChatObserver) OnNotify(e Event) {
	observer.EventList.PushBack(e)
}
