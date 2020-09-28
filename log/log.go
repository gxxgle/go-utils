package log

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gxxgle/go-utils/flag"
	"github.com/gxxgle/go-utils/path"

	"github.com/phuslu/log"
	"github.com/urfave/cli/v2"
)

var (
	fileWriter *log.FileWriter
)

func init() {
	log.DefaultLogger.SetLevel(log.InfoLevel)
	log.DefaultLogger.Caller = 1
}

func InitFromFlag() {
	flag.Add(&cli.StringFlag{
		Name:    "log_level",
		Usage:   "set log level. (debug, info, error)",
		EnvVars: []string{"LOG_LEVEL"},
		Value:   "info",
	}, func(ctx *cli.Context) error {
		lvl := log.ParseLevel(ctx.String("log_level"))
		if lvl >= log.TraceLevel && lvl <= log.PanicLevel {
			log.DefaultLogger.SetLevel(lvl)
		}
		return nil
	})

	flag.Add(&cli.StringFlag{
		Name:    "log_type",
		Usage:   "set log type. (json, console, color_console)",
		EnvVars: []string{"LOG_TYPE"},
		Value:   "json",
	}, func(ctx *cli.Context) error {
		switch ctx.String("log_type") {
		case "console":
			Console()
		case "color_console":
			ColorConsole()
		}
		return nil
	})
}

func Console() {
	log.DefaultLogger.Writer = &log.ConsoleWriter{
		ColorOutput:    false,
		QuoteString:    false,
		EndWithMessage: false,
	}
}

func ColorConsole() {
	log.DefaultLogger.Writer = &log.ConsoleWriter{
		ColorOutput:    true,
		QuoteString:    false,
		EndWithMessage: false,
	}
}

func File(logpaths ...string) {
	logpath := ""
	if len(logpaths) > 0 {
		logpath = logpaths[0]
	}

	if len(logpath) == 0 {
		dir := path.TopLevelDir(path.CurrentDir())
		filename := path.CurrentFilename()
		logpath = fmt.Sprintf("%s/log/%s.log", dir, filename)
	}

	err := os.MkdirAll(filepath.Dir(logpath), os.ModePerm)
	if err != nil {
		log.Error().Err(err).Str("path", logpath).Msg("log mkdir failed")
		return
	}

	fileWriter = &log.FileWriter{
		Filename:  logpath,
		FileMode:  0600,
		LocalTime: true,
	}
	log.DefaultLogger.Writer = fileWriter
}

func SetLevel(lvl log.Level) {
	log.DefaultLogger.SetLevel(lvl)
}

func SetDebug() {
	SetLevel(log.DebugLevel)
}

func LogIfError(err error, msgs ...string) {
	if err == nil {
		return
	}

	msg := ""
	if len(msgs) > 0 {
		msg = msgs[0]
	}

	log.Error().Caller(log.DefaultLogger.Caller + 1).Err(err).Msg(msg)
}

func LogIfFuncError(fn func() error, msgs ...string) {
	LogIfError(fn(), msgs...)
}

func FatalIfError(err error, msgs ...string) {
	if err == nil {
		return
	}

	msg := ""
	if len(msgs) > 0 {
		msg = msgs[0]
	}

	log.Fatal().Caller(log.DefaultLogger.Caller + 1).Err(err).Msg(msg)
}

func Close() error {
	if fileWriter != nil {
		return fileWriter.Close()
	}

	return nil
}
