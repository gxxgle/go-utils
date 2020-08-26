package log

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gxxgle/go-utils/path"

	"github.com/phuslu/log"
)

var (
	fileWriter *log.FileWriter
)

func init() {
	log.DefaultLogger.Level = log.InfoLevel
	log.DefaultLogger.Caller = 1
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
		log.Error().Err(err).Str("path", logpath).Msg("go-utils log mkdir failed")
		return
	}

	fileWriter = &log.FileWriter{
		Filename:  logpath,
		FileMode:  0600,
		LocalTime: true,
	}
	log.DefaultLogger.Writer = fileWriter
}

func SetDebug() {
	log.DefaultLogger.Level = log.DebugLevel
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
