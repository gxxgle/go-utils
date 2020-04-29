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
	L = logrus.WithField("@pid", os.Getpid())
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

func File(logfiles ...string) {
	logfile := ""
	if len(logfiles) > 0 {
		logfile = logfiles[0]
	}

	if len(logfile) == 0 {
		dir := path.TopLevelDir(path.CurrentDir())
		filename := path.CurrentFilename()
		logfile = fmt.Sprintf("%s/log/%s.log", dir, filename)
	}

	if err := os.MkdirAll(filepath.Dir(logfile), os.ModePerm); err != nil {
		L.WithField("path", logfile).WithError(err).Error("go-utils log mkdir failed")
		return
	}

	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		L.WithField("path", logfile).WithError(err).Error("go-utils log open file failed")
		return
	}

	logrus.SetOutput(file)
}

func Debug() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
}
