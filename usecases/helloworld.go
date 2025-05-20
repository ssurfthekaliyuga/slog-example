package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slog-example/sl"
)

type HelloWorld struct {
	logger *slog.Logger
}

func NewHelloWorld(logger *slog.Logger) (*HelloWorld, error) {
	if logger == nil {
		return nil, errors.New("logger must not be nil")
	}

	return &HelloWorld{
		logger: logger.With(sl.Component("usecases.HelloWorld")),
	}, nil
}

func (hw *HelloWorld) SayHello(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", errors.New("empty name")
	}

	hw.logger.InfoContext(ctx, "hello was sayed",
		slog.String("name", name),
	)

	return fmt.Sprintf("Hello %s", name), nil
}
