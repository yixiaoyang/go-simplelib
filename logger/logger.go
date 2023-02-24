package logger

import (
	"sync"
)

// Logger is a general logger interface.
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(errno int, args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(errno int, template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	// SetLevel 作用域：全局，会更改同一个name的logger的日志等级，同时设置console和file输出的level
	// 可选输出等级："debug", "info", "warn","error","fatal"
	SetLevel(level string)

	// SetLevel 作用域：全局，会更改同一个name的logger的日志等级，仅设置console输出的level
	// 可选输出等级："debug", "info", "warn","error","fatal"
	SetConsoleLevel(level string)

	// SetLevel 作用域：全局，会更改同一个name的logger的日志等级，仅设置file输出的level
	// 可选输出等级："debug", "info", "warn","error","fatal"
	SetFileLevel(level string)

	// 返回当前level值。
	// 当file和console的level设置不一致时，返回level较低的值（debug < info < warn < error）
	GetLevel() string

	// 检查level是否使能
	// 当file和console的level设置不一致时，对比level较低的值（debug < info < warn < error）
	IsLevelEnabled(level string) bool

	// SetPrefix 作用域：当前实例，仅设置当前实例的前缀
	SetPrefix(prefix string)
	GetPrefix() string

	GetName() string
}

var (
	loggers     sync.Map
	loggersLock sync.Mutex
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// Configuration stores the config for the Logger
// For some loggers there can only be one level across writers,for such the level of Console is picked by default
type Configuration struct {
	// logger name
	Name string

	// Console print options
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string

	// EnableFile that whether to log information to a file
	EnableFile bool

	// FileLevel is the log need json format for file
	FileJSONFormat bool

	// FileLevel is the log level that for file
	FileLevel string

	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.
	Filename string

	// FileMaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	FileMaxSize int

	// FileMaxNum is the maximum number of file of log file before it gets rotated.
	// Default value is to 10
	FileMaxNum int

	// FileMaxAge is the maximum number of days to retain old log files based
	// on the timestamp encoded in their filename.  Note that a day is defined
	// as 24 hours and may not exactly correspond to calendar days due to
	// daylight savings, leap seconds, etc. The default is not to remove old
	// log files based on age.
	FileMaxAge int

	// Compress support, default false
	Compress bool

	ConsoleSeparator string

	Prefix string
}

// GetLogger 获取名为name的logger，如果不存在则新建默认logger。
// logger默认文件名为'name.log'
func GetLogger(name, prefix string) Logger {
	loggersLock.Lock()
	defer loggersLock.Unlock()
	val, ok := loggers.Load(name)
	if ok {
		zlog, ok := val.(*ZapLogger)
		if !ok {
			return nil
		}
		log := cloneZapLogger(zlog)
		log.SetPrefix(prefix)
		return log
	}

	config := Configuration{
		Name:              name,
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      "debug",
		EnableFile:        true,
		FileJSONFormat:    false,
		FileLevel:         "debug",
		FileMaxSize:       100,
		FileMaxNum:        10,
		Compress:          false,
		ConsoleSeparator:  " | ",
		Prefix:            prefix,
	}
	if config.Filename == "" {
		if name != "" {
			config.Filename = "./log/" + name + ".log"
		} else {
			config.Filename = "./log/" + "default.log"
		}
	}
	log, err := newZapLogger(config, false)
	if err != nil {
		log = &ZapLogger{}
	}
	loggers.Store(name, log)
	return log
}

// SetLogger 设置名为name，配置为config的logger，如果存在则重建
func SetLogger(name string, config Configuration) (Logger, error) {
	if !(config.EnableConsole || config.EnableFile) {
		panic("Either EnableConsole or EnableFile should be true")
	}
	if config.ConsoleSeparator == "" {
		config.ConsoleSeparator = " | "
	}

	config.Name = name
	if config.Filename == "" {
		if name != "" {
			config.Filename = "log/" + name + ".log"
		} else {
			config.Filename = "log/default.log"
		}
	}

	zLogger, err := newZapLogger(config, true)
	if err != nil {
		zLogger = &ZapLogger{}
	}

	loggersLock.Lock()
	defer loggersLock.Unlock()

	_, ok := loggers.Load(name)
	if ok {
		loggers.Delete(name)
	}
	loggers.Store(name, zLogger)
	return zLogger, nil
}

// WrappedLogger返回一个自定义前缀的新包装logger，底层仍使用同一个name的sugerLogger。
func WrappedLogger(logger Logger, prefix string) Logger {
	if logger == nil {
		return GetLogger("", prefix)
	}
	return GetLogger(logger.GetName(), prefix)
}
