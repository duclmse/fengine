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
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/subosito/gotenv"
	jconfig "github.com/uber/jaeger-client-go/config"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	service "github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/fengine/api"
	"github.com/duclmse/fengine/fengine/api/grpc"
	http_api "github.com/duclmse/fengine/fengine/api/http"
	"github.com/duclmse/fengine/fengine/db/cache"
	"github.com/duclmse/fengine/fengine/db/sql"
	"github.com/duclmse/fengine/fengine/tracing"
	pb "github.com/duclmse/fengine/pb"
	"github.com/duclmse/fengine/pkg/logger"
	. "github.com/duclmse/fengine/viot"
)

const (
	envLogLevel      = "VT_FENGINE_LOG_LEVEL"
	envDBUrl         = "VT_FENGINE_DB_URL"
	envDBHost        = "VT_FENGINE_DB_HOST"
	envDBPort        = "VT_FENGINE_DB_PORT"
	envDBUser        = "VT_FENGINE_DB_USER"
	envDBPass        = "VT_FENGINE_DB_PASS"
	envDBName        = "VT_FENGINE_DB_NAME"
	envDBSSLMode     = "VT_FENGINE_DB_SSL_MODE"
	envDBSSLCert     = "VT_FENGINE_DB_SSL_CERT"
	envDBSSLKey      = "VT_FENGINE_DB_SSL_KEY"
	envDBSSLRootCert = "VT_FENGINE_DB_SSL_ROOT_CERT"
	envCACerts       = "VT_FENGINE_CA_CERTS"
	envCacheURL      = "VT_FENGINE_CACHE_URL"
	envCachePass     = "VT_FENGINE_CACHE_PASS"
	envCacheDB       = "VT_FENGINE_CACHE_DB"
	envHTTPPort      = "VT_FENGINE_HTTP_PORT"
	envAuthGRPCPort  = "VT_FENGINE_AUTH_GRPC_PORT"
	envServerCert    = "VT_FENGINE_SERVER_CERT"
	envServerKey     = "VT_FENGINE_SERVER_KEY"
	envGrpcUrl       = "VT_%s_GRPC_URL"
	envGrpcTimeout   = "VT_%s_GRPC_TIMEOUT"
	envJaegerURL     = "JAEGER_URL"
)

const (
	defLogLevel      = "debug"
	defDBUrl         = ""
	defDBHost        = ""
	defDBPort        = ""
	defDBUser        = ""
	defDBPass        = ""
	defDBName        = ""
	defDBSSLMode     = ""
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defCACerts       = ""
	defCacheURL      = ""
	defCachePass     = ""
	defCacheDB       = ""
	defHTTPPort      = ""
	defAuthGRPCPort  = ""
	defServerCert    = ""
	defServerKey     = ""
	defJaegerURL     = ""
	defGrpcURL       = ""
	defGrpcTimeout   = "1"
)

//#endregion CONSTANT

type Config struct {
	DbConfig     sql.Config
	Cache        CacheConfig
	LogLevel     string
	CaCerts      string
	HttpPort     string
	AuthGRPCPort string
	ServerCert   string
	ServerKey    string
	JaegerURL    string
	GrpcServices map[string]grpc.GrpcService
}

type CacheConfig struct {
	URL  string
	Pass string
	DB   string
}

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

	db, err := sql.Connect(cfg.DbConfig, log)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %s", err)
	}
	defer db.Close()

	// Connect to User service
	executorConn := ConnectToGrpcService("executor", cfg, log)
	defer Close(log, "executor connection")(executorConn)

	executorTracer, executorCloser := InitJaeger("executor", cfg.JaegerURL, log)
	defer Close(log, "executor")(executorCloser)
	exeClient := grpc.NewExecutorClient(executorConn, executorTracer, cfg.GrpcServices["executor"])

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

