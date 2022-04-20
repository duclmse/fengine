package http

import (
	"context"
	"errors"
	"net/http"

	trace "github.com/go-kit/kit/tracing/opentracing"
	kit "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/viot"
)

const (
	contentType = "application/json"
)

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	errInvalidQueryParams     = errors.New("invalid query params")
)

func MakeHandler(svc fengine.Service, component fengine.ServiceComponent) http.Handler {
	opts := []kit.ServerOption{
		kit.ServerErrorEncoder(encodeError),
	}

	mux := bone.New()
	tracer := component.Tracer

	mux.Get("/fe/thing/:id/service", kit.NewServer(
		trace.TraceClient(tracer, "fe_exec")(execEndpoint(svc, component)),
		decodeAllServiceRequest, encodeAllServiceResponse, opts...))
	mux.Get("/fe/thing/:id/service/:service", kit.NewServer(
		trace.TraceClient(tracer, "fe_exec")(execEndpoint(svc, component)),
		decodeServiceRequest, encodeServiceResponse, opts...))
	mux.Post("/fe/exec", kit.NewServer(
		trace.TraceClient(tracer, "fe_exec")(execEndpoint(svc, component)),
		decodeExecRequest, encodeExecResponse, opts...))

	mux.GetFunc("/version", viot.Version("fengine"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
