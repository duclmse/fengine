package main

import (
	"fmt"
	"github.com/duclmse/fengine/viot"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	service "github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/fengine/api"
	grpc_api "github.com/duclmse/fengine/fengine/api/grpc"
	http_api "github.com/duclmse/fengine/fengine/api/http"
	"github.com/duclmse/fengine/fengine/db/sql"
	pb "github.com/duclmse/fengine/pb"
	_logger "github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/pkg/uuid"
)

//#region CONSTANT
const (
	defLogLevel      = "debug"
	defDBHost        = "localhost"
	defDBPort        = "5433"
	defDBUser        = "viot"
	defDBPass        = "viot"
	defDB            = "fengine"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defClientTLS     = "false"
	defCACerts       = ""
	defCacheURL      = "localhost:6380"
	defCachePass     = ""
	defCacheDB       = "0"
	defESURL         = "localhost:6380"
	defESPass        = ""
	defESDB          = "0"
	defHTTPPort      = "8080"
	defAuthHTTPPort  = "8669"
	defAuthGRPCPort  = "8667"
	defServerCert    = ""
	defServerKey     = ""
	defJaegerURL     = ""
	defUserURL       = "viot-organizations:8230"
	defUserTimeout   = "1"

	envLogLevel      = "VT_FENGINE_LOG_LEVEL"
	envDBHost        = "VT_FENGINE_DB_HOST"
	envDBPort        = "VT_FENGINE_DB_PORT"
	envDBUser        = "VT_FENGINE_DB_USER"
	envDBPass        = "VT_FENGINE_DB_PASS"
	envDB            = "VT_FENGINE_DB"
	envDBSSLMode     = "VT_FENGINE_DB_SSL_MODE"
	envDBSSLCert     = "VT_FENGINE_DB_SSL_CERT"
	envDBSSLKey      = "VT_FENGINE_DB_SSL_KEY"
	envDBSSLRootCert = "VT_FENGINE_DB_SSL_ROOT_CERT"
	envClientTLS     = "VT_FENGINE_CLIENT_TLS"
	envCACerts       = "VT_FENGINE_CA_CERTS"
	envCacheURL      = "VT_FENGINE_CACHE_URL"
	envCachePass     = "VT_FENGINE_CACHE_PASS"
	envCacheDB       = "VT_FENGINE_CACHE_DB"
	envESURL         = "VT_FENGINE_ES_URL"
	envESPass        = "VT_FENGINE_ES_PASS"
	envESDB          = "VT_FENGINE_ES_DB"
	envHTTPPort      = "VT_FENGINE_HTTP_PORT"
	envAuthHTTPPort  = "VT_FENGINE_AUTH_HTTP_PORT"
	envAuthGRPCPort  = "VT_FENGINE_AUTH_GRPC_PORT"
	envServerCert    = "VT_FENGINE_SERVER_CERT"
	envServerKey     = "VT_FENGINE_SERVER_KEY"
	envJaegerURL     = "JAEGER_URL"
	envUserURL       = "VT_USERS_GRPC_URL"
	envUserTimeout   = "VT_USER_GRPC_TIMEOUT"
)

//#endregion CONSTANT

type config struct {
	logLevel     string
	dbConfig     sql.Config
	clientTLS    bool
	caCerts      string
	cacheURL     string
	cachePass    string
	cacheDB      string
	esURL        string
	esPass       string
	esDB         string
	httpPort     string
	authHTTPort  string
	authGRPCPort string
	serverCert   string
	serverKey    string
	jaegerURL    string
	userURL      string
	userTimeout  time.Duration
}

