package applog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adnanahmady/go-rest-api-blog/config"
	"github.com/adnanahmady/go-rest-api-blog/pkg/app"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string, args ...any)
	With(args ...any) Logger
}

type AppLog struct {
	lgr zerolog.Logger
}

func NewLog(cfg *config.Config) *AppLog {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	writers := make([]io.Writer, 0, 2)
	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	if !cfg.App.IsTesting() {
		writers = append(writers, &lumberjack.Logger{
			Filename: filepath.Join(
				filepath.Join(app.GetRootPath(), cfg.Log.Dir),
				fmt.Sprintf("request-%s.log", time.Now().Format("2006-01-02")),
			),
			MaxSize:  3000, // 3GB
			MaxAge:   cfg.Log.MaxAge,
			Compress: true,
		})
	}
	ml := zerolog.MultiLevelWriter(writers...)
	level := getLevel(cfg.Log.Level)
	zerolog.SetGlobalLevel(level)

	z := zerolog.New(ml).With().Timestamp().Logger()

	return &AppLog{
		lgr: z,
	}
}

func getLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error", "err":
		return zerolog.ErrorLevel
	case "fatal", "crit", "critical":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *AppLog) Debug(msg string, args ...any) {
	printLog(l.lgr.Debug(), msg, args)
}

func (l *AppLog) Info(msg string, args ...any) {
	printLog(l.lgr.Info(), msg, args)
}

func (l *AppLog) Warn(msg string, args ...any) {
	printLog(l.lgr.Warn(), msg, args)
}

func (l *AppLog) Error(msg string, args ...any) {
	printLog(l.lgr.Error(), msg, args)
}

func (l *AppLog) Fatal(msg string, args ...any) {
	printLog(l.lgr.Fatal(), msg, args)
}

func (l *AppLog) Panic(msg string, args ...any) {
	printLog(l.lgr.Panic(), msg, args)
}

func (l *AppLog) With(args ...any) Logger {
	return &AppLog{lgr: l.lgr.With().Fields(args).Logger()}
}

func printLog(lgr *zerolog.Event, msg string, fields []any) {
	if len(fields) > 0 {
		if err, ok := fields[0].(error); ok {
			lgr = lgr.Err(err)
			fields = fields[1:]
		}
	}

	args := []any{}
	count := strings.Count(msg, "%")
	if count > 0 {
		args = fields[:count]
		fields = fields[count:]
	}

	lgr.Fields(fields).Msgf(msg, args...)
}
