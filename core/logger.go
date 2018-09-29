package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggingLevel はloggingのsuppressするlevelを表します.
// 指定されたlevel未満のlogは出力されません.
type LoggingLevel int

const (
	// LoggingLvlDebug -
	LoggingLvlDebug LoggingLevel = iota - 1
	// LoggingLvlInfo -
	LoggingLvlInfo
	// LoggingLvlWarn -
	LoggingLvlWarn
	// LoggingLvlError -
	LoggingLvlError
)

func toLoggingLevel(s string) LoggingLevel {
	var lvl LoggingLevel
	switch s = strings.ToLower(s); s {
	case "debug":
		lvl = LoggingLvlDebug
	case "info":
		lvl = LoggingLvlInfo
	case "warn":
		lvl = LoggingLvlWarn
	case "error":
		lvl = LoggingLvlError
	default:
		lvl = LoggingLvlInfo
	}
	return lvl
}

// LoggingEncode -
type LoggingEncode string

const (
	// LoggingEncJSON -
	LoggingEncJSON LoggingEncode = "json"
	// LoggingEncText -
	LoggingEncText LoggingEncode = "text"
	// LoggingEncColorText -
	LoggingEncColorText LoggingEncode = "color"
)

func toLoggingEncode(s string) LoggingEncode {
	var enc LoggingEncode
	switch s = strings.ToLower(s); s {
	case "json":
		enc = LoggingEncJSON
	case "text", "console":
		enc = LoggingEncText
	case "color":
		enc = LoggingEncColorText
	default:
		enc = LoggingEncText
	}
	return enc
}

// LoggerConfig store logging configurations
type LoggerConfig struct {
	LoggingLevel  LoggingLevel
	LoggingEncode LoggingEncode
	Out           io.Writer
	NoTimestamp   bool

	level      zap.AtomicLevel
	encode     string
	encoderCfg zapcore.EncoderConfig
}

// NewLogger -
func NewLogger(level, encode string) (*zap.Logger, error) {
	return newLogger(&LoggerConfig{
		LoggingLevel:  toLoggingLevel(level),
		LoggingEncode: toLoggingEncode(encode),
		Out:           os.Stdout,
	})
}

func newLogger(c *LoggerConfig) (*zap.Logger, error) {
	return c.build()
}

func (c *LoggerConfig) build() (*zap.Logger, error) {
	switch c.LoggingLevel {
	case LoggingLvlDebug:
		c.level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case LoggingLvlInfo:
		c.level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case LoggingLvlWarn:
		c.level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case LoggingLvlError:
		c.level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		c.level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	c.encoderCfg = zapcore.EncoderConfig{
		TimeKey:        "timelocal",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	const (
		encodeJSON    = "json"
		encodeConsole = "text"
	)
	switch c.LoggingEncode {
	case LoggingEncJSON:
		c.encode = encodeJSON
	case LoggingEncText:
		c.encode = encodeConsole
	case LoggingEncColorText:
		c.encode = encodeConsole
		c.encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		c.encode = encodeConsole
	}

	// for testability
	if c.NoTimestamp {
		c.encoderCfg.TimeKey = ""
	}

	if c.Out == nil {
		return nil, errors.New("nil logger output")
	}

	encoder, err := c.buildEncoder()
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(c.Out), c.level)
	return zap.New(core, c.options()...), nil
}

func (c *LoggerConfig) buildEncoder() (zapcore.Encoder, error) {
	switch c.LoggingEncode {
	case LoggingEncJSON:
		return zapcore.NewJSONEncoder(c.encoderCfg), nil
	case LoggingEncText, LoggingEncColorText:
		return zapcore.NewConsoleEncoder(c.encoderCfg), nil
	}
	return nil, fmt.Errorf("unexpected logging encode %v", c.LoggingEncode)
}

func (c *LoggerConfig) options() []zap.Option {
	options := []zap.Option{
		zap.AddCallerSkip(0),
	}
	if c.LoggingLevel == LoggingLvlDebug {
		options = append(options, zap.AddCaller())
	}
	return options
}
