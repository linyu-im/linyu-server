package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"sync"

	mqclient "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/event"
)

type RocketMQEventBus struct {
	producer  mqclient.Producer
	namesrv   []string
	consumers []mqclient.PushConsumer
	mu        sync.Mutex
}

func NewRocketMQEventBus() (*RocketMQEventBus, error) {
	p, err := mqclient.NewProducer(
		producer.WithNameServer(config.C.EventBus.RocketEventBus.Namesrv),
		producer.WithRetry(config.C.EventBus.RocketEventBus.RetryTimes),
	)
	if err != nil {
		return nil, err
	}
	if err := p.Start(); err != nil {
		return nil, err
	}
	return &RocketMQEventBus{
		producer: p,
		namesrv:  config.C.EventBus.RocketEventBus.Namesrv,
	}, nil
}

func (b *RocketMQEventBus) Subscribe(eventName string, factory Factory, handler event.Handler) {
	c, err := mqclient.NewPushConsumer(
		consumer.WithGroupName("g_"+eventName),
		consumer.WithNameServer(b.namesrv),
		consumer.WithConsumerModel(consumer.BroadCasting),
	)
	if err != nil {
		fmt.Printf("create consumer error: %v\n", err)
		return
	}

	err = c.Subscribe(eventName, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			newEv := factory()
			if err := json.Unmarshal(msg.Body, newEv); err != nil {
				continue
			}
			go func(ev event.Event) {
				_ = handler(ev)
			}(newEv)
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		fmt.Printf("subscribe error: %v\n", err)
		return
	}
	if err := c.Start(); err != nil {
		fmt.Printf("start consumer error: %v\n", err)
		return
	}

	b.mu.Lock()
	b.consumers = append(b.consumers, c)
	b.mu.Unlock()
}

func (b *RocketMQEventBus) Publish(e event.Event) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: e.EventName(),
		Body:  data,
	}
	_, err = b.producer.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("RocketMQ Publish Error: %v\n", err)
		return err
	}
	return nil
}

func (b *RocketMQEventBus) Close() {
	if b.producer != nil {
		_ = b.producer.Shutdown()
	}
	b.mu.Lock()
	for _, c := range b.consumers {
		_ = c.Shutdown()
	}
	b.mu.Unlock()
}
