package mylogger

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger //nolint:gochecknoglobals
	once   sync.Once          //nolint:gochecknoglobals

	ErrNoLogger = errors.New("my_logger.Get fot not existing logger")
)

type MyLogger struct {
	*zap.SugaredLogger
}

func NewNop() *MyLogger {
	once.Do(func() {
		logger = zap.NewNop().Sugar()
	})

	return &MyLogger{logger}
}

func New(
	outputPaths []string,
	errorOutputPaths []string,
	productionMode bool,
	options ...zap.Option,
) (*MyLogger, error) {
	var err error
	var logger *zap.SugaredLogger

	once.Do(func() {
		var config zap.Config

		if productionMode {
			config = zap.NewProductionConfig()
		} else {
			config = zap.NewDevelopmentConfig()
		}

		config.OutputPaths = outputPaths
		config.ErrorOutputPaths = errorOutputPaths

		zapLogger, innerErr := config.Build(options...)
		if innerErr != nil {
			err = innerErr

			return
		}

		logger = zapLogger.Sugar()
	})

	if err != nil {
		return nil, err
	}

	return &MyLogger{logger}, nil
}

func Get() (*MyLogger, error) {
	if logger == nil {
		fmt.Println(ErrNoLogger)

		return nil, ErrNoLogger
	}

	return &MyLogger{logger}, nil
}

func (m *MyLogger) WithCtx(ctx context.Context) *MyLogger {
	return &MyLogger{logger.With(zap.String("request_id", middleware.GetReqID(ctx)))}
}
