package logr

/*
Script that watched etcd and rewrites configuration files on change in etcd
*/

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/
import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/tcotav/etcdhooks/config"
	"os"
	"runtime"
)

var log = logrus.New()

const Linfo = "info"
const Lfatal = "fatal"
const Lwarn = "lwarn"
const Ldebug = "debug"
const Lpanic = "panic"
const Lerror = "error"

var stackTrace = false

func init() {
	logcfg := config.ParseConfig("log.cfg")

	// do we want to dump stack traces?
	stackTrace, _ := logcfg["stacktrace"]

	// put overrides here
	return
}

func LogFatal(tagsrc string, msg string) {
	l := log.WithFields(logrus.Fields{
		"src": tagsrc,
	})

	l.Fatal(msg)

	if stackTrace {
		lstack := log.WithFields(logrus.Fields{
			"src":  tagsrc,
			"func": functionSrc,
			"data": "stack",
		})
		//stack trace
		var stack [4096]byte
		runtime.Stack(stack[:], false)
		lstack.Fatal(fmt.Sprintf("%s", stack[:]))
	}
}

func LogLine(lvl string, tagsrc string, msg string) {
	l := log.WithFields(logrus.Fields{
		"src": tagsrc,
	})
	switch lvl {
	case Linfo:
		l.Info(msg)
	case Lfatal:
		l.Fatal(msg)
		os.Exit(3)
	case Lwarn:
		l.Warn(msg)
	case Ldebug:
		l.Debug(msg)
	case Lpanic:
		l.Panic(msg)
		os.Exit(4)
	default:
		l.Info(msg)
	}
}
