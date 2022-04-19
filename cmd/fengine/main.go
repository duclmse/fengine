package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-redis/redis/v8"
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
	"github.com/duclmse/fengine/pkg/logger"
	. "github.com/duclmse/fengine/viot"
)

func main() {
	cfg := LoadConfig("./.env", "executor")

	log, err := logger.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	serviceTracer, closer := InitJaeger("VTFEngine", cfg.JaegerURL, log)
	defer Close(log, "jeager")(closer)

	cacheClient := ConnectToCache(cfg.Cache, log)
	cacheTracer, cacheCloser := InitJaeger("fengine_cache", cfg.JaegerURL, log)
	defer Close(log, "cache")(cacheCloser)

	db := ConnectToDB(cfg.DbConfig, log)
	defer Close(log, "db")(db)

	// Connect to User service
	executorConn := ConnectToGrpcService("executor", cfg, log)
	defer Close(log, "executor connection")(executorConn)

	executorTracer, executorCloser := InitJaeger("executor", cfg.JaegerURL, log)
	defer Close(log, "executor")(executorCloser)
	exeClient := grpc_api.NewExecutorClient(executorConn, executorTracer, cfg.GrpcServices["executor"])

	// Create FEngine Service
	components := service.ServiceComponent{
		Tracer:      serviceTracer,
		Cache:       cacheClient,
		CacheTracer: cacheTracer,
		ExeClient:   exeClient,
		DB:          db,
		Log:         log,
	}
	svc := newService(components)

	errs := make(chan error, 2)
	go startHTTPServer(http_api.MakeHandler(svc, components), cfg, log, errs)
	go startGRPCServer(svc, serviceTracer, cfg, log, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	log.Error("FEngine service terminated: %s", err)
}

func InitJaeger(svcName, url string, log logger.Logger) (opentracing.Tracer, io.Closer) {
	if url == "" {
		return opentracing.NoopTracer{}, ioutil.NopCloser(nil)
	}

	tracer, closer, err := jconfig.Configuration{
		ServiceName: svcName,
		Sampler:     &jconfig.SamplerConfig{Type: "const", Param: 1},
		Reporter:    &jconfig.ReporterConfig{LocalAgentHostPort: url, LogSpans: true},
	}.NewTracer()
	if err != nil {
		log.Error("Failed to init jaeger client: %s", err)
		os.Exit(1)
	}

	return tracer, closer
}

func ConnectToCache(cache CacheConfig, log logger.Logger) *redis.Client {
	db, err := strconv.Atoi(cache.DB)
	if err != nil {
		log.Error("Failed to connect to cache: %s", err)
		os.Exit(1)
	}

	return redis.NewClient(&redis.Options{
		Addr:     cache.URL,
		Password: cache.Pass,
		DB:       db,
	})
}

func ConnectToDB(cfg sql.Config, log logger.Logger) *sqlx.DB {
	log.Info("db info: %s:%s/%s user: %s pass: %s", cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Pass)
	db, err := sql.Connect(cfg, log)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %s", err)
	}
	return db
}

func ConnectToGrpcService(name string, cfg Config, log logger.Logger) *grpc.ClientConn {
	var opts grpc.DialOption
	if cfg.CaCerts != "" {
		tpc, err := credentials.NewClientTLSFromFile(cfg.CaCerts, "")
		if err != nil {
			log.Error("Failed to create tls credentials: %s", err)
			os.Exit(1)
		}
		opts = grpc.WithTransportCredentials(tpc)
	} else {
		log.Info("gRPC communication is not encrypted")
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	url := cfg.GrpcServices[name].URL
	conn, err := grpc.Dial(url, opts)
	if err != nil {
		log.Fatalf("Failed to connect to %s service %s", name, err)
	}
	log.Info("Connected to %s service by GRPC: %s", name, url)
	return conn
}

// newService create new instantiate
func newService(component service.ServiceComponent) service.Service {
	database := sql.NewDatabase(component.DB)
	repo := sql.NewFEngineRepository(database)
	repo = tracing.FEngineRepositoryMiddleware(component.Tracer, repo)

	serviceCache := cache.NewFEngineCache(component.Cache)
	serviceCache = tracing.FEngineCacheMiddleware(component.CacheTracer, serviceCache)

	// Create new fengine service
	svc := service.FengineService{
		Repository: repo,
		Cache:      serviceCache,
		ExecClient: component.ExeClient,
		Log:        component.Log,
	}.New()
	//svc = serviceCache.NewEventStoreMiddleware(svc, serviceCache)
	svc = api.LoggingMiddleware(svc, component.Log)
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
func startHTTPServer(handler http.Handler, cfg Config, log logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.HttpPort)
	if cfg.ServerCert != "" || cfg.ServerKey != "" {
		log.Info("FEngine HTTP service started on port %s (cert: %s; key: %s)", p, cfg.ServerCert, cfg.ServerKey)
		errs <- http.ListenAndServeTLS(p, cfg.ServerCert, cfg.ServerKey, handler)
		return
	}
	log.Info("FEngine HTTP service started on port %s", p)
	errs <- http.ListenAndServe(p, handler)
}

//Start GRPC server
func startGRPCServer(svc service.Service, tracer opentracing.Tracer, cfg Config, log logger.Logger, errs chan error) {
	port := fmt.Sprintf(":%s", cfg.AuthGRPCPort)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Error("Failed to listen on port %s: %s", port, err)
		os.Exit(1)
	}

	var server *grpc.Server
	if cfg.ServerCert != "" || cfg.ServerKey != "" {
		credential, err := credentials.NewServerTLSFromFile(cfg.ServerCert, cfg.ServerKey)
		if err != nil {
			log.Error("Failed to load FEngine certificates: %s", err)
			os.Exit(1)
		}
		log.Info("FEngine gRPC service started on port %s (cert: %s; key %s)", port, cfg.ServerCert, cfg.ServerKey)
		server = grpc.NewServer(grpc.Creds(credential))
	} else {
		log.Info("FEngine gRPC service started on port %s", port)
		server = grpc.NewServer()
	}
	pb.RegisterFEngineDataServer(server, grpc_api.NewDataServer(tracer, svc))
	pb.RegisterFEngineThingServer(server, grpc_api.NewThingServer(tracer, svc))
	errs <- server.Serve(listener)
}