func LoadConfig(envFile string, grpcServices ...string) Config {
	abs, _ := filepath.Abs(envFile)
	fmt.Printf("Loading env file from %s", abs)
	if err := gotenv.Load(envFile); err != nil {
		fmt.Printf("Cannot load env file from %s", abs)
	}
	return Config{
		DbConfig:     getDbConfig(),
		Cache:        getCacheConfig(),
		LogLevel:     Env(envLogLevel, defLogLevel),
		CaCerts:      Env(envCACerts, defCACerts),
		HttpPort:     Env(envHTTPPort, defHTTPPort),
		AuthGRPCPort: Env(envAuthGRPCPort, defAuthGRPCPort),
		ServerCert:   Env(envServerCert, defServerCert),
		ServerKey:    Env(envServerKey, defServerKey),
		JaegerURL:    Env(envJaegerURL, defJaegerURL),
		GrpcServices: getGrpcConfig(grpcServices...),
	}
}

func getGrpcConfig(services ...string) map[string]grpc.GrpcService {
	config := make(map[string]grpc.GrpcService, len(services))
	for _, svc := range services {
		svc := strings.ToUpper(svc)
		url := Env(fmt.Sprintf(envGrpcUrl, svc), defGrpcURL)
		timeoutCfg := fmt.Sprintf(envGrpcTimeout, svc)
		timeout, err := strconv.ParseInt(Env(timeoutCfg, defGrpcTimeout), 10, 32)
		if err != nil {
			log.Fatalf("Invalid %s value: %s", timeoutCfg, err.Error())
		}
		config[svc] = grpc.GrpcService{URL: url, Timeout: int(timeout)}
	}
	return config
}

func getDbConfig() sql.Config {
	return sql.Config{
		Url:         Env(envDBUrl, defDBUrl),
		Host:        Env(envDBHost, defDBHost),
		Port:        Env(envDBPort, defDBPort),
		User:        Env(envDBUser, defDBUser),
		Pass:        Env(envDBPass, defDBPass),
		Name:        Env(envDBName, defDBName),
		SSLMode:     Env(envDBSSLMode, defDBSSLMode),
		SSLCert:     Env(envDBSSLCert, defDBSSLCert),
		SSLKey:      Env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: Env(envDBSSLRootCert, defDBSSLRootCert),
	}
}

func getCacheConfig() CacheConfig {
	return CacheConfig{
		URL:  Env(envCacheURL, defCacheURL),
		Pass: Env(envCachePass, defCachePass),
		DB:   Env(envCacheDB, defCacheDB),
	}
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
		log.Fatalf("Failed to connect to cache: %s", err)
	}

	return redis.NewClient(&redis.Options{
		Addr:     cache.URL,
		Password: cache.Pass,
		DB:       db,
	})
}

func ConnectToGrpcService(name string, cfg Config, log logger.Logger) *ggrpc.ClientConn {
	var opts ggrpc.DialOption
	if cfg.CaCerts != "" {
		tpc, err := credentials.NewClientTLSFromFile(cfg.CaCerts, "")
		if err != nil {
			log.Error("Failed to create tls credentials: %s", err)
			os.Exit(1)
		}
		opts = ggrpc.WithTransportCredentials(tpc)
	} else {
		log.Info("gRPC communication is not encrypted")
		opts = ggrpc.WithTransportCredentials(insecure.NewCredentials())
	}

	url := cfg.GrpcServices[name].URL
	conn, err := ggrpc.Dial(url, opts)
	if err != nil {
		log.Fatalf("Failed to connect to %s service %s", name, err)
	}
	log.Info("Connected to %s service by GRPC: %s", name, url)
	return conn
}

// newService create new instantiate
func newService(component service.ServiceComponent) service.Service {
	repo := sql.NewFEngineRepository(component.DB, component.Log)
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

	var server *ggrpc.Server
	if cfg.ServerCert != "" || cfg.ServerKey != "" {
		credential, err := credentials.NewServerTLSFromFile(cfg.ServerCert, cfg.ServerKey)
		if err != nil {
			log.Error("Failed to load FEngine certificates: %s", err)
			os.Exit(1)
		}
		log.Info("FEngine gRPC service started on port %s (cert: %s; key %s)", port, cfg.ServerCert, cfg.ServerKey)
		server = ggrpc.NewServer(ggrpc.Creds(credential))
	} else {
		log.Info("FEngine gRPC service started on port %s", port)
		server = ggrpc.NewServer()
	}
	pb.RegisterFEngineDataServer(server, grpc.NewDataServer(tracer, svc))
	pb.RegisterFEngineThingServer(server, grpc.NewThingServer(tracer, svc))
	errs <- server.Serve(listener)
}
