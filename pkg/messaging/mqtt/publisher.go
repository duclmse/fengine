package mqtt

import (
	"errors"
	"time"

	"github.com/duclmse/fengine/pkg/messaging"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ messaging.Publisher = (*publisher)(nil)

var errPublishTimeout = errors.New("failed to publish due to timeout reached")

type publisher struct {
	client  mqtt.Client
	timeout time.Duration
}

// NewPublisher returns a new MQTT message publisher.
func NewPublisher(address string, username string, password string, timeout time.Duration) (messaging.Publisher, error) {
	client, err := newPubClient(address, username, password, timeout)
	if err != nil {
		return nil, err
	}

	ret := publisher{
		client:  client,
		timeout: timeout,
	}
	return ret, nil
}

func (pub publisher) Publish(topic string, msg messaging.Message) error {
	token := pub.client.Publish(topic, byte(msg.Qos), msg.Retain, msg.Payload)
	if token.Error() != nil {
		return token.Error()
	}
	ok := token.WaitTimeout(pub.timeout)
	if ok && token.Error() != nil {
		return token.Error()
	}
	if !ok {
		return errPublishTimeout
	}
	return nil
}
