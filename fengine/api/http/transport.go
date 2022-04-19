package http

import (
	"context"
	"errors"
	"net/http"

	kitot "github.com/go-kit/kit/tracing/opentracing"
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

	r := bone.New()
	tracer := component.Tracer
	r.Post("/fe/exec", kit.NewServer(
		kitot.TraceClient(tracer, "fe_exec")(execEndpoint(svc, component)),
		decodeExecRequest, encodeExecResponse, opts...))

	r.GetFunc("/version", viot.Version("fengine"))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
