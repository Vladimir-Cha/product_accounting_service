package logger

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) log(ctx context.Context, level, msg string, fields ...interface{}) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s", level, msg))

	for i := 0; i < len(fields); i += 2 {
		if key, ok := fields[i].(string); ok {
			sb.WriteString(fmt.Sprintf(" %s=%v", key, fields[i+1]))
		}
	}

	log.Println(sb.String())
}

func (l *ConsoleLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.log(ctx, "DEBUG", msg, fields...)
}

func (l *ConsoleLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.log(ctx, "INFO", msg, fields...)
}

func (l *ConsoleLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.log(ctx, "WARN", msg, fields...)
}

func (l *ConsoleLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.log(ctx, "ERROR", msg, fields...)
}

func (l *ConsoleLogger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	l.log(ctx, "FATAL", msg, fields...)
	log.Fatal()
}
