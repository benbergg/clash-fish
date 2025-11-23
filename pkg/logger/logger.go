package logger

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init 初始化日志系统
func Init(logDir string, debug bool) error {
	// 创建日志目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 设置日志级别
	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	// 配置日志文件
	logFile := filepath.Join(logDir, "clash-fish.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// 多输出：控制台 + 文件
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, file)

	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	return nil
}

// Info 记录 Info 级别日志
func Info() *zerolog.Event {
	return log.Info()
}

// Debug 记录 Debug 级别日志
func Debug() *zerolog.Event {
	return log.Debug()
}

// Warn 记录 Warn 级别日志
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error 记录 Error 级别日志
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal 记录 Fatal 级别日志并退出
func Fatal() *zerolog.Event {
	return log.Fatal()
}
