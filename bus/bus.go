package bus

import (
	"sync"

	"github.com/michilu/bazel-bin-go/log"
	"github.com/vardius/message-bus"
)

var (
	bus messagebus.MessageBus
	wg  sync.WaitGroup

	Publish func(topic string, args ...interface{})
)

func init() {
	bus = messagebus.New()
	Publish = bus.Publish
}

func Subscribe(topic string, fn interface{}) error {
	const op = "bus.Subscribe"

	log.Debug().
		Str("op", op).
		Str("topic", topic).
		Msg("start")

	bus.Subscribe(topic, fn)
	wg.Add(1)

	log.Debug().
		Str("op", op).
		Str("topic", topic).
		Msg("end")

	return nil
}

func Unsubscribe(topic string, fn interface{}) error {
	const op = "bus.Unsubscribe"
	defer wg.Done()

	log.Debug().
		Str("op", op).
		Str("topic", topic).
		Msg("start")

	bus.Unsubscribe(topic, fn)

	log.Debug().
		Str("op", op).
		Str("topic", topic).
		Msg("end")

	return nil
}

func Wait() {
	wg.Wait()
}
