package util_logger

import (
	"io"
	defaultLog "log"
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

func (l Log) Trace(pairs ...any) {
	event := l.Logger.Trace()
	l.publish(event, lvlTrace, pairs...)
}

func (l Log) Debug(pairs ...any) {
	event := l.Logger.Debug()
	l.publish(event, lvlDebug, pairs...)
}

func (l Log) Info(pairs ...any) {
	event := l.Logger.Info()
	l.publish(event, lvlInfo, pairs...)
}

func (l Log) Warn(pairs ...any) {
	event := l.Logger.Warn()
	l.publish(event, lvlWarn, pairs...)
}

func (l Log) Error(err error, pairs ...any) {
	event := l.Logger.Error().Err(err)
	l.publish(event, lvlError, pairs...)
}

func (l Log) Fatal(err error, pairs ...any) {
	event := l.Logger.Fatal().Err(err)
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
		key, ok := pairs[i].(string)
		if !ok {
			continue
		}
		e = e.Interface(key, pairs[i+1])
	}
	if len(pairs)%2 != 0 {
		l.Warn("msg", "Logger received an odd number of arguments. Last argument ignored.")
	}
	return e
}
