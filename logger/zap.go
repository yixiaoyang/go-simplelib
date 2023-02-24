package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	sugaredLogger *zap.SugaredLogger
	fileLevel     *zap.AtomicLevel
	consoleLevel  *zap.AtomicLevel
	prefix        string
	name          string
}

var (
	zaps     sync.Map
	zapsLock sync.Mutex
)

const (
	FieldDelimiter  = " | "
	CustomDelimiter = " "
)

func zyCapitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}

func getEncoder(config Configuration) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.TimeKey = "time"
	encoderConfig.CallerKey = "line"
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.ConsoleSeparator = config.ConsoleSeparator
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	if config.ConsoleJSONFormat {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(l string) zapcore.Level {
	ll := strings.ToLower(l)
	switch ll {
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
func cloneZapLogger(z *ZapLogger) *ZapLogger {
	log := &ZapLogger{
		sugaredLogger: z.sugaredLogger,
		fileLevel:     z.fileLevel,
		consoleLevel:  z.consoleLevel,
		prefix:        z.prefix,
		name:          z.name,
	}
	return log
}

func newZapLogger(config Configuration, newCore bool) (*ZapLogger, error) {
	zapsLock.Lock()
	defer zapsLock.Unlock()

	if !newCore {
		val, ok := zaps.Load(config.Name)
		if ok {
			zapLogger, ok := val.(*ZapLogger)
			if ok {
				zlog := &ZapLogger{
					sugaredLogger: zapLogger.sugaredLogger,
					fileLevel:     zapLogger.fileLevel,
					consoleLevel:  zapLogger.consoleLevel,
					prefix:        config.Prefix,
					name:          config.Name,
				}
				return zlog, nil
			}
		}
	}

	var cores []zapcore.Core

	consoleLevelAtom := zap.NewAtomicLevelAt(getZapLevel(config.ConsoleLevel))
	fileLevelAtom := zap.NewAtomicLevelAt(getZapLevel(config.FileLevel))

	if config.EnableConsole {
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(config), writer, consoleLevelAtom)
		cores = append(cores, core)
	}

	if config.EnableFile {
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.FileMaxSize,
			Compress:   config.Compress,
			MaxAge:     config.FileMaxAge,
			MaxBackups: config.FileMaxNum,
		})
		core := zapcore.NewCore(getEncoder(config), writer, fileLevelAtom)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	logger := zap.New(combinedCore,
		// 需要开启此选项跳过函数的封装以获取期望的文件名和行号
		zap.AddCallerSkip(1),
		zap.AddCaller(),
	).Sugar()

	log := &ZapLogger{
		sugaredLogger: logger,
		fileLevel:     &fileLevelAtom,
		consoleLevel:  &consoleLevelAtom,
		prefix:        config.Prefix,
		name:          config.Name,
	}

	zaps.Store(config.Name, log)
	return log, nil
}

func (l *ZapLogger) SetLevel(level string) {
	l.consoleLevel.SetLevel(getZapLevel(level))
	l.fileLevel.SetLevel(getZapLevel(level))
}

func (l *ZapLogger) SetConsoleLevel(level string) {
	l.consoleLevel.SetLevel(getZapLevel(level))
}

func (l *ZapLogger) SetFileLevel(level string) {
	l.fileLevel.SetLevel(getZapLevel(level))
}

func (l *ZapLogger) GetLevel() string {
	if l.fileLevel.Level() > l.consoleLevel.Level() {
		return l.consoleLevel.Level().String()
	} else {
		return l.fileLevel.String()
	}
}

func (l *ZapLogger) IsLevelEnabled(level string) bool {
	if l.fileLevel.Level() > l.consoleLevel.Level() {
		return l.consoleLevel.Enabled(getZapLevel(level))
	} else {
		return l.fileLevel.Enabled(getZapLevel(level))
	}
}

func (l *ZapLogger) getPrefix() string {
	return l.getErrnoPrefix(0)
}

func (l *ZapLogger) levelLargeThan(level zapcore.Level) bool {
	return (l.fileLevel.Level() > level) && (l.consoleLevel.Level() > level)
}

func sprintf(args ...interface{}) string {
	str := ""
	for _, arg := range args {
		str += fmt.Sprintf("%v ", arg)
	}
	return str
}

func (l *ZapLogger) getErrnoPrefix(errno int) string {
	if l.prefix != "" {
		return strconv.Itoa(os.Getpid()) + FieldDelimiter + strconv.Itoa(errno) + FieldDelimiter + l.prefix + CustomDelimiter
	}
	return strconv.Itoa(os.Getpid()) + FieldDelimiter + strconv.Itoa(errno) + FieldDelimiter
}

func (l *ZapLogger) Debug(args ...interface{}) {
	if l.levelLargeThan(zapcore.DebugLevel) {
		return
	}
	l.sugaredLogger.Debug(l.getPrefix() + sprintf(args...))
}

func (l *ZapLogger) Info(args ...interface{}) {
	if l.levelLargeThan(zapcore.InfoLevel) {
		return
	}
	l.sugaredLogger.Info(l.getPrefix() + sprintf(args...))
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(l.getPrefix() + sprintf(args...))
}

func (l *ZapLogger) Error(errno int, args ...interface{}) {
	l.sugaredLogger.Error(l.getErrnoPrefix(errno) + sprintf(args...))
}

func (l *ZapLogger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(l.getPrefix() + sprintf(args...))
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(l.getPrefix() + sprintf(args...))
}

func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	if l.levelLargeThan(zapcore.DebugLevel) {
		return
	}
	l.sugaredLogger.Debugf(l.getPrefix()+template, args...)
}

func (l *ZapLogger) Infof(template string, args ...interface{}) {
	if l.levelLargeThan(zapcore.InfoLevel) {
		return
	}
	l.sugaredLogger.Infof(l.getPrefix()+template, args...)
}

func (l *ZapLogger) Warnf(template string, args ...interface{}) {
	l.sugaredLogger.Warnf(l.getPrefix()+template, args...)
}

func (l *ZapLogger) Errorf(errno int, template string, args ...interface{}) {
	l.sugaredLogger.Errorf(l.getErrnoPrefix(errno) + fmt.Sprintf(template, args...))
}

func (l *ZapLogger) Panicf(template string, args ...interface{}) {
	l.sugaredLogger.Panicf(l.getPrefix()+template, args...)
}

func (l *ZapLogger) Fatalf(template string, args ...interface{}) {
	l.sugaredLogger.Fatalf(l.getPrefix()+template, args...)
}

func (l *ZapLogger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *ZapLogger) GetPrefix() string {
	return l.prefix
}

func (l *ZapLogger) GetName() string {
	return l.name
}
