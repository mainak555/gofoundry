package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"http/handlers"
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type StatusResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (r *StatusResponseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func KitHttpServerOptions(logger kitlog.Logger) []kithttp.ServerOption {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(handlers.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(EncodeErrorResponse),
		kithttp.ServerFinalizer(NewServerFinalizer(logger)),
	}
	return options
}

func EncodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}

	var errMsg string
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	//w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	if strings.HasPrefix(err.Error(), "mongo: no documents in result") {
		errMsg = "no documents found"
		w.WriteHeader(404)
	} else if strings.HasPrefix(err.Error(), "invalid") {
		errMsg = strings.TrimPrefix(err.Error(), "invalid")
		w.WriteHeader(400)
	} else if strings.HasPrefix(err.Error(), "unauthorize") {
		errMsg = strings.TrimPrefix(err.Error(), "unauthorize")
		w.WriteHeader(401)
	} else if strings.HasPrefix(err.Error(), "no content err") {
		errMsg = strings.TrimPrefix(err.Error(), "no content err")
		w.WriteHeader(404)
	} else if strings.HasPrefix(err.Error(), "no content") {
		errMsg = strings.TrimPrefix(err.Error(), "no content")
		w.WriteHeader(204)
	} else {
		errMsg = err.Error()
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": strings.TrimSpace(strings.TrimLeft(errMsg, ":")),
	})
	//w.Write([]byte(err.Error()))
}

func NewServerFinalizer(logger kitlog.Logger) kithttp.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		level.Info(logger).Log("status", code, "path", r.RequestURI, "method", r.Method)
	}
}