func main() {
	cfg := loadConfig()

	logger, err := _logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	serviceTracer, closer := initJaeger("VTFEngine", cfg.jaegerURL, logger)
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			logger.Error("cannot close jaeger %s", err.Error())
		}
	}(closer)

	// Connect to Redis cache
	//cacheClient := connectToCache(cfg.cacheURL, cfg.cachePass, cfg.cacheDB, logger)
	//esClient := connectToCache(cfg.esURL, cfg.esPass, cfg.esDB, logger)

	cacheTracer, cacheCloser := initJaeger("fengine_cache", cfg.jaegerURL, logger)
	defer cacheCloser.Close()

	// Connect to Database
	//db := connectToDB(cfg.dbConfig, logger)
	//defer db.Close()

	serviceTracer, dbCloser := initJaeger("vtfengine_db", cfg.jaegerURL, logger)
	defer dbCloser.Close()

	// Connect to User service
	userConn := connectToUserService(cfg, logger)
	defer userConn.Close()

	// userTracer, userCloser := initJaeger("user", cfg.jaegerURL, logger)
	// defer userCloser.Close()
	// uc := usrapi.NewClient(userConn, userTracer, cfg.userTimeout)

	// Create FEngine Service
	svc := newService(serviceTracer, cacheTracer, logger) //db, cacheClient, esClient,

	errs := make(chan error, 2)
	// Start HTTP server
	go startHTTPServer(http_api.MakeHandler(serviceTracer, svc), cfg.httpPort, cfg, logger, errs)
	// At present, There is nothing api expose so stop grpc server.
	//go startGRPCServer(svc, serviceTracer, cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error("FEngine service terminated: %s", err)
}

func loadConfig() config {
	tls, err := strconv.ParseBool(viot.Env(envClientTLS, defClientTLS))
	if err != nil {
		log.Fatalf("Invalid value passed for %s\n", envClientTLS)
	}

	userTimeout, err := strconv.ParseInt(viot.Env(envUserTimeout, defUserTimeout), 10, 64)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", envUserTimeout, err.Error())
	}

	dbConfig := sql.Config{
		Host:        viot.Env(envDBHost, defDBHost),
		Port:        viot.Env(envDBPort, defDBPort),
		User:        viot.Env(envDBUser, defDBUser),
		Pass:        viot.Env(envDBPass, defDBPass),
		Name:        viot.Env(envDB, defDB),
		SSLMode:     viot.Env(envDBSSLMode, defDBSSLMode),
		SSLCert:     viot.Env(envDBSSLCert, defDBSSLCert),
		SSLKey:      viot.Env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: viot.Env(envDBSSLRootCert, defDBSSLRootCert),
	}

	return config{
		logLevel:     viot.Env(envLogLevel, defLogLevel),
		dbConfig:     dbConfig,
		clientTLS:    tls,
		caCerts:      viot.Env(envCACerts, defCACerts),
		cacheURL:     viot.Env(envCacheURL, defCacheURL),
		cachePass:    viot.Env(envCachePass, defCachePass),
		cacheDB:      viot.Env(envCacheDB, defCacheDB),
		esURL:        viot.Env(envESURL, defESURL),
		esPass:       viot.Env(envESPass, defESPass),
		esDB:         viot.Env(envESDB, defESDB),
		httpPort:     viot.Env(envHTTPPort, defHTTPPort),
		authHTTPort:  viot.Env(envAuthHTTPPort, defAuthHTTPPort),
		authGRPCPort: viot.Env(envAuthGRPCPort, defAuthGRPCPort),
		serverCert:   viot.Env(envServerCert, defServerCert),
		serverKey:    viot.Env(envServerKey, defServerKey),
		jaegerURL:    viot.Env(envJaegerURL, defJaegerURL),
		userURL:      viot.Env(envUserURL, defUserURL),
		userTimeout:  time.Duration(userTimeout) * time.Second,
	}
}

// init jaeger
func initJaeger(svcName, url string, logger _logger.Logger) (opentracing.Tracer, io.Closer) {
	if url == "" {
		return opentracing.NoopTracer{}, ioutil.NopCloser(nil)
	}

	tracer, closer, err := jconfig.Configuration{
		ServiceName: svcName,
		Sampler:     &jconfig.SamplerConfig{Type: "const", Param: 1},
		Reporter:    &jconfig.ReporterConfig{LocalAgentHostPort: url, LogSpans: true},
	}.NewTracer()
	if err != nil {
		logger.Error("Failed to init jaeger client: %s", err)
		os.Exit(1)
	}

	return tracer, closer
}

