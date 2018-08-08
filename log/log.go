package log

import (
	"fmt"
	"log"
	"os"

	"github.com/gxxgle/go-utils/path"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// default options
var (
	pid               = os.Getpid()
	defaultLogger     *Logger
	DefaultZapConf    = zap.NewProductionConfig()
	DefaultPath       = ""
	DefaultMaxSize    = 512 // 512 MB
	DefaultMaxAge     = 0
	DefaultMaxBackups = 0
	DefaultLocalTime  = true
	DefaultCompress   = true
)

// Logger provides fast, leveled, structured logging
type Logger struct {
	opts    []zap.Option
	zapConf *zap.Config
	lcConf  *lumberjack.Logger
	Log     *zap.SugaredLogger
}

// Option to logger
type Option func(*Logger)

func init() {
	var err error

	defaultLogger, err = NewLogger(ZapOption(zap.AddCallerSkip(1)))
	if err != nil {
		log.Fatalln(err)
	}
}

// Path is the file path to write logs.
// logs write to console if not set path
func Path(path string) Option {
	return func(o *Logger) {
		o.lcConf.Filename = path
	}
}

// MaxSize is the maximum size in MB of the log file before it gets rotated
// default is 512 MB
func MaxSize(size int) Option {
	return func(o *Logger) {
		o.lcConf.MaxSize = size
	}
}

// MaxAge is the maximum number of days to retain old log files based on the
// timestamp encoded in their filename.  Note that a day is defined as 24
// hours and may not exactly correspond to calendar days due to daylight
// savings, leap seconds, etc. The default is not to remove old log files
// based on age.
func MaxAge(age int) Option {
	return func(o *Logger) {
		o.lcConf.MaxAge = age
	}
}

// MaxBackups is the maximum number of old log files to retain.  The default
// is to retain all old log files (though MaxAge may still cause them to get
// deleted.)
func MaxBackups(backups int) Option {
	return func(o *Logger) {
		o.lcConf.MaxBackups = backups
	}
}

// LocalTime determines if the time used for formatting the timestamps in backup
// files is the computer's local time.
// default is enabled
func LocalTime(enabled bool) Option {
	return func(o *Logger) {
		o.lcConf.LocalTime = enabled
	}
}

// Compress determines if the rotated log files should be compressed using gzip
// default is enabled
func Compress(enabled bool) Option {
	return func(o *Logger) {
		o.lcConf.Compress = enabled
	}
}

// Development puts the logger in development mode, which changes the
// behavior of DPanicLevel and takes stacktraces more liberally.
func Development(enabled bool) Option {
	return func(o *Logger) {
		o.zapConf.Development = enabled
	}
}

// SetLevel set logger level
func SetLevel(level zapcore.Level) Option {
	return func(o *Logger) {
		o.zapConf.Level.SetLevel(level)
	}
}

// ZapOption for zap
func ZapOption(opt zap.Option) Option {
	return func(o *Logger) {
		o.opts = append(o.opts, opt)
	}
}

// NewLogger return a logger
func NewLogger(opts ...Option) (*Logger, error) {
	DefaultZapConf.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	DefaultZapConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	DefaultZapConf.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	out := &Logger{
		zapConf: &DefaultZapConf,
		lcConf: &lumberjack.Logger{
			Filename:   DefaultPath,
			MaxSize:    DefaultMaxSize,
			MaxAge:     DefaultMaxAge,
			MaxBackups: DefaultMaxBackups,
			LocalTime:  DefaultLocalTime,
			Compress:   DefaultCompress,
		},
	}

	for _, opt := range opts {
		opt(out)
	}

	out.useLumberjack()

	out.opts = append(out.opts, zap.AddCallerSkip(1))
	logger, err := out.zapConf.Build(out.opts...)
	if err != nil {
		return nil, err
	}

	out.Log = logger.Sugar().With("pid", pid)
	return out, nil
}

func (l *Logger) useLumberjack() {
	if l.lcConf.Filename == "" {
		return
	}

	l.opts = append(
		l.opts,
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(l.zapConf.EncoderConfig),
				zapcore.AddSync(l.lcConf),
				l.zapConf.Level,
			)
		}),
	)
}

func InitFileLogger() {
	dir := path.TopLevelDir(path.CurrentDir())
	filename := path.CurrentFilename()
	l, err := NewLogger(
		Path(dir+"/log/"+filename+".log"),
		ZapOption(zap.AddCallerSkip(1)),
	)
	if err != nil {
		log.Fatalln(err)
	}

	SetDefaultLogger(l)
}

func SetDefaultLogger(l *Logger) {
	defaultLogger = l
}

func (l *Logger) With(keysAndValues ...interface{}) {
	l.Log.With(keysAndValues...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	defer l.Log.Sync()
	l.Log.Debugw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	defer l.Log.Sync()
	l.Log.Infow(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	defer l.Log.Sync()
	l.Log.Warnw(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	defer l.Log.Sync()
	l.Log.Errorw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

func Errorln(a ...interface{}) {
	defaultLogger.Errorw(fmt.Sprint(a...))
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	defer l.Log.Sync()
	l.Log.DPanicw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	defer l.Log.Sync()
	l.Log.Panicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	defer l.Log.Sync()
	l.Log.Fatalw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

func Fatalln(a ...interface{}) {
	defaultLogger.Fatalw(fmt.Sprint(a...))
}
