package nats

import (
	"fmt"
	"github.com/duclmse/fengine/pkg/messaging"
	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
	"log"
)

var _ messaging.Publisher = (*stanpublisher)(nil)

type stanpublisher struct {
	conn stan.Conn
}

// NewStanPublisher returns NATS streaming server message Publisher.
func NewStanPublisher(url, clusterid, clientid string) (Publisher, error) {
	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		fmt.Printf("Can not connect to nats: %s. Error %v \n", url, err)
		return nil, err
	}
	conn, err := stan.Connect(clusterid, clientid, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		fmt.Printf("Can not connect to stan: [%s] [%s]. Error: %v  \n", clusterid, clientid, err)
		return nil, err
	}
	ret := &stanpublisher{
		conn: conn,
	}
	return ret, nil
}

func (pub *stanpublisher) Publish(topic string, msg messaging.Message) error {
	fmt.Printf("Forward message to STAN with topic: %s \n", topic)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	if err := pub.conn.Publish(topic, data); err != nil {
		return err
	}

	return nil
}

func (pub *stanpublisher) Close() {
	pub.conn.NatsConn().Close()
	pub.conn.Close()
}
