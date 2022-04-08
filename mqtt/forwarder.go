package mqttv2

import (
	"context"
	"errors"
	"fmt"
	"github.com/duclmse/fengine/pb"
	"hash/fnv"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/duclmse/fengine/pkg/common"
	log "github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/pkg/messaging"
	uuid2 "github.com/duclmse/fengine/pkg/uuid"
)

const (
	channels = "channels"
	messages = "messages"
)

var (
	channelRegExp         = regexp.MustCompile(`^/?messages/([\w-]+)(/[^?]*)?(\?.*)?$`)
	errMalformedTopic     = errors.New("malformed topic")
	errMalformedData      = errors.New("malformed request data")
	errMalformedSubtopic  = errors.New("malformed subtopic")
	errUnauthorizedAccess = errors.New("missing or invalid credentials provided")
	errNilClient          = errors.New("using nil client")
	errInvalidConnect     = errors.New("CONNECT request with invalid username or client ID")
	errNilTopicPub        = errors.New("PUBLISH to nil topic")
	errNilTopicSub        = errors.New("SUB to nil topic")
	ErrGetUserInfo        = errors.New("non-existent user")
)

// Forwarder specifies MQTT forwarder interface API.
type Forwarder interface {
	// Forward subscribes to the Subscriber and
	// publishes messages using provided Publisher.
	Forward(sub messaging.Subscriber, writer *kafka.Writer) error
}

type forwarder struct {
	topic                string
	dv                   pb.DevicesServiceClient
	logger               log.Logger
	msgWriterTopic       string
	varsParserTopic      string
	kafkaFEPartitionSize string
	kafkaTopicFlowEngine string
	us                   viot.UserServiceClient
	organizations        viot.OrganizationServiceClient
}

// NewForwarder returns new Forwarder implementation.
func NewForwarder(
	topic string, dv viot.DevicesServiceClient, logger log.Logger, msgWriterTopic, varsParserTopic, kafkaFEPartitionSize,
	kafkaTopicFlowEngine string, us viot.UserServiceClient, organizations viot.OrganizationServiceClient,
) Forwarder {
	return forwarder{
		topic:                topic,
		logger:               logger,
		dv:                   dv,
		us:                   us,
		msgWriterTopic:       msgWriterTopic,
		varsParserTopic:      varsParserTopic,
		kafkaFEPartitionSize: kafkaFEPartitionSize,
		kafkaTopicFlowEngine: kafkaTopicFlowEngine,
		organizations:        organizations,
	}
}

