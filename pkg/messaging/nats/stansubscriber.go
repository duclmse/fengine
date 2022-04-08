package nats

import (
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"

	log "github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/pkg/messaging"
	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

var _ messaging.Subscriber = (*stanSubs)(nil)

type stanSubs struct {
	conn          stan.Conn
	logger        log.Logger
	mu            sync.Mutex
	queue         string
	subscriptions map[string]stan.Subscription
}

type Subscriber interface {
	messaging.Subscriber
	Close()
}

// NewStanSub returns NATS streaming server (STAN) subscriber.
// Parameter queue specifies the queue for the Subscribe method.
// If queue is specified (is not an empty string), Subscribe method
// will execute NATS QueueSubscribe which is conceptually different
// from ordinary subscribe. For more information, please take a look
// here: https://docs.nats.io/developing-with-nats/receiving/queues.
// If the queue is empty, Subscribe will be used.
func NewStanSub(url, queue, clusterid, clientid string, logger log.Logger) (Subscriber, error) {
	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		fmt.Printf("Can not connect to nats: %s. Error %v \n", url, err)
		return nil, err
	}
	conn, err := stan.Connect(clusterid, clientid, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			logger.Error(fmt.Sprintf("Connection lost, reason: %v", reason))
		}))
	if err != nil {
		return nil, err
	}
	ret := &stanSubs{
		conn:          conn,
		logger:        logger,
		queue:         queue,
		subscriptions: make(map[string]stan.Subscription),
	}
	return ret, nil
}

func (ps *stanSubs) Subscribe(topic string, handler messaging.MessageHandler) error {
	if topic == "" {
		return errEmptyTopic
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, ok := ps.subscriptions[topic]; ok {
		return errAlreadySubscribed
	}
	ps.logger.Info(fmt.Sprintf("Subcribe for topic: %s", topic))
	nh := ps.stanHandler(handler)

	if ps.queue != "" {
		sub, err := ps.conn.QueueSubscribe(topic, ps.queue, nh)
		if err != nil {
			return err
		}
		ps.subscriptions[topic] = sub
		return nil
	} else {
		sub, err := ps.conn.Subscribe(topic, nh)
		if err != nil {
			return err
		}
		ps.subscriptions[topic] = sub
	}
	return nil
}

func (ps *stanSubs) Unsubscribe(topic string) error {
	if topic == "" {
		return errEmptyTopic
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()

	sub, ok := ps.subscriptions[topic]
	if !ok {
		return errNotSubscribed
	}

	if err := sub.Unsubscribe(); err != nil {
		return err
	}

	delete(ps.subscriptions, topic)
	return nil
}

func (ps *stanSubs) Close() {
	ps.conn.Close()
}

func (ps *stanSubs) stanHandler(h messaging.MessageHandler) stan.MsgHandler {
	return func(m *stan.Msg) {
		ps.logger.Info(fmt.Sprintf("Received message: "))
		var msg messaging.Message
		if err := proto.Unmarshal(m.Data, &msg); err != nil {
			ps.logger.Warn(fmt.Sprintf("Failed to unmarshal received message: %s", err))
			return
		}
		if err := h(msg); err != nil {
			ps.logger.Warn(fmt.Sprintf("Failed to handle message: %s", err))
		}
	}
}
