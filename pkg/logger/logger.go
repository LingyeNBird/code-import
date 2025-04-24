package logger

import (
	"ccrctl/pkg/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

var (
	Logger                *zap.SugaredLogger
	SuccessfulLogFilePath string
	SuccessfulLogFile     *os.File
)

const (
	SuccessfulLog = "successful.log"
	InfoLog       = "migrate.log"
)

func init() {
	logFile, err := os.OpenFile(InfoLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Error("打开日志文件失败:", err)
	}
	//defer logFile.Close()

	SuccessfulLogFile, err = os.OpenFile(SuccessfulLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Error("打开日志文件失败:", err)
	}

	//创建一个多写入器，同时写入标准输出和日志文件
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(logFile),
	)

	// 创建一个JSON编码器配置
	//encoderConfig := zapcore.EncoderConfig{
	//	TimeKey:        "time",
	//	LevelKey:       "level",
	//	NameKey:        "logger",
	//	CallerKey:      "caller",
	//	MessageKey:     "msg",
	//	StacktraceKey:  "stacktrace",
	//	LineEnding:     zapcore.DefaultLineEnding,
	//	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	//	EncodeTime:     zapcore.ISO8601TimeEncoder,
	//	EncodeDuration: zapcore.StringDurationEncoder,
	//	EncodeCaller:   zapcore.ShortCallerEncoder,
	//}

	level := zap.NewAtomicLevel()
	switch config.Cfg.GetString("migrate.log_level") {
	case "debug":
		level.SetLevel(zap.DebugLevel)
	case "info":
		level.SetLevel(zap.InfoLevel)
	case "warn":
		level.SetLevel(zap.WarnLevel)
	case "error":
		level.SetLevel(zap.ErrorLevel)
	default:
		fmt.Println("未知的日志级别:", config.Cfg.GetString("migrate.log_level"))
		return
	}

	// 创建一个核心，使用多写入器和JSON编码器配置
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		multiWriteSyncer,
		level,
	)

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		Logger.Error("Error:", err)
	}

	// 拼接日志文件的绝对路径
	SuccessfulLogFilePath = filepath.Join(currentDir, SuccessfulLog)

	logger := zap.New(core)
	logger = logger.WithOptions(zap.AddCaller())
	Logger = logger.Sugar()

}

// 记录迁移成功的仓库
func RecordSuccessfulRepo(repoPath string) {
	_, err := SuccessfulLogFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + " " + repoPath + "\n")
	if err != nil {
		Logger.Error("Error:", err)
	}
}
