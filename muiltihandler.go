package multihandler

import (
	"github.com/sirupsen/logrus"
	"io"
)

type Handler struct {
	Formatter logrus.Formatter
	Level     logrus.Level
	Writer    io.Writer
	Logger    *logrus.Logger
}

type MultiHandler struct {
	handlers []*Handler
}

type HandlerConfig struct {
	Hooks        logrus.LevelHooks
	ReportCaller bool
	ExitFunc     func(int)
	BufferPool   logrus.BufferPool
}

func NewMultiHandler(handler ...*Handler) *MultiHandler {
	return &MultiHandler{handlers: handler}
}

func NewHandler(formatter logrus.Formatter, level logrus.Level, writer io.Writer, config *HandlerConfig) *Handler {
	return &Handler{
		Formatter: formatter,
		Level:     level,
		Writer:    writer,
		Logger: &logrus.Logger{
			Out:          writer,
			Hooks:        config.Hooks,
			Formatter:    formatter,
			ReportCaller: config.ReportCaller,
			Level:        level,
			ExitFunc:     config.ExitFunc,
			BufferPool:   config.BufferPool,
		},
	}
}

func (h *MultiHandler) Format(e *logrus.Entry) ([]byte, error) {
	for _, handler := range h.handlers {
		if e.Level <= handler.Level {
			newEntry := e.Dup()
			newEntry.Logger = handler.Logger
			go newEntry.Log(e.Level, e.Message)
		}
	}
	return make([]byte, 0), nil
}
