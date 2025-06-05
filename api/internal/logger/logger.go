package logger

import (
	"context"
	"io"
	defaultLog "log"
	"maps"
	"net/http"
	"sync"
	"via/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Log struct {
	zerolog.Logger
	iconEnabled bool
}

var (
	instance *Log
	mutex    sync.RWMutex
	lvlTrace = "trace"
	lvlDebug = "debug"
	lvlInfo  = "info"
	lvlWarn  = "warn"
	lvlError = "error"
	lvlFatal = "fatal"
	incons   = map[string]string{
		lvlTrace: "🔍",
		lvlDebug: "🐛",
		lvlInfo:  "ℹ️",
		lvlWarn:  "⚠️",
		lvlError: "❌",
		lvlFatal: "💀"}
)

func Init(logCfg config.Log) {
	mutex.Lock()
	defer mutex.Unlock()

	var writers []io.Writer
	if logCfg.DefaultWriter.Enabled {
		writers = append(writers, DefaultWriter{logCfg.DefaultWriter}.Writer())
	}
	if logCfg.FileWriter.Enabled {
		writers = append(writers, FileWriter{logCfg.FileWriter}.Writer())
	}
	if logCfg.ConsoleWriter.Enabled {
		writers = append(writers, ConsoleWriter{}.Writer())
	}
	logLevel, err := zerolog.ParseLevel(logCfg.Level)
	if err != nil {
		defaultLog.Fatalf("unkonwn log level: %s, err:%v", logCfg.Level, err)
	}
	zerolog.SetGlobalLevel(logLevel)
	log.Logger = zerolog.New(io.MultiWriter(writers...)).
		With().
		Timestamp().
		CallerWithSkipFrameCount(4).
		Logger()
	instance = &Log{Logger: log.Logger, iconEnabled: logCfg.IconEnabled}
}

func Get() *Log {
	mutex.RLock()
	if instance != nil {
		mutex.RUnlock()
		return instance
	}
	mutex.RUnlock()
	Init(config.Get().Log)
	return instance
}

func (l Log) Trace(ctx context.Context, pairs ...any) {
	event := l.Logger.Trace()
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlTrace, pairs...)
}

func (l Log) Debug(ctx context.Context, pairs ...any) {
	event := l.Logger.Debug()
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlDebug, pairs...)
}

func (l Log) Info(ctx context.Context, pairs ...any) {
	event := l.Logger.Info()
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlInfo, pairs...)
}

func (l Log) Warn(ctx context.Context, pairs ...any) {
	event := l.Logger.Warn()
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlWarn, pairs...)
}

func (l Log) Error(ctx context.Context, err error, pairs ...any) {
	event := l.Logger.Error().Err(err)
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlError, pairs...)
}

func (l Log) Fatal(ctx context.Context, err error, pairs ...any) {
	event := l.Logger.Fatal().Err(err)
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlFatal, pairs...)
}

func (l Log) addIcon(e *zerolog.Event, lvl string) *zerolog.Event {
	icon, ok := incons[lvl]
	if l.iconEnabled && ok {
		e.Str("lvl", icon)
	}
	return e
}

func (l Log) publish(e *zerolog.Event, level string, pairs ...any) {
	e = l.addIcon(e, level)
	e = l.addFields(e, pairs...)
	e.Msg("")
}

// arg conversion to key value
func (l Log) addFields(e *zerolog.Event, pairs ...any) *zerolog.Event {
	for i := 0; i+1 < len(pairs); i += 2 {
		//if key is not a string, skip it
		key, ok := pairs[i].(string)
		if !ok {
			continue
		}
		e = e.Interface(key, pairs[i+1])
	}
	if len(pairs)%2 != 0 {
		l.Warn(context.Background(), "msg", "Logger received an odd number of arguments. Last argument ignored.", "key", pairs[len(pairs)-1])
		// If the last argument is not a key-value pair, we ignore it value
	}
	return e
}

type ctxKey string

const logContextKey ctxKey = "logFields"

type LogFields map[string]any

func (l Log) WithLogFieldsInRequest(r *http.Request, pairs ...any) *http.Request {
	return r.WithContext(l.WithLogFields(r.Context(), pairs...))
}

// WithLogFields returns a new context with logging fields added
func (l Log) WithLogFields(ctx context.Context, pairs ...any) context.Context {
	ctxFields, _ := ctx.Value(logContextKey).(LogFields)
	if ctxFields == nil {
		ctxFields = make(LogFields)
	}
	maps.Copy(ctxFields, l.pairsToLogFields(pairs...))
	return context.WithValue(ctx, logContextKey, ctxFields)
}

func (l Log) pairsToLogFields(pairs ...any) LogFields {
	fields := make(LogFields)
	for i := 0; i+1 < len(pairs); i += 2 {
		//if key is not a string, skip it
		key, ok := pairs[i].(string)
		if !ok {
			continue
		}
		fields[key] = pairs[i+1]
	}
	if len(pairs)%2 != 0 {
		if len(pairs)%2 != 0 {
			l.Warn(context.Background(), "msg", "Logger received an odd number of arguments. Last argument ignored.", "key", pairs[len(pairs)-1])
			// If the last argument is not a key-value pair, we ignore it value
		}
	}
	return fields
}

func (l Log) addContextFields(ctx context.Context, pairs ...any) []any {
	if ctx == nil {
		return pairs
	}
	fields, ok := ctx.Value(logContextKey).(LogFields)
	if !ok || len(fields) == 0 {
		return pairs
	}
	maps.Copy(fields, l.pairsToLogFields(pairs...))
	result := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		result = append(result, k, v)
	}
	return result
}
