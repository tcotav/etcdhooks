package logr

/*
Script that watched etcd and rewrites configuration files on change in etcd
*/

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/
import (
	"github.com/Sirupsen/logrus"
	"os"
)

var log = logrus.New()

const Linfo = "info"
const Lfatal = "fatal"
const Lwarn = "lwarn"
const Ldebug = "debug"
const Lpanic = "panic"
const Lerror = "error"

func init() {
	// put overrides here
	return
}

func LogLine(lvl string, tagsrc string, o string) {
	l := log.WithFields(logrus.Fields{
		"src": tagsrc,
	})
	switch lvl {
	case Linfo:
		l.Info(o)
	case Lfatal:
		l.Fatal(o)
		os.Exit(3)
	case Lwarn:
		l.Warn(o)
	case Ldebug:
		l.Debug(o)
	case Lpanic:
		l.Panic(o)
		os.Exit(4)
	default:
		l.Info(o)
	}
}
