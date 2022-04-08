package nats

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gogo/protobuf/proto"

	log "github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/pkg/messaging"
	broker "github.com/nats-io/nats.go"
)

const (
	chansPrefix        = "channels" //top prefix for all messages
	userDefineTopic    = "userdefined"
	dataTopic          = "data"
	eventTopic         = "events"
	appCtrlTopic       = "app.controls"
	deviceCtrlTopic    = "device.controls"
	SubjectAllChannels = "channels.>"
)

type NatsTopic struct {
	Userdefined *string
	Subtopic    string
	DeviceId    string
	DeviceGroup string
	DeviceType  string
}

// SubjectAllChannels represents subject to subscribe for all the channels.

var (
	errAlreadySubscribed = errors.New("already subscribed to topic")
	errNotSubscribed     = errors.New("not subscribed")
	errEmptyTopic        = errors.New("empty topic")
)

var _ messaging.PubSub = (*pubsub)(nil)

// PubSub wraps messaging Publisher exposing
// Close() method for NATS connection.
type PubSub interface {
	messaging.PubSub
	Close()
}

type pubsub struct {
	conn          *broker.Conn
	logger        log.Logger
	mu            sync.Mutex
	queue         string
	subscriptions map[string]*broker.Subscription
}

// NewPubSub returns NATS message publisher/subscriber.
// Parameter queue specifies the queue for the Subscribe method.
// If queue is specified (is not an empty string), Subscribe method
// will execute NATS QueueSubscribe which is conceptually different
// from ordinary subscribe. For more information, please take a look
// here: https://docs.nats.io/developing-with-nats/receiving/queues.
// If the queue is empty, Subscribe will be used.
func NewPubSub(url, queue string, logger log.Logger) (PubSub, error) {
	conn, err := broker.Connect(url)
	if err != nil {
		return nil, err
	}
	ret := &pubsub{
		conn:          conn,
		queue:         queue,
		logger:        logger,
		subscriptions: make(map[string]*broker.Subscription),
	}
	return ret, nil
}

func (ps *pubsub) Publish(topic string, msg messaging.Message) error {
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("%s.%s", chansPrefix, topic)
	if msg.Subtopic != "" {
		subject = fmt.Sprintf("%s.%s", subject, msg.Subtopic)
	}
	if err := ps.conn.Publish(subject, data); err != nil {
		return err
	}

	return nil
}

func (ps *pubsub) Subscribe(topic string, handler messaging.MessageHandler) error {
	if topic == "" {
		return errEmptyTopic
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, ok := ps.subscriptions[topic]; ok {
		return errAlreadySubscribed
	}
	nh := ps.natsHandler(handler)

	if ps.queue != "" {
		sub, err := ps.conn.QueueSubscribe(topic, ps.queue, nh)
		if err != nil {
			return err
		}
		ps.subscriptions[topic] = sub
		return nil
	}
	sub, err := ps.conn.Subscribe(topic, nh)
	if err != nil {
		return err
	}
	ps.subscriptions[topic] = sub
	return nil
}

func (ps *pubsub) Unsubscribe(topic string) error {
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

func (ps *pubsub) Close() {
	ps.conn.Close()
}

func (ps *pubsub) natsHandler(h messaging.MessageHandler) broker.MsgHandler {
	return func(m *broker.Msg) {
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

func ParseTopic(topic string) NatsTopic {
	var sTopic NatsTopic
	arrTopics := strings.Split(topic, ".")
	if len(arrTopics) < 4 { //invalid topic, should be in format channels....subtopic.deviceid.devicegroup.devicetype
		return sTopic
	}
	if arrTopics[0] == userDefineTopic { //topic in form: channel.userdefined.subtopic.deviceid.devicegroup.devicetype
		if len(arrTopics) == 5 { //subtopic should be events or data
			sTopic.Userdefined = &arrTopics[0]
			sTopic.Subtopic = arrTopics[1]
			sTopic.DeviceId = arrTopics[2]
			sTopic.DeviceGroup = arrTopics[3]
			sTopic.DeviceType = arrTopics[4]
		}
		if len(arrTopics) == 6 { //subtopic should be app.controls or device.controls
			sTopic.Userdefined = &arrTopics[0]
			sTopic.Subtopic = arrTopics[1] + "." + arrTopics[2]
			sTopic.DeviceId = arrTopics[3]
			sTopic.DeviceGroup = arrTopics[4]
			sTopic.DeviceType = arrTopics[5]
		}
	} else { //topic in form: channel.subtopic.deviceid.devicegroup.devicetype
		if len(arrTopics) == 4 { //subtopic should be events or data
			sTopic.Subtopic = arrTopics[0]
			sTopic.DeviceId = arrTopics[1]
			sTopic.DeviceGroup = arrTopics[2]
			sTopic.DeviceType = arrTopics[3]
		}
		if len(arrTopics) == 5 { //subtopic should be app.controls or device.controls
			sTopic.Subtopic = arrTopics[0] + "." + arrTopics[1]
			sTopic.DeviceId = arrTopics[2]
			sTopic.DeviceGroup = arrTopics[3]
			sTopic.DeviceType = arrTopics[4]
		}
	}
	return sTopic

}
