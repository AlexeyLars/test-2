package main

import "fmt"

// ExternalLogger - сторонняя библиотека логирования
type ExternalLogger struct{}

func (el *ExternalLogger) LogMessage(msg string) {
	fmt.Printf("External log: %s\n", msg)
}

// Logger - целевой интерфейс
type Logger interface {
	Log(message string)
}

// LoggerAdapter - адаптер для интеграции внешнего логгера
type LoggerAdapter struct {
	externalLogger *ExternalLogger
}

func NewLoggerAdapter(externalLogger *ExternalLogger) *LoggerAdapter {
	return &LoggerAdapter{
		externalLogger: externalLogger,
	}
}

func (adapter *LoggerAdapter) Log(message string) {
	adapter.externalLogger.LogMessage(message)
}

func main() {
	fmt.Println("=== Adapter Pattern ===")

	externalLogger := &ExternalLogger{}
	logger := NewLoggerAdapter(externalLogger)

	logger.Log("This is a test message.")
	fmt.Println()
}
