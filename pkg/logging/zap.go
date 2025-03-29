package logging

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerImpl struct {
	params    LogParams
	logger    *zap.Logger
	formatter *customJSONFormatter
}

func newLogger(params LogParams) *loggerImpl {
	if params.AddAttributesFromContext == nil {
		params.AddAttributesFromContext = func(ctx context.Context) []LogField {
			return nil
		}
	}

	formatter := newCustomJSONFormatter(params)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    formatter.formatLevel,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	level := zap.InfoLevel
	if params.DebugLevel {
		level = zap.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&colorWriter{os.Stdout, formatter}),
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	return &loggerImpl{
		params:    params,
		logger:    logger,
		formatter: formatter,
	}
}

type colorWriter struct {
	w         io.Writer
	formatter *customJSONFormatter
}

func (cw *colorWriter) Write(p []byte) (n int, err error) {
	colored := cw.formatter.Format(string(p))
	return cw.w.Write([]byte(colored))
}

func (l *loggerImpl) log(ctx context.Context, level zapcore.Level, msg string, fields ...LogField) {
	ctxFields := l.params.AddAttributesFromContext(ctx)
	allFields := make([]zap.Field, 0, len(fields)+len(ctxFields)+len(l.formatter.getDefaultFields()))

	for _, field := range fields {
		allFields = append(allFields, field.ToZapField())
	}
	for _, field := range ctxFields {
		allFields = append(allFields, field.ToZapField())
	}

	allFields = append(allFields, l.formatter.getDefaultFields()...)

	l.logger.Log(level, msg, allFields...)
}

func (l *loggerImpl) Info(ctx context.Context, msg string) {
	l.log(ctx, zapcore.InfoLevel, msg)
}

func (l *loggerImpl) Infof(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zapcore.InfoLevel, fmt.Sprintf(msg, args...))
}

func (l *loggerImpl) Infow(ctx context.Context, msg string, fields ...LogField) {
	l.log(ctx, zapcore.InfoLevel, msg, fields...)
}

func (l *loggerImpl) Debug(ctx context.Context, msg string) {
	l.log(ctx, zapcore.DebugLevel, msg)
}

func (l *loggerImpl) Debugf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zapcore.DebugLevel, fmt.Sprintf(msg, args...))
}

func (l *loggerImpl) Debugw(ctx context.Context, msg string, fields ...LogField) {
	l.log(ctx, zapcore.DebugLevel, msg, fields...)
}

func (l *loggerImpl) Warn(ctx context.Context, msg string) {
	l.log(ctx, zapcore.WarnLevel, msg)
}

func (l *loggerImpl) Warnf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zapcore.WarnLevel, fmt.Sprintf(msg, args...))
}

func (l *loggerImpl) Warnw(ctx context.Context, msg string, fields ...LogField) {
	l.log(ctx, zapcore.WarnLevel, msg, fields...)
}

func (l *loggerImpl) Error(ctx context.Context, msg string) {
	l.log(ctx, zapcore.ErrorLevel, msg)
}

func (l *loggerImpl) Errorf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zapcore.ErrorLevel, fmt.Sprintf(msg, args...))
}

func (l *loggerImpl) Errorw(ctx context.Context, msg string, fields ...LogField) {
	l.log(ctx, zapcore.ErrorLevel, msg, fields...)
}

func (l *loggerImpl) Fatal(ctx context.Context, msg string) {
	l.log(ctx, zapcore.FatalLevel, msg)
	os.Exit(1)
}

func (l *loggerImpl) Fatalf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zapcore.FatalLevel, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func (l *loggerImpl) Fatalw(ctx context.Context, msg string, fields ...LogField) {
	l.log(ctx, zapcore.FatalLevel, msg, fields...)
	os.Exit(1)
}

const (
	LevelCritical = zapcore.Level(60) // high number to avoid conflicts with other levels
)

func (l *loggerImpl) Critical(ctx context.Context, msg string) {
	l.log(ctx, LevelCritical, msg)
}

func (l *loggerImpl) Criticalf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelCritical, fmt.Sprintf(msg, args...))
}

func (l *loggerImpl) Criticalw(ctx context.Context, msg string, fields ...LogField) {
	l.log(ctx, LevelCritical, msg, fields...)
}

func (l *loggerImpl) Print(args ...any) {
	l.log(context.Background(), zapcore.InfoLevel, fmt.Sprint(args...))
}

func (l *loggerImpl) Printf(msg string, args ...any) {
	l.log(context.Background(), zapcore.InfoLevel, fmt.Sprintf(msg, args...))
}
