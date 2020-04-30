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
	L       = logrus.WithField("@pid", os.Getpid())
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@time",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyFunc:  "@func",
			logrus.FieldKeyFile:  "@file",
			logrus.FieldKeyMsg:   "msg",
		},
		CallerPrettyfier: func(rf *runtime.Frame) (function string, file string) {
			file = fmt.Sprintf(":L%d", rf.Line)
			files := strings.Split(rf.File, "/")
			if len(files) > 0 {
				file = files[len(files)-1] + file
			}
			return rf.Function, file
		},
	})
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

func Close() {
	if logfile != nil {
		logfile.Close()
		logfile = nil
	}
}
