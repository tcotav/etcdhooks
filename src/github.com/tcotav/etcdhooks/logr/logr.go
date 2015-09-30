package logr

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/tcotav/etcdhooks/config"
	"os"
	"runtime"
)

var stackTrace = false
var configFile = "log.cfg"
var log = logrus.New()

const Linfo = "info"
const Lfatal = "fatal"
const Lwarn = "lwarn"
const Ldebug = "debug"
const Lpanic = "panic"
const Lerror = "error"

// SetConfig sets the path to the log config file we want to use AND resets the
// module to use it.
func SetConfig(path string) {
	configFile = path
	log = logrus.New()
}

func setLogLevel(lvl string) {
	switch lvl {
	case Linfo:
		logrus.SetLevel(logrus.InfoLevel)
	case Lerror:
		logrus.SetLevel(logrus.ErrorLevel)
	case Lfatal:
		logrus.SetLevel(logrus.FatalLevel)
	case Lwarn:
		logrus.SetLevel(logrus.WarnLevel)
	case Ldebug:
		logrus.SetLevel(logrus.DebugLevel)
	case Lpanic:
		logrus.SetLevel(logrus.PanicLevel)
	}
}

func init() {
	logcfg, err := config.ParseConfig(configFile)
	if err != nil {
		log.Printf("%s logging config not found")
		return
	}

	// do we want to dump stack traces?
	s, _ := logcfg["stacktrace"]

	if s == "true" {
		stackTrace = true
	}

	// put overrides here
	logLevel := logcfg["loglevel"]
	if logLevel != "" {
		setLogLevel(logLevel)
	}
	return
}

func DumpStackTrace(lvl string, tagsrc string, msg string) {
	if stackTrace {
		//stack trace
		var stack [4096]byte
		runtime.Stack(stack[:], false)
		logLine(lvl, tagsrc, fmt.Sprintf("%s", stack[:]))
	}
}

func LogLine(lvl string, tagsrc string, msg string) {
	logLine(lvl, tagsrc, msg)
	// additional actions
	switch lvl {
	case Lerror:
		DumpStackTrace(lvl, tagsrc, msg)
	case Lfatal:
		DumpStackTrace(lvl, tagsrc, msg)
		os.Exit(3)
	}
}

// logLine only handles the logging to target
func logLine(lvl string, tagsrc string, msg string) {
	l := log.WithFields(logrus.Fields{
		"src": tagsrc,
	})
	switch lvl {
	case Linfo:
		l.Info(msg)
	case Lerror:
		l.Error(msg)
	case Lfatal:
		l.Fatal(msg)
	case Lwarn:
		l.Warn(msg)
	case Ldebug:
		l.Debug(msg)
	case Lpanic:
		l.Panic(msg)
	default:
		l.Info(msg)
	}
}
