package sl

import (
	"errors"
	"log/slog"
)

var (
	ErrorKey     = "error"
	PanicKey     = "panic"
	ComponentKey = "component"
	MethodKey    = "method"
	RequestIDKey = "request_id"
)

func Error(err error) slog.Attr {
	if err == nil {
		err = errors.New("sl.Error: got nil error")
	}
	return slog.String(ErrorKey, err.Error())
}

func Panic(p any) slog.Attr {
	return slog.Any("panic", p)
}

func Component(s string) slog.Attr {
	return slog.String("component", s)
}

func Method(s string) slog.Attr {
	return slog.String("method", s)
}

func RequestID(id string) slog.Attr {
	return slog.String("request_id", id)
}
