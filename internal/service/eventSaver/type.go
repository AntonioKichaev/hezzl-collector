package eventSaver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antoniokichaev/hezzl-collector/internal/entity/event"
	"github.com/antoniokichaev/hezzl-collector/internal/repo"
	"github.com/nats-io/nats.go"
	"time"
)

type EventSaver struct {
	repo repo.Event
	js   nats.JetStreamContext
}

func (es *EventSaver) CheckBatch(ctx context.Context) {

	for {
		sub, _ := es.js.PullSubscribe("events.item",
			"worker",
			nats.PullMaxWaiting(128),
			nats.BindStream("EVENTS"),
		)
		if _, ok := ctx.Deadline(); ok {
			break
		}

		time.Sleep(time.Second * 10)
		msgs, err := sub.FetchBatch(100, nats.Context(ctx))
		if err != nil {
			continue
		}

		batch := make([]event.ClickhouseEvent, 0, 100)
		for msg := range msgs.Messages() {

			ev := event.ClickhouseEvent{}
			_ = json.Unmarshal(msg.Data, &ev)
			_ = msg.Ack()
			batch = append(batch, ev)
		}
		if len(batch) != 0 {
			err = es.repo.CreateEvent(ctx, batch)
		}
		fmt.Println(batch, err)
	}
}

func (es *EventSaver) Start(ctx context.Context) {
	go es.CheckBatch(ctx)
}

func New(evRepo repo.Event, js nats.JetStreamContext) *EventSaver {
	es := &EventSaver{
		js:   js,
		repo: evRepo,
	}
	return es
}
