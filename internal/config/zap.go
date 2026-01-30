package config

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(v *viper.Viper) *zap.Logger {
	levelInt := v.GetInt("log.level")

	var level zapcore.Level
	switch levelInt {
	case 0:
		level = zapcore.PanicLevel
	case 1:
		level = zapcore.FatalLevel
	case 2:
		level = zapcore.ErrorLevel
	case 3:
		level = zapcore.WarnLevel
	case 4:
		level = zapcore.InfoLevel
	case 5:
		level = zapcore.DebugLevel
	default:
		level = zapcore.InfoLevel
	}

	encCfg := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	if v.GetString("app.env") == "dev" {
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encCfg.ConsoleSeparator = "  "
		encCfg.EncodeCaller = shortCallerFixed(30)
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encCfg),
			zapcore.AddSync(os.Stdout),
			level,
		)

		return zap.New(core)
	}

	encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)

	return zap.New(core)
}

// ===== helper padding caller (dev only) =====
func shortCallerFixed(width int) zapcore.CallerEncoder {
	return func(c zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		s := c.TrimmedPath()
		if len(s) >= width {
			enc.AppendString(s[:width])
			return
		}
		for len(s) < width {
			s += " "
		}
		enc.AppendString(s)
	}
}
