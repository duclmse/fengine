package viot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/goccy/go-json"

	"github.com/duclmse/fengine/pkg/logger"
)

const version string = "0.10.0"

// VersionInfo contains version endpoint response.
type VersionInfo struct {
	// Service contains service name.
	Service string `json:"service"`
	// Version contains service current version value.
	Version string `json:"version"`
}

// Version exposes an HTTP handler for retrieving service version.
func Version(service string) http.HandlerFunc {
	return func(rw http.ResponseWriter, _ *http.Request) {
		res := VersionInfo{service, version}
		data, _ := json.Marshal(res)
		_, _ = rw.Write(data)
	}
}

// Env reads specified environment variable. If no value has been found, fallback is returned.
func Env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}

// Response contains HTTP response specific methods.
type Response interface {
	// Code returns HTTP response code.
	Code() int
	// Headers returns map of HTTP headers with their values.
	Headers() map[string]string
	// Empty indicates if HTTP response has content.
	Empty() bool
}

func Close(log logger.Logger, name string) func(io.Closer) {
	return func(closer io.Closer) {
		if err := closer.Close(); err != nil {
			if log == nil {
				fmt.Printf("cannot close %s: %s\n", name, err.Error())
			} else {
				log.Error("cannot close %s: %s", name, err.Error())
			}
		}
	}
}

func CloseCtx(
	closer interface {
		Close(ctx context.Context) error
	},
	ctx context.Context,
) func() {
	return func() {
		if err := closer.Close(ctx); err != nil {
			fmt.Printf("cannot close %s\n", err.Error())
		}
	}
}
