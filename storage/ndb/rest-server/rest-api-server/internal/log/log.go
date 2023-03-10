package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"hopsworks.ai/rdrs/internal/config"
)

func SetupLogger(conf config.AllConfigs) (*zap.Logger, error) {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.InfoLevel)
	return zap.Config{
		Level:            atom,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			CallerKey:    "caller",
			TimeKey:      "time",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeTime:   zapcore.RFC3339NanoTimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()
}
