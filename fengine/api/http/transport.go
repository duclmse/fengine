package http

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"net/http"

	"github.com/duclmse/fengine/fengine"
	"github.com/duclmse/fengine/viot"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	kit "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	contentType = "application/json"
)

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	errInvalidQueryParams     = errors.New("invalid query params")
)

func MakeHandler(tracer opentracing.Tracer, svc fengine.Service) http.Handler {
	opts := []kit.ServerOption{
		kit.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/fe/exec", kit.NewServer(
		kitot.TraceClient(tracer, "fe_exec")(getEndpoint(svc)),
		decodeRequest, encodeResponse, opts...))

	r.GetFunc("/version", viot.Version("pricing"))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func serve(
	name string, tracer opentracing.Tracer, endpoint endpoint.Endpoint,
	decoder kit.DecodeRequestFunc, encoder kit.EncodeResponseFunc, opts ...kit.ServerOption,
) *kit.Server {
	return kit.NewServer(kitot.TraceClient(tracer, name)(endpoint), decoder, encoder, opts...)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
