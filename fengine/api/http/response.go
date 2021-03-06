package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/duclmse/fengine/viot"
)

func encodeResponse(_ context.Context, w http.ResponseWriter, response any) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(viot.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}
		w.WriteHeader(ar.Code())
		if ar.Empty() {
			return nil
		}
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Printf("err respond %s\n", err)
	}
	return err
}
