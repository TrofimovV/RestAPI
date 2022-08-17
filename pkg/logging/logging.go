package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

var logger *logrus.Entry

func GetLogger() *logrus.Entry {
	return logger
}

func init() {
	l := logrus.New()

	l.SetReportCaller(true)

	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%v()", frame.Function), fmt.Sprintf("%v :%v", filename, frame.Line)
		},
		FullTimestamp: true,
		DisableColors: false,
	}
	l.SetOutput(os.Stdout)

	l.SetLevel(logrus.TraceLevel)

	logger = logrus.NewEntry(l)
}
