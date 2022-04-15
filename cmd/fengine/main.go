package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
	"github.com/duclmse/fengine/fengine/db/cache"
	"github.com/duclmse/fengine/fengine/db/sql"
	"github.com/duclmse/fengine/fengine/tracing"
	pb "github.com/duclmse/fengine/pb"
	_logger "github.com/duclmse/fengine/pkg/logger"
	"github.com/duclmse/fengine/pkg/uuid"
	. "github.com/duclmse/fengine/viot"
)

func main() {
	cfg := LoadConfig("./.env", "executor")

	logger, err := _logger.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	serviceTracer, closer := InitJaeger("VTFEngine", cfg.JaegerURL, logger)
	defer Close(logger, "jeager")(closer)

	cacheClient := ConnectToCache(cfg.Cache, logger)
	cacheTracer, cacheCloser := InitJaeger("fengine_cache", cfg.JaegerURL, logger)
	defer Close(logger, "cache")(cacheCloser)

	db := ConnectToDB(cfg.DbConfig, logger)
	defer Close(logger, "db")(db)

	serviceTracer, dbCloser := InitJaeger("vtfengine_db", cfg.JaegerURL, logger)
	defer Close(logger, "vtfengine_db")(dbCloser)

	// Connect to User service
	userConn := ConnectToGrpcService("user", cfg, logger)
	defer Close(logger, "user service")(userConn)

	// userTracer, userCloser := InitJaeger("user", cfg.JaegerURL, logger)
	// defer userCloser.Close()
	// uc := usrapi.NewClient(userConn, userTracer, cfg.UserTimeout)

	// Create FEngine Service
	svc := newService(serviceTracer, cacheClient, db, cacheTracer, logger) //db, , esClient,

	errs := make(chan error, 2)
	go startHTTPServer(http_api.MakeHandler(serviceTracer, svc), cfg.HttpPort, cfg, logger, errs)
	go startGRPCServer(svc, serviceTracer, cfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error("FEngine service terminated: %s", err)
}

func InitJaeger(svcName, url string, logger _logger.Logger) (opentracing.Tracer, io.Closer) {
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

func ConnectToCache(cache CacheConfig, logger _logger.Logger) *redis.Client {
	db, err := strconv.Atoi(cache.DB)
	if err != nil {
		logger.Error("Failed to connect to cache: %s", err)
		os.Exit(1)
	}

	return redis.NewClient(&redis.Options{
		Addr:     cache.URL,
		Password: cache.Pass,
		DB:       db,
	})
}

func ConnectToDB(dbCfg sql.Config, logger _logger.Logger) *sqlx.DB {
	logger.Info("db info: %s:%s/%s user: %s pass: %s", dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.User, dbCfg.Pass)
	db, err := sql.Connect(dbCfg)
	if err != nil {
		logger.Fatalf("Failed to connect to postgres: %s", err)
	}
	return db
}

func ConnectToGrpcService(name string, cfg Config, logger _logger.Logger) *grpc.ClientConn {
	var opts grpc.DialOption
	if cfg.CaCerts != "" {
		tpc, err := credentials.NewClientTLSFromFile(cfg.CaCerts, "")
		if err != nil {
			logger.Error("Failed to create tls credentials: %s", err)
			os.Exit(1)
		}
		opts = grpc.WithTransportCredentials(tpc)
	} else {
		logger.Info("gRPC communication is not encrypted")
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	url := cfg.GrpcServices[name].URL
	conn, err := grpc.Dial(url, opts)
	if err != nil {
		log.Panicf("Failed to connect to %s service %s", name, err)
	}
	logger.Info("Connected to %s service by GRPC: %s", name, url)
	return conn
}

// newService create new instantiate
func newService(
	pTracer opentracing.Tracer, cacheClient *redis.Client, db *sqlx.DB, cacheTracer opentracing.Tracer, logger _logger.Logger,
) service.Service {
	database := sql.NewDatabase(db)
	repo := sql.NewFEngineRepository(database)
	repo = tracing.FEngineRepositoryMiddleware(pTracer, repo)

	serviceCache := cache.NewFEngineCache(cacheClient)
	serviceCache = tracing.FEngineCacheMiddleware(cacheTracer, serviceCache)

	// Create new fengine service
	svc := service.New(uuid.New(), repo, serviceCache)
	//svc = serviceCache.NewEventStoreMiddleware(svc, serviceCache)
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
func startHTTPServer(handler http.Handler, port string, cfg Config, logger _logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if cfg.ServerCert != "" || cfg.ServerKey != "" {
		logger.Info("vtFEngine service started using https on port %s with cert %s key %s",
			port, cfg.ServerCert, cfg.ServerKey)
		errs <- http.ListenAndServeTLS(p, cfg.ServerCert, cfg.ServerKey, handler)
		return
	}
	logger.Info("FEngine HTTP service started on port %s", cfg.HttpPort)
	errs <- http.ListenAndServe(p, handler)
}

//Start GRPC server
func startGRPCServer(svc service.Service, tracer opentracing.Tracer, cfg Config, logger _logger.Logger, errs chan error) {
	port := fmt.Sprintf(":%s", cfg.AuthGRPCPort)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("Failed to listen on port %s: %s", cfg.AuthGRPCPort, err)
		os.Exit(1)
	}

	var server *grpc.Server
	if cfg.ServerCert != "" || cfg.ServerKey != "" {
		credential, err := credentials.NewServerTLSFromFile(cfg.ServerCert, cfg.ServerKey)
		if err != nil {
			logger.Error("Failed to load VTFEngine certificates: %s", err)
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("VT-FEngine gRPC service started using https on port %s with cert %s key %s",
			cfg.AuthGRPCPort, cfg.ServerCert, cfg.ServerKey))
		server = grpc.NewServer(grpc.Creds(credential))
	} else {
		logger.Info("FEngine gRPC service started on port %s", cfg.AuthGRPCPort)
		server = grpc.NewServer()
	}
	pb.RegisterFEngineDataServer(server, grpc_api.NewDataServer(tracer, svc))
	pb.RegisterFEngineThingServer(server, grpc_api.NewThingServer(tracer, svc))
	errs <- server.Serve(listener)
}
