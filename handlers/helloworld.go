package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"net/url"
	"slog-example/sl"
)

type HelloWorldResponse struct {
	Msg string `json:"msg"`
}

type IHelloWorld interface {
	SayHello(ctx context.Context, name string) (string, error)
}

type HelloWorld struct {
	logger     *slog.Logger
	helloWorld IHelloWorld
}

func NewHelloWorld(logger *slog.Logger, helloWorld IHelloWorld) (*HelloWorld, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	if helloWorld == nil {
		return nil, errors.New("helloWorld is nil")
	}

	return &HelloWorld{
		logger:     logger.With(sl.Component("handlers.HelloWorld")),
		helloWorld: helloWorld,
	}, nil
}

func (hw *HelloWorld) SayHello(w http.ResponseWriter, r *http.Request) {
	ctx := sl.ContextWithAttrs(r.Context(), sl.RequestID(uuid.NewString())) //move it to middleware

	hw.logger.InfoContext(ctx, "got request", //move it to middleware
		slog.Group("http_request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("host", r.Host),
		),
	)

	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		writeError(ctx, w, hw.logger, err)
		return
	}

	msg, err := hw.helloWorld.SayHello(ctx, values.Get("name"))
	if err != nil {
		writeError(ctx, w, hw.logger, err)
		return
	}

	writeJSON(ctx, w, hw.logger, HelloWorldResponse{
		Msg: msg,
	})
}

func writeError(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(err.Error()))
	logger.ErrorContext(ctx, "cannot handle query", sl.Error(err))
}

func writeJSON(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(payload)
	if err != nil {
		writeError(ctx, w, logger, err)
		return
	}

	_, _ = w.Write(b)
}
