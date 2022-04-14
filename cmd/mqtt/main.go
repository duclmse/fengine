package mqtt

//
//import (
//	"fmt"
//	mqttv2 "github.com/duclmse/fengine/mqtt"
//	"github.com/duclmse/fengine/viot"
//	"io"
//	"io/ioutil"
//	"log"
//	"net"
//	"os"
//	"os/signal"
//	"strconv"
//	"syscall"
//	"time"
//
//	"github.com/go-redis/redis"
//	"github.com/opentracing/opentracing-go"
//	jconfig "github.com/uber/jaeger-client-go/config"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials"
//
//	//"github.com/viot/viot"
//	//devicesapi "github.com/viot/viot/devices/api/auth/grpc"
//	//"github.com/viot/viot/mqttv2"
//	//organizationsapi "github.com/viot/viot/organizations/api/grpc"
//	//mflog "github.com/viot/viot/pkg/logger"
//	//mqttsub "github.com/viot/viot/pkg/messaging/mqtt"
//	//usersapi "github.com/viot/viot/users/api/grpc"
//)
//
//const (
//	// Logging
//	defLogLevel = "error"
//	envLogLevel = "VT_MQTT_LOG_LEVEL"
//	// MQTT
//	defMQTTHost             = "0.0.0.0"
//	defMQTTPort             = "1885"
//	defMQTTTargetHost       = "0.0.0.0"
//	defMQTTTargetPort       = "1885"
//	defMQTTForwarderTimeout = "30"    // in seconds
//	defMQTTToken            = "token" // device token
//
//	envMQTTHost             = "VT_MQTT_SUB_HOST"
//	envMQTTPort             = "VT_MQTT_SUB_PORT"
//	envMQTTToken            = "VT_MQTT_SUB_TOKEN"
//	envMQTTTargetHost       = "VT_MQTT_TARGET_HOST"
//	envMQTTTargetPort       = "VT_MQTT_TARGET_PORT"
//	envMQTTForwarderTimeout = "VT_MQTT_FORWARDER_TIMEOUT"
//	// HTTP
//	defHTTPHost       = "0.0.0.0"
//	defHTTPPort       = "8080"
//	defHTTPScheme     = "ws"
//	defHTTPTargetHost = "localhost"
//	defHTTPTargetPort = "8080"
//	defHTTPTargetPath = "/vtmqtt"
//	envHTTPHost       = "VT_MQTT_WS_HOST"
//	envHTTPPort       = "VT_MQTT_WS_PORT"
//	envHTTPScheme     = "VT_MQTT_WS_SCHEMA"
//	envHTTPTargetHost = "VT_MQTT_WS_TARGET_HOST"
//	envHTTPTargetPort = "VT_MQTT_WS_TARGET_PORT"
//	envHTTPTargetPath = "VT_MQTT_WS_TARGET_PATH"
//
//	// Devices
//	defDevicesURL     = "localhost:8251"
//	defDevicesTimeout = "1" // in seconds
//	envDevicesURL     = "VT_DEVICES_AUTH_GRPC_URL"
//	envDevicesTimeout = "VT_DEVICES_AUTH_GRPC_TIMEOUT"
//
//	// Jaeger
//	defJaegerURL = ""
//	envJaegerURL = "JAEGER_URL"
//	// TLS
//	defClientTLS = "false"
//	defCACerts   = ""
//	envClientTLS = "VT_MQTT_CLIENT_TLS"
//	envCACerts   = "VT_MQTT_CA_CERTS"
//	// Instance
//	envInstance = "VT_MQTT_INSTANCE"
//	defInstance = ""
//	// ES
//	envESURL  = "VT_MQTT_ES_URL"
//	envESPass = "VT_MQTT_ES_PASS"
//	envESDB   = "VT_MQTT_ES_DB"
//	defESURL  = "localhost:6379"
//	defESPass = ""
//	defESDB   = "0"
//
//	//kafka
//	defKafkaURL                    = "viot-kafka.viot-kafka:9092"
//	defKafkaTopicMsgWriter         = "msgwriter-test1"
//	defKafkaTopicVarsParser        = "varsparser-test1"
//	defKafkaTopicPartition         = "10"
//	defKafkaTopicReplicationFactor = "1"
//
//	envKafkaURL                    = "KAFKA_URL"
//	envKafkaTopicMsgWriter         = "KAFKA_TOPIC_MSG_WRITER"
//	envKafkaTopicVarsParser        = "KAFKA_TOPIC_VARS_PARSER"
//	envKafkaTopicPartition         = "KAFKA_TOPIC_PARTITION"
//	envKafkaTopicReplicationFactor = "KAFKA_TOPIC_REPLICATION_FACTOR"
//
//	defKafkaTopicRuleEngine = "tb_rule_engine"
//	defKafkaPartitionSizeRE = "10"
//	envKafkaTopicRuleEngine = "KAFKA_TOPIC_RULE_ENGINE"
//	envKafkaPartitionSizeRE = "KAFKA_PARTITION_SIZE_RE"
//
//	//grpc users
//	defUsersURL     = "viot-users:8298"
//	defUsersTimeout = "5" // in seconds
//	envUsersURL     = "VT_USERS_AUTH_GRPC_URL"
//	envUsersTimeout = "VT_USERS_AUTH_GRPC_TIMEOUT"
//
//	//org grpc
//	defOrgGrpcURL     = "localhost:8230"
//	defOrgGrpcTimeout = "1" // in seconds
//	envOrgGrpcURL     = "VT_ORGANIZATIONS_GRPC_URL"
//	envOrgGrpcTimeout = "VT_ORGANIZATIONS_GRPC_TIMEOUT"
//)
//
//type config struct {
//	mqttHost             string
//	mqttPort             string
//	mqttToken            string
//	mqttTargetHost       string
//	mqttTargetPort       string
//	mqttForwarderTimeout time.Duration
//	httpHost             string
//	httpPort             string
//	httpScheme           string
//	httpTargetHost       string
//	httpTargetPort       string
//	httpTargetPath       string
//	jaegerURL            string
//	logLevel             string
//	devicesURL           string
//	devicesTimeout       time.Duration
//	clientTLS            bool
//	caCerts              string
//	instance             string
//	esURL                string
//	esPass               string
//	esDB                 string
//	//kafka
//	kafkaURL               string
//	kafkaTopicMsgWriter    string
//	kafkaTopicVarsParser   string
//	topicPartition         int
//	topicReplicationFactor int
//
//	kafkaTopicFlowEngine string
//	kafkaFEPartitionSize string
//
//	usersURL     string
//	usersTimeout time.Duration
//
//	orgGrpcURL     string
//	orgGrpcTimeout time.Duration
//}
//
//func main() {
//	cfg := loadConfig()
//
//	logger, err := mflog.New(os.Stdout, cfg.logLevel)
//	if err != nil {
//		log.Fatalf(err.Error())
//	}
//	deviceConn := connectToDevices(cfg, logger)
//	defer deviceConn.Close()
//
//	userConn := connectToUsers(cfg, logger)
//	defer userConn.Close()
//
//	orgConn := connectToOrg(cfg, logger)
//	defer orgConn.Close()
//
//	devicesTracer, devicesCloser := initJaeger("devices", cfg.jaegerURL, logger)
//	defer devicesCloser.Close()
//
//	usersTracer, usersCloser := initJaeger("users", cfg.jaegerURL, logger)
//	defer usersCloser.Close()
//
//	orgTracer, orgCloser := initJaeger("organizations", cfg.jaegerURL, logger)
//	defer orgCloser.Close()
//
//	rc := connectToRedis(cfg.esURL, cfg.esPass, cfg.esDB, logger)
//	defer rc.Close()
//
//	dv := devicesapi.NewClient(deviceConn, devicesTracer, cfg.devicesTimeout)
//
//	us := usersapi.NewClient(userConn, usersTracer, cfg.usersTimeout)
//
//	orgs := organizationsapi.NewClient(orgConn, orgTracer, cfg.orgGrpcTimeout)
//
//	//=======================mqtt subscriber=======================//
//	ms, err := mqttsub.NewSubscriber(fmt.Sprintf("%s:%s", cfg.mqttHost, cfg.mqttPort), cfg.mqttToken, cfg.mqttForwarderTimeout, "mqtt-subscriber", logger)
//	logger.Info("Connect to mqtt subsriber: %s:%s\n", cfg.mqttHost, cfg.mqttPort)
//	if err != nil {
//		logger.Error("Failed to create MQTT subsriber: %s", err)
//		os.Exit(1)
//	}
//
//	topicConfigs := []kafka.TopicConfig{
//		{
//			Topic:             cfg.kafkaTopicMsgWriter,
//			NumPartitions:     cfg.topicPartition,
//			ReplicationFactor: cfg.topicReplicationFactor,
//		},
//		{
//			Topic:             cfg.kafkaTopicVarsParser,
//			NumPartitions:     cfg.topicPartition,
//			ReplicationFactor: cfg.topicReplicationFactor,
//		},
//	}
//
//	createKafkaTopic(cfg.kafkaURL, topicConfigs...)
//
//	//intialize the kafka writer with the broker addresses, and the topic
//	w := &kafka.Writer{
//		//Balancer: &kafka.Hash{},
//		Addr: kafka.TCP(cfg.kafkaURL),
//		//Async: true,
//		//BatchBytes: 102400,
//		BatchTimeout: 100 * time.Millisecond,
//		RequiredAcks: -1,
//	}
//
//	//w, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.kafkaURL})
//	//if err != nil {
//	//	panic(err)
//	//}
//	//defer w.Close()
//
//	fwd := mqttv2.NewForwarder(mqttsub.Alltopics, dv, logger, cfg.kafkaTopicMsgWriter, cfg.kafkaTopicVarsParser, cfg.kafkaFEPartitionSize, cfg.kafkaTopicFlowEngine, us, orgs)
//	if err := fwd.Forward(ms, w); err != nil {
//		logger.Error("Failed to start MQTT forwarder %s", err)
//		os.Exit(1)
//	}
//
//	errs := make(chan error, 2)
//
//	go func() {
//		c := make(chan os.Signal, 1)
//		signal.Notify(c, syscall.SIGINT)
//		errs <- fmt.Errorf("%s", <-c)
//	}()
//
//	err = <-errs
//	logger.Error("mqttv2 terminated: %s", err)
//}
//
//func loadConfig() config {
//	tls, err := strconv.ParseBool(viot.Env(envClientTLS, defClientTLS))
//	if err != nil {
//		log.Fatalf("Invalid value passed for %s\n", envClientTLS)
//	}
//
//	mqttTimeout, err := strconv.ParseInt(viot.Env(envMQTTForwarderTimeout, defMQTTForwarderTimeout), 10, 64)
//	if err != nil {
//		log.Fatalf("Invalid %s value: %s", envMQTTForwarderTimeout, err.Error())
//	}
//
//	devicesTimeout, err := strconv.ParseInt(viot.Env(envDevicesTimeout, defDevicesTimeout), 10, 64)
//	if err != nil {
//		log.Fatalf("Invalid %s value: %s", envDevicesTimeout, err.Error())
//	}
//
//	topicPartition, err := strconv.ParseInt(viot.Env(envKafkaTopicPartition, defKafkaTopicPartition), 10, 64)
//	if err != nil {
//		log.Fatalf("Invalid %s value: %s", topicPartition, err.Error())
//	}
//
//	topicReplicationFactor, err := strconv.ParseInt(viot.Env(envKafkaTopicReplicationFactor, defKafkaTopicReplicationFactor), 10, 64)
//	if err != nil {
//		log.Fatalf("Invalid %s value: %s", topicReplicationFactor, err.Error())
//	}
//
//	usersTimeout, err := strconv.ParseInt(viot.Env(envUsersTimeout, defUsersTimeout), 10, 64)
//	if err != nil {
//		log.Fatalf("Invalid %s value: %s", envUsersTimeout, err.Error())
//	}
//
//	orgTimeout, err := strconv.ParseInt(viot.Env(envOrgGrpcTimeout, defOrgGrpcTimeout), 10, 64)
//	if err != nil {
//		log.Fatalf("Invalid %s value: %s", envOrgGrpcTimeout, err.Error())
//	}
//
//	return config{
//		mqttHost:             viot.Env(envMQTTHost, defMQTTHost),
//		mqttPort:             viot.Env(envMQTTPort, defMQTTPort),
//		mqttToken:            viot.Env(envMQTTToken, defMQTTToken),
//		mqttTargetHost:       viot.Env(envMQTTTargetHost, defMQTTTargetHost),
//		mqttTargetPort:       viot.Env(envMQTTTargetPort, defMQTTTargetPort),
//		mqttForwarderTimeout: time.Duration(mqttTimeout) * time.Second,
//		httpHost:             viot.Env(envHTTPHost, defHTTPHost),
//		httpPort:             viot.Env(envHTTPPort, defHTTPPort),
//		httpScheme:           viot.Env(envHTTPScheme, defHTTPScheme),
//		httpTargetHost:       viot.Env(envHTTPTargetHost, defHTTPTargetHost),
//		httpTargetPort:       viot.Env(envHTTPTargetPort, defHTTPTargetPort),
//		httpTargetPath:       viot.Env(envHTTPTargetPath, defHTTPTargetPath),
//		jaegerURL:            viot.Env(envJaegerURL, defJaegerURL),
//		devicesURL:           viot.Env(envDevicesURL, defDevicesURL),
//		devicesTimeout:       time.Duration(devicesTimeout) * time.Second,
//		logLevel:             viot.Env(envLogLevel, defLogLevel),
//		clientTLS:            tls,
//		caCerts:              viot.Env(envCACerts, defCACerts),
//		instance:             viot.Env(envInstance, defInstance),
//		esURL:                viot.Env(envESURL, defESURL),
//		esPass:               viot.Env(envESPass, defESPass),
//		esDB:                 viot.Env(envESDB, defESDB),
//		usersURL:             viot.Env(envUsersURL, defUsersURL),
//		usersTimeout:         time.Duration(usersTimeout) * time.Second,
//
//		//kafka
//		kafkaURL:               viot.Env(envKafkaURL, defKafkaURL),
//		kafkaTopicMsgWriter:    viot.Env(envKafkaTopicMsgWriter, defKafkaTopicMsgWriter),
//		kafkaTopicVarsParser:   viot.Env(envKafkaTopicVarsParser, defKafkaTopicVarsParser),
//		topicPartition:         int(topicPartition),
//		topicReplicationFactor: int(topicReplicationFactor),
//		kafkaTopicFlowEngine:   viot.Env(envKafkaTopicRuleEngine, defKafkaTopicRuleEngine),
//		kafkaFEPartitionSize:   viot.Env(envKafkaPartitionSizeRE, defKafkaPartitionSizeRE),
//		orgGrpcURL:             viot.Env(envOrgGrpcURL, defOrgGrpcURL),
//		orgGrpcTimeout:         time.Duration(orgTimeout) * time.Second,
//	}
//}
//
//func initJaeger(svcName, url string, logger mflog.Logger) (opentracing.Tracer, io.Closer) {
//	if url == "" {
//		return opentracing.NoopTracer{}, ioutil.NopCloser(nil)
//	}
//
//	tracer, closer, err := jconfig.Configuration{
//		ServiceName: svcName,
//		Sampler: &jconfig.SamplerConfig{
//			Type:  "const",
//			Param: 1,
//		},
//		Reporter: &jconfig.ReporterConfig{
//			LocalAgentHostPort: url,
//			LogSpans:           true,
//		},
//	}.NewTracer()
//	if err != nil {
//		logger.Error("Failed to init Jaeger client: %s", err)
//		os.Exit(1)
//	}
//
//	return tracer, closer
//}
//
//func connectToRedis(redisURL, redisPass, redisDB string, logger mflog.Logger) *redis.Client {
//	db, err := strconv.Atoi(redisDB)
//	if err != nil {
//		logger.Error("Failed to connect to redis: %s", err)
//		os.Exit(1)
//	}
//
//	return redis.NewClient(&redis.Options{
//		Addr:     redisURL,
//		Password: redisPass,
//		DB:       db,
//	})
//}
//
//func connectToDevices(cfg config, logger mflog.Logger) *grpc.ClientConn {
//	var opts []grpc.DialOption
//	if cfg.clientTLS {
//		if cfg.caCerts != "" {
//			tpc, err := credentials.NewClientTLSFromFile(cfg.caCerts, "")
//			if err != nil {
//				logger.Error("Failed to create tls credentials: %s", err)
//				os.Exit(1)
//			}
//			opts = append(opts, grpc.WithTransportCredentials(tpc))
//		}
//	} else {
//		opts = append(opts, grpc.WithInsecure())
//		logger.Info("gRPC communication with devices service is not encrypted")
//	}
//	logger.Info("Connect to devices url %s", cfg.devicesURL)
//	conn, err := grpc.Dial(cfg.devicesURL, opts...)
//	if err != nil {
//		logger.Error("Failed to connect to devices service: %s", err)
//		os.Exit(1)
//	}
//
//	return conn
//}
//
//func createKafkaTopic(kafkaUrl string, topicConfig ...kafka.TopicConfig) {
//	conn, err := kafka.Dial("tcp", kafkaUrl)
//	if err != nil {
//		fmt.Printf("Error dial tcp to Kafka server %s\n", kafkaUrl)
//	}
//	defer conn.Close()
//
//	controller, err := conn.Controller()
//	if err != nil {
//		panic(err.Error())
//	}
//	var controllerConn *kafka.Conn
//	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
//	if err != nil {
//		panic(err.Error())
//	}
//	defer controllerConn.Close()
//
//	err = controllerConn.CreateTopics(topicConfig...)
//	if err != nil {
//		fmt.Printf("Error create topic [%v] to Kafka server %s\n", topicConfig, kafkaUrl)
//	}
//}
//
//func connectToUsers(cfg config, logger logger.Logger) *grpc.ClientConn {
//	var opts []grpc.DialOption
//	if cfg.clientTLS {
//		if cfg.caCerts != "" {
//			tpc, err := credentials.NewClientTLSFromFile(cfg.caCerts, "")
//			if err != nil {
//				logger.Error("Failed to create tls credentials: %s", err)
//				os.Exit(1)
//			}
//			opts = append(opts, grpc.WithTransportCredentials(tpc))
//		}
//	} else {
//		opts = append(opts, grpc.WithInsecure())
//		logger.Info("gRPC communication with users service is not encrypted")
//	}
//	fmt.Printf("Connect to users url %s\n", cfg.usersURL)
//	conn, err := grpc.Dial(cfg.usersURL, opts...)
//	if err != nil {
//		logger.Error("Failed to connect to users service: %s", err)
//		os.Exit(1)
//	} else {
//		logger.Info("Connected to users service by gRPC : %s", cfg.usersURL)
//	}
//
//	return conn
//}
//
//func connectToOrg(cfg config, logger logger.Logger) *grpc.ClientConn {
//	var opts []grpc.DialOption
//	if cfg.clientTLS {
//		if cfg.caCerts != "" {
//			tpc, err := credentials.NewClientTLSFromFile(cfg.caCerts, "")
//			if err != nil {
//				logger.Error("Failed to create tls credentials: %s", err)
//				os.Exit(1)
//			}
//			opts = append(opts, grpc.WithTransportCredentials(tpc))
//		}
//	} else {
//		opts = append(opts, grpc.WithInsecure())
//		logger.Info("gRPC communication is not encrypted")
//	}
//
//	conn, err := grpc.Dial(cfg.orgGrpcURL, opts...)
//	if err != nil {
//		logger.Error("Failed to connect to organization service: %s", err)
//		os.Exit(1)
//	} else {
//		logger.Info("Connected to organization service by gRPC : %s", cfg.orgGrpcURL)
//	}
//
//	return conn
//}
