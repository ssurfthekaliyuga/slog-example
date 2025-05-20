package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slog-example/handlers"
	"slog-example/sl"
	"slog-example/usecases"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: nil,
	})

	logger := slog.New(&sl.ContextHandler{
		Handler: textHandler,
	})

	confLogger := logger.With(sl.Component("configuring"))

	confLogger.InfoContext(ctx, "start configuring", slog.Group("config",
		slog.Group("logger",
			slog.String("level", slog.LevelInfo.String()),
		),
		slog.String("http_address", "localhost:8080"),
	))

	hwuc, err := usecases.NewHelloWorld(logger)
	if err != nil {
		confLogger.ErrorContext(ctx, "create hello world", sl.Error(err))
		return
	}

	hwhd, err := handlers.NewHelloWorld(logger, hwuc)
	if err != nil {
		confLogger.ErrorContext(ctx, "create hello world", sl.Error(err))
		return
	}

	http.HandleFunc("/", hwhd.SayHello)

	go func() {
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			logger.ErrorContext(ctx, "cannot listen", sl.Error(err))
		}
	}()

	<-ctx.Done()
}
