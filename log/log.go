package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gxxgle/go-utils/path"

	"github.com/sirupsen/logrus"
)

type (
	F = logrus.Fields
)

var (
	logfile *os.File
	Logger  = logrus.StandardLogger()
	L       = logrus.WithField("@pid", os.Getpid())
)

func init() {
	JSONFormat()
}

func TextFormat() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02T15:04:05.000",
		CallerPrettyfier: callerPrettyfier,
	})
}

func JSONFormat() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@time",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyFunc:  "@func",
			logrus.FieldKeyFile:  "@file",
			logrus.FieldKeyMsg:   "msg",
		},
		CallerPrettyfier: callerPrettyfier,
	})
}

func callerPrettyfier(rf *runtime.Frame) (string, string) {
	file := fmt.Sprintf(":%d", rf.Line)
	files := strings.Split(rf.File, "/")
	if len(files) > 0 {
		file = files[len(files)-1] + file
	}
	if len(files) > 1 {
		file = files[len(files)-2] + "/" + file
	}
	function := ""
	// functions := strings.Split(rf.Function, "/")
	// if len(functions) > 0 {
	// 	function = functions[len(functions)-1] + function
	// }
	// if len(functions) > 1 {
	// 	function = functions[len(functions)-2] + "/" + function
	// }
	return function, file
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
		L.WithField("path", logpath).WithError(err).Error("go-utils log mkdir failed")
		return
	}

	logfile, err = os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		L.WithField("path", logpath).WithError(err).Error("go-utils log open file failed")
		return
	}

	logrus.SetOutput(logfile)
}

func Debug() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
}

func LogIfError(err error, msg ...interface{}) {
	if err == nil {
		return
	}

	L.WithError(err).Error(msg...)
}

func FatalIfError(err error, msg ...interface{}) {
	if err == nil {
		return
	}

	L.WithError(err).Fatal(msg...)
}

func Close() {
	if logfile != nil {
		logfile.Close()
		logfile = nil
	}
}