func connectToCache(cacheURL, cachePass, cacheDB string, logger _logger.Logger) *redis.Client {
	db, err := strconv.Atoi(cacheDB)
	if err != nil {
		logger.Error("Failed to connect to cache: %s", err)
		os.Exit(1)
	}

	return redis.NewClient(&redis.Options{
		Addr:     cacheURL,
		Password: cachePass,
		DB:       db,
	})
}

// connect to postgres DB
func connectToDB(dbConfig sql.Config, logger _logger.Logger) *sqlx.DB {
	fmt.Printf("db info: host: %s port: %s db: %s user: %s pass: %s \n", dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.User, dbConfig.Pass)
	db, err := sql.Connect(dbConfig)
	if err != nil {
		logger.Error("Failed to connect to postgres: %s", err)
		os.Exit(1)
	}
	return db
}

// Connect to User service
func connectToUserService(cfg config, logger _logger.Logger) *grpc.ClientConn {
	var opts []grpc.DialOption
	if cfg.clientTLS {
		if cfg.caCerts != "" {
			tpc, err := credentials.NewClientTLSFromFile(cfg.caCerts, "")
			if err != nil {
				logger.Error("Failed to create tls credentials: %s", err)
				os.Exit(1)
			}
			opts = append(opts, grpc.WithTransportCredentials(tpc))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		logger.Info("gRPC communication is not encrypted")
	}

	conn, err := grpc.Dial(cfg.userURL, opts...)
	if err != nil {
		logger.Error("Failed to connect to User service %s", err)
		os.Exit(1)
	} else {
		logger.Info(" Connected to User service by GPRC: %s", cfg.userURL)
	}
	return conn
}

// newService create new instantiate
func newService( /*usc viot.UserServiceClient,*/
	pTracer opentracing.Tracer,
	cacheTracer opentracing.Tracer,
	//db *sqlx.DB,
	//cacheClient *redis.Client,
	//esClient *redis.Client,
	logger _logger.Logger,
) service.Service {
	//database := sql.NewDatabase(db)
	//repo := sql.NewFEngineRepository(database)
	//repo = tracing.FEngineRepositoryMiddleware(pTracer, repo)

	//cache := cache.NewFEngineCache(cacheClient)
	//cache = tracing.FEngineCacheMiddleware(cacheTracer, cache)

	// Create new fengine service
	svc := service.New(uuid.New()) //repo, cache,
	//svc = cache.NewEventStoreMiddleware(svc, esClient)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "vtFEngine",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "vtFEngine",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

// Start HTTP Server
func startHTTPServer(handler http.Handler, port string, cfg config, logger _logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if cfg.serverCert != "" || cfg.serverKey != "" {
		logger.Info("vtFEngine service started using https on port %s with cert %s key %s",
			port, cfg.serverCert, cfg.serverKey)
		errs <- http.ListenAndServeTLS(p, cfg.serverCert, cfg.serverKey, handler)
		return
	}
	logger.Info("FEngine service started using HTTP on port %s", cfg.httpPort)
	errs <- http.ListenAndServe(p, handler)
}

//Start GRPC server
func startGRPCServer(svc service.Service, tracer opentracing.Tracer, cfg config, logger _logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.authGRPCPort)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error("Failed to listen on port %s: %s", cfg.authGRPCPort, err)
		os.Exit(1)
	}

	var server *grpc.Server
	if cfg.serverCert != "" || cfg.serverKey != "" {
		creds, err := credentials.NewServerTLSFromFile(cfg.serverCert, cfg.serverKey)
		if err != nil {
			logger.Error("Failed to load VTFEngine certificates: %s", err)
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("VT-FEngine gRPC service started using https on port %s with cert %s key %s",
			cfg.authGRPCPort, cfg.serverCert, cfg.serverKey))
		server = grpc.NewServer(grpc.Creds(creds))
	} else {
		logger.Info("VT-FEngine gRPC service started using http on port %s", cfg.authGRPCPort)
		server = grpc.NewServer()
	}
	pb.RegisterFEngineDataServer(server, grpc_api.NewDataServer(tracer, svc))
	pb.RegisterFEngineThingServer(server, grpc_api.NewThingServer(tracer, svc))
	errs <- server.Serve(listener)
}
