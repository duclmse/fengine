package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/duclmse/fengine/fengine/db/sql"
	. "github.com/duclmse/fengine/viot"
)

const (
	envLogLevel      = "VT_FENGINE_LOG_LEVEL"
	envDBHost        = "VT_FENGINE_DB_HOST"
	envDBPort        = "VT_FENGINE_DB_PORT"
	envDBUser        = "VT_FENGINE_DB_USER"
	envDBPass        = "VT_FENGINE_DB_PASS"
	envDBName        = "VT_FENGINE_DB"
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
	defDBHost        = "localhost"
	defDBPort        = "5432"
	defDBUser        = "postgres"
	defDBPass        = "1"
	defDBName        = "fengine"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defCACerts       = ""
	defCacheURL      = "localhost:6379"
	defCachePass     = ""
	defCacheDB       = "0"
	defHTTPPort      = "8080"
	defAuthGRPCPort  = "8667"
	defServerCert    = ""
	defServerKey     = ""
	defJaegerURL     = ""
	defGrpcURL       = "viot-organizations:8230"
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
	GrpcServices map[string]GrpcService
}

type GrpcService struct {
	URL     string
	Timeout int64
}

type CacheConfig struct {
	URL  string
	Pass string
	DB   string
}

func LoadConfig(envFile string, grpcServices ...string) Config {
	err := LoadEnvFile(envFile)
	if err != nil {
		abs, _ := filepath.Abs(envFile)
		log.Fatalf("Cannot load env file from %s", abs)
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

func getGrpcConfig(services ...string) map[string]GrpcService {
	config := make(map[string]GrpcService, len(services))
	for _, service := range services {
		svc := strings.ToUpper(service)
		url := Env(fmt.Sprintf(envGrpcUrl, svc), defGrpcURL)
		timeoutCfg := fmt.Sprintf(envGrpcTimeout, svc)
		timeout, err := strconv.ParseInt(Env(timeoutCfg, defGrpcTimeout), 10, 64)
		if err != nil {
			log.Fatalf("Invalid %s value: %s", timeoutCfg, err.Error())
		}
		config[service] = GrpcService{URL: url, Timeout: timeout}
	}
	return config
}

func getDbConfig() sql.Config {
	return sql.Config{
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
