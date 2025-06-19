package app_log

import (
	"context"
	"io"
	defaultLog "log"
	"maps"
	"net/http"
	"via/internal/log"

	"github.com/rs/zerolog"
	zero_log "github.com/rs/zerolog/log"
)

type AppLogger struct {
	zerolog.Logger
	iconEnabled bool
}

type LogCfg struct {
	FileWriter    FileWriterCfg    `envPrefix:"FILE_WRITER_"    json:"fileWriter"`
	ConsoleWriter ConsoleWriterCfg `envPrefix:"CONSOLE_WRITER_" json:"consoleWriter"`
	DefaultWriter DefaultWriterCfg `envPrefix:"DEFAULT_WRITER_" json:"defaultWriter"`
	IconEnabled   bool             `env:"ICON_ENABLED"          envDefault:"false"  json:"iconEnabled"`
	Level         string           `env:"LEVEL"                 envDefault:"debug" json:"level"`
}

var (
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

func New(logCfg LogCfg) log.Logger {
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
	zero_log.Logger = zerolog.New(io.MultiWriter(writers...)).
		With().
		Timestamp().
		CallerWithSkipFrameCount(4).
		Logger()
	return &AppLogger{Logger: zero_log.Logger, iconEnabled: logCfg.IconEnabled}
}

func (l AppLogger) Trace(ctx context.Context, pairs ...any) {
	event := l.Logger.Trace()
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlTrace, pairs...)
}

func (l AppLogger) Debug(ctx context.Context, pairs ...any) {
	event := l.Logger.Debug()
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlDebug, pairs...)
}

func (l AppLogger) Info(ctx context.Context, pairs ...any) {
	event := l.Logger.Info()
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlInfo, pairs...)
}

func (l AppLogger) Warn(ctx context.Context, pairs ...any) {
	event := l.Logger.Warn()
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlWarn, pairs...)
}

func (l AppLogger) Error(ctx context.Context, err error, pairs ...any) {
	event := l.Logger.Error().Err(err)
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlError, pairs...)
}

func (l AppLogger) Fatal(ctx context.Context, err error, pairs ...any) {
	event := l.Logger.Fatal().Err(err)
	pairs = l.addContextFields(ctx, pairs...)
	l.publish(event, lvlFatal, pairs...)
}

func (l AppLogger) addIcon(e *zerolog.Event, lvl string) *zerolog.Event {
	icon, ok := incons[lvl]
	if l.iconEnabled && ok {
		e.Str("lvl", icon)
	}
	return e
}

func (l AppLogger) publish(e *zerolog.Event, level string, pairs ...any) {
	e = l.addIcon(e, level)
	e = l.addFields(e, pairs...)
	e.Msg("")
}

// arg conversion to key value
func (l AppLogger) addFields(e *zerolog.Event, pairs ...any) *zerolog.Event {
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

func (l AppLogger) WithLogFieldsInRequest(r *http.Request, pairs ...any) *http.Request {
	return r.WithContext(l.WithLogFields(r.Context(), pairs...))
}

// WithLogFields returns a new context with logging fields added
func (l AppLogger) WithLogFields(ctx context.Context, pairs ...any) context.Context {
	if ctx == nil {
		return ctx
	}
	ctxFields, _ := ctx.Value(logContextKey).(LogFields)
	if ctxFields == nil {
		ctxFields = make(LogFields)
	}
	maps.Copy(ctxFields, l.pairsToLogFields(pairs...))
	return context.WithValue(ctx, logContextKey, ctxFields)
}

func (l AppLogger) pairsToLogFields(pairs ...any) LogFields {
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

func (l AppLogger) addContextFields(ctx context.Context, pairs ...any) []any {
	if ctx == nil {
		return pairs
	}
	fieldsInContext, ok := ctx.Value(logContextKey).(LogFields)
	if !ok || len(fieldsInContext) == 0 {
		return pairs
	}
	fields := make(LogFields)
	maps.Copy(fields, fieldsInContext)
	maps.Copy(fields, l.pairsToLogFields(pairs...))
	result := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		result = append(result, k, v)
	}
	return result
}
