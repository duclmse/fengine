package mqtt

import (
	"fmt"
	"time"

	"github.com/duclmse/fengine/pkg/errors"
	log "github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/pkg/messaging"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ messaging.Publisher = (*publisher)(nil)

var (
	errSubscribeTimeout   = errors.New("failed to subscribe due to timeout reached")
	errUnsubscribeTimeout = errors.New("failed to unsubscribe due to timeout reached")
)

type subscriber struct {
	client  mqtt.Client
	timeout time.Duration
	logger  log.Logger
}

// NewSubscriber returns a new MQTT message subscriber.
func NewSubscriber(address string, password string, timeout time.Duration, clientid string, logger log.Logger) (messaging.Subscriber, error) {
	client, err := newSubClient(address, password, clientid, timeout)
	if err != nil {
		return nil, err
	}

	ret := subscriber{
		client:  client,
		timeout: timeout,
		logger:  logger,
	}
	return ret, nil
}

func (sub subscriber) Subscribe(topic string, handler messaging.MessageHandler) error {
	token := sub.client.Subscribe(topic, qos, sub.mqttHandler(handler))
	if token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("Subscribe token %v\n", token)
	//ok := token.WaitTimeout(sub.timeout)
	ok := token.Wait()
	if ok && token.Error() != nil {
		return token.Error()
	}
	if !ok {
		return errSubscribeTimeout
	}
	return nil
}

func (sub subscriber) Unsubscribe(topic string) error {
	token := sub.client.Unsubscribe(topic)
	if token.Error() != nil {
		return token.Error()
	}
	ok := token.WaitTimeout(sub.timeout)
	if ok && token.Error() != nil {
		return token.Error()
	}
	if !ok {
		return errUnsubscribeTimeout
	}
	return nil
}

func (sub subscriber) mqttHandler(h messaging.MessageHandler) mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		//sub.logger.Warn(fmt.Sprintf("mqttHandler received message: %v", string(m.Payload())))

		var msg messaging.Message
		msg.Payload = m.Payload()
		msg.Subtopic = m.Topic()
		msg.Protocol = "mqtt"
		msg.Qos = int32(m.Qos())
		msg.Retain = m.Retained()
		msg.Created = time.Now().UnixNano()

		//if err := proto.Unmarshal(m.Payload(), &msg); err != nil {
		//	sub.logger.Warn(fmt.Sprintf("Failed to unmarshal received message: %s", err))
		//	return
		//}
		if err := h(msg); err != nil {
			sub.logger.Warn(fmt.Sprintf("Failed to handle Mainflux message: %s", err))
		}
	}
}
