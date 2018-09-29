package core

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		desc   string
		config *LoggerConfig
		msg    string
		want   string
	}{
		{
			desc: "json logging",
			config: &LoggerConfig{
				LoggingLevel:  LoggingLvlInfo,
				LoggingEncode: LoggingEncJSON,
			},
			msg:  "golang or gohome",
			want: "{\"level\":\"INFO\",\"msg\":\"golang or gohome\"}\n",
		},
		{
			desc: "text logging",
			config: &LoggerConfig{
				LoggingLevel:  LoggingLvlInfo,
				LoggingEncode: LoggingEncText,
			},
			msg:  "golang or gohome",
			want: "INFO\tgolang or gohome\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			b := &bytes.Buffer{}
			tt.config.Out = b
			tt.config.NoTimestamp = true
			logger, err := newLogger(tt.config)
			if err != nil {
				t.Fatalf("failed to new zap %s", err)
			}

			logger.Info(tt.msg)
			if err := logger.Sync(); err != nil {
				t.Fatalf("Logger.Sync() return err: %s", err)
			}

			if diff := cmp.Diff(b.String(), tt.want); diff != "" {
				t.Errorf("(-got +want)\n%s", diff)
			}
		})
	}
}