func (f forwarder) Forward(sub messaging.Subscriber, writer *kafka.Writer) error {
	fmt.Printf("Forward %s\n", f.topic)
	return sub.Subscribe(f.topic, f.handle(writer))
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func (f forwarder) handle(writer *kafka.Producer) messaging.MessageHandler {
	return func(msg messaging.Message) error {
		//fmt.Printf("Receive message from mqtt broker %s, publisher %s\n", msg.Protocol, msg.Publisher)
		topic := msg.Subtopic
		channelParts := channelRegExp.FindStringSubmatch(topic)
		if len(channelParts) < 1 {
			return nil
		}
		deviceID := channelParts[1]
		subtopic := channelParts[2]
		subtopic, err := parseSubtopic(subtopic) //Topic
		if err != nil {
			f.logger.Info("Error parsing subtopic: " + err.Error())
			return nil
		}
		fmt.Printf("PublishTopic: %s , deviceID %s\n", subtopic, deviceID)
		arDevice := &viot.InfoReqPriv{Id: deviceID}
		//fmt.Printf("Publish %v\n", arDevice)
		deviceInfo, err := f.dv.IdentifyInfoPrivate(context.Background(), arDevice)
		if err == nil {
			f.logger.Info(fmt.Sprintf("Got device info from device service %v", deviceInfo))
		} else {
			f.logger.Error(fmt.Sprintf("Error get device info %s", err))
			return nil
		}
		adminInfo, err := f.us.IdentifyUser(context.Background(), &viot.UserID{Value: deviceInfo.Owner})
		if err != nil {
			fmt.Printf("Error get user info %s\n", err)
		}
		projectInfo, err := f.organizations.IdentifyProject(context.Background(), &viot.IdentifyReq{Owner: deviceInfo.Owner, Id: deviceInfo.ProjectId})
		if err != nil {
			fmt.Printf("Error get project info %s\n", err)
		}
		userInfo, err := f.us.IdentifyUser(context.Background(), &viot.UserID{Value: deviceInfo.Createby})
		if err != nil {
			fmt.Printf("Error get user info %s\n", err)
		}
		userinfo := userInfo.GetPhone()
		if userinfo == "" {
			userinfo = userInfo.GetEmail()
		}
		var ruleChainId string
		if deviceInfo.GetTemplateId() != "" {
			templateInfo, err := f.dv.IdentifyTemplateInfo(context.Background(), &viot.IdentifyReq{
				Id:        deviceInfo.GetTemplateId(),
				Owner:     deviceInfo.GetOwner(),
				ProjectId: deviceInfo.GetProjectId(),
			})
			if err != nil {
				fmt.Printf("Get device template id %s got error %s\n ", deviceInfo.TemplateId, err)
			}
			if templateInfo != nil && templateInfo.GetRuleChainId() != "" {
				fmt.Printf("Got device template rule chain %s\n", templateInfo.GetRuleChainId())
				ruleChainId = templateInfo.GetRuleChainId()
			}
		}

		m := map[string]string{
			"entity_type":    common.DEVICE_TYPE,
			"entity_id":      deviceID,
			"deviceID":       deviceID,
			"topic":          subtopic,
			"protocol":       "mqtt",
			"createdby":      deviceInfo.Createby,
			"userInfo":       userinfo,
			"owner":          deviceInfo.Owner,
			"ownerInfo":      adminInfo.Email,
			"deviceName":     deviceInfo.Name,
			"group_name":     deviceInfo.Groupname,
			"group_id":       deviceInfo.Group,
			"devicemetadata": string(deviceInfo.GetMetadata()),
		}
		//metadataHeader := kafka.Header{
		//	Key: "metadata",
		//	Value: deviceInfo.GetMetadata(),
		//}
		headers := parseHeaders(m)
		//headers = append(headers,metadataHeader)
		//create message for flow engine

		//Check root rule chain of user
		if (ruleChainId != "" || projectInfo.RootRuleChain != "") && strings.Compare(subtopic, "events") != 0 {
			//create message for rule engine

			metaData := &messaging.TbMsgMetaDataProto{Data: m}
			uuid, err := uuid2.New().ID()
			if err != nil {
				fmt.Printf("error UUID: %s\n", err.Error())
			}
			tbMsgProto := &messaging.TbMsgProto{
				Id:                  uuid,
				RelationType:        "POST_VARIABLES",
				EntityType:          "DEVICE",
				EntityId:            deviceID,
				RuleChainId:         ruleChainId,
				RuleNodeId:          "",
				ClusterPartition:    10,
				MetaData:            metaData,
				DataType:            0,
				Data:                string(msg.Payload),
				Ts:                  time.Now().UnixNano(),
				RuleNodeExecCounter: 0,
			}

			TbMsg, err := tbMsgProto.Marshal()
			if err != nil {
				fmt.Printf("error marshal proto message %s\n", err.Error())
			}
			toRuleEngineMsg := &messaging.ToRuleEngineMsg{
				TenantId: deviceInfo.ProjectId,
				TbMsg:    TbMsg,
			}
			msgValue, err := toRuleEngineMsg.Marshal()
			if err != nil {
				fmt.Printf("error marshal proto message %s\n", err.Error())
			}

			// flow engine topic
			flowEnginePartitionSize, err := strconv.Atoi(f.kafkaFEPartitionSize)
			if err != nil {
				fmt.Printf("error parse flowEnginePartitionSize value %s to int %s\n", f.kafkaFEPartitionSize, err.Error())
			}
			partition := math.Abs(float64(int(hash(deviceID)) % flowEnginePartitionSize))
			flowEngineTopic := fmt.Sprintf("%s.main.%d", f.kafkaTopicFlowEngine, int(partition))

			toFlowEngineMsg := kafka.Message{
				Value: msgValue,
				TopicPartition: kafka.TopicPartition{
					Topic: &flowEngineTopic,
				},
				Key: []byte(uuid),
			}
			go func() {
				err := writer.WriteMessages(context.Background(), toFlowEngineMsg)
				if err != nil {
					fmt.Printf("send message to flow-engine got error %s \n", err.Error())
				} else {
					// log a confirmation once the message is written
					fmt.Printf("send message with payload %v to flow-engine with topic %s \n", string(msg.Payload), flowEngineTopic)
				}
			}()
		}

		//send message to core
		var messages []kafka.Message

		//toMsgWriterMsg := kafka.Message{
		//	Value:   msg.Payload,
		//	Headers: headers,
		//	Topic:   f.msgWriterTopic,
		//}
		//
		//messages = append(messages, toMsgWriterMsg)

		if subtopic == "events" || subtopic == "attribute" || subtopic == "attributets" || subtopic == "data" {
			toVarsParserMsg := kafka.Message{
				Value: msg.Payload,
				//Headers: parseHeaders(map[string]string{"deviceID": deviceID}),
				Headers: headers,
				Topic:   f.varsParserTopic,
			}
			messages = append(messages, toVarsParserMsg)
		}

		//Publish to Kafka

		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err = writer.WriteMessages(context.Background(), messages...)
		if err != nil {
			fmt.Printf("could not write message %s \n", err.Error())
		}
		// log a confirmation once the message is written
		fmt.Println("write success ")
		_, err = f.us.IncreaseMsg(context.Background(), &viot.UserID{Value: deviceInfo.Createby})
		if err != nil {
			fmt.Printf("err IncreaseMsg: %s", err)
		}
		_, err = f.organizations.IncreaseMsg(context.Background(), &viot.IncreaseReq{ProjectId: deviceInfo.ProjectId, Increase: true})
		if err != nil {
			fmt.Printf("err IncreaseMsg project: %s", err)
		}
		return nil
	}
}

func parseSubtopic(subtopic string) (string, error) {
	if subtopic == "" {
		return subtopic, nil
	}

	subtopic, err := url.QueryUnescape(subtopic)
	if err != nil {
		return "", errMalformedSubtopic
	}
	subtopic = strings.Replace(subtopic, "/", ".", -1)

	elems := strings.Split(subtopic, ".")
	filteredElems := []string{}
	for _, elem := range elems {
		if elem == "" {
			continue
		}
		if len(elem) > 1 && (strings.Contains(elem, "*") || strings.Contains(elem, ">")) {
			return "", errMalformedSubtopic
		}
		filteredElems = append(filteredElems, elem)
	}

	subtopic = strings.Join(filteredElems, ".")
	return subtopic, nil
}

func parseHeaders(keyvalue map[string]string) []kafka.Header {
	var res []kafka.Header
	var header kafka.Header
	for key, value := range keyvalue {
		header = kafka.Header{
			Key:   key,
			Value: []byte(value),
		}
		res = append(res, header)
	}

	return res
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0
	}
	return h.Sum32()
}
