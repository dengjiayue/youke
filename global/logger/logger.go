package logger

import (
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/natefinch/lumberjack"
)

type LoggerConfig struct {
	OutPath string `yaml:"OutPath"`
}

const (
	ErrorLevel   = "ERROR"
	WarningLevel = "WARNING"
	InfoLevel    = "INFO"
)

type Logger struct {
	warningLog *log.Logger
	errorLog   *log.Logger
}

func NewLogger(logConfig *LoggerConfig) (*Logger, error) {
	// 获取当前时间用于日志文件名
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(logConfig.OutPath, currentTime+".log")

	// 创建lumberjack日志轮转器
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10,   // 最大文件大小，单位MB
		MaxBackups: 5,    // 保留旧文件的最大个数
		MaxAge:     30,   // 保留旧文件的最大天数
		Compress:   true, // 是否压缩旧文件
	}

	return &Logger{
		warningLog: log.New(lumberjackLogger, WarningLevel+" ", log.Ldate|log.Ltime),
		errorLog:   log.New(lumberjackLogger, ErrorLevel+" ", log.Ldate|log.Ltime),
	}, nil
}

func (l *Logger) Info(msg ...any) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("%s %s:%d: %v\n", InfoLevel, filepath.Base(file), line, msg)
	}
}

func (l *Logger) Warning(msg ...any) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("%s %s:%d: %v\n", WarningLevel, filepath.Base(file), line, msg)
		l.warningLog.Printf("%s:%d: %v\n", filepath.Base(file), line, msg) // 使用 %v 以支持多参数
	}
}

func (l *Logger) Error(msg ...any) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("%s %s:%d: %v\n", ErrorLevel, filepath.Base(file), line, msg)
		l.errorLog.Printf("%s:%d: %v\n", filepath.Base(file), line, msg) // 使用 %v 以支持多参数
	}
}
