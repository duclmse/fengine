package mqtt

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	protocol = "mqtt"
	idPub    = "mqtt-publisher"
	idSub    = "mqtt-subscriber"
	qos      = 0
	//Alltopics = "messages/06ea1181-e045-47fe-87de-4ba87cfe3983/data"
	Alltopics      = "messages/#"
	Ont2mqttTopics = "ont2mqtt/#"
)

var errConnect = errors.New("failed to connect to MQTT broker")

func newPubClient(address string, username string, password string, timeout time.Duration) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(address).
		SetUsername(username).
		SetPassword(password).
		SetClientID(username).
		SetCleanSession(false)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Error() != nil {
		return nil, token.Error()
	}

	ok := token.WaitTimeout(timeout)
	if ok && token.Error() != nil {
		return nil, token.Error()
	}
	if !ok {
		return nil, errConnect
	}

	return client, nil
}

func newSubClient(address string, password string, clientid string, timeout time.Duration) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(address).
		SetUsername(idSub).
		SetPassword(password).
		SetClientID(clientid).
		SetCleanSession(false)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Error() != nil {
		return nil, token.Error()
	}

	ok := token.WaitTimeout(timeout)
	if ok && token.Error() != nil {
		return nil, token.Error()
	}
	if !ok {
		return nil, errConnect
	}

	return client, nil
}
