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

	mux.Get("/fe/thing/:id", kit.NewServer(
		trace.TraceClient(tracer, "fe_get_thing")(getEntityEndpoint(svc, component)),
		decodeEntityRequest, encodeResponse, opts...))
	mux.Post("/fe/thing/:id", kit.NewServer(
		trace.TraceClient(tracer, "fe_get_thing")(upsertEntityEndpoint(svc, component)),
		decodeUpsertEntityRequest, encodeResponse, opts...))
	mux.Delete("/fe/thing/:id", kit.NewServer(
		trace.TraceClient(tracer, "fe_get_thing")(deleteEntityEndpoint(svc, component)),
		decodeEntityRequest, encodeResponse, opts...))

	mux.Get("/fe/thing/:id/service", kit.NewServer(
		trace.TraceClient(tracer, "fe_get_all_services")(getAllServicesEndpoint(svc, component)),
		decodeEntityRequest, encodeResponse, opts...))
	mux.Get("/fe/thing/:id/service/:service", kit.NewServer(
		trace.TraceClient(tracer, "fe_get_service")(getServiceEndpoint(svc, component)),
		decodeServiceRequest, encodeResponse, opts...))
	mux.Post("/fe/thing/:id/service/:service", kit.NewServer(
		trace.TraceClient(tracer, "fe_exec_service")(execServiceEndpoint(svc, component)),
		decodeExecRequest, encodeResponse, opts...))

	mux.GetFunc("/version", viot.Version("fengine"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
