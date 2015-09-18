package config

import (
	"bufio"
	"fmt"
	"github.com/tcotav/etcdhooks/logr"
	"os"
	"strings"
)

const ltagsrc = "etcconf"

// ParseConfig parse a simple K=V pair based config file and return
// a map equivalent.
func ParseConfig(fileName string) map[string]string {
	file, err := os.Open(fileName)
	if err != nil {
		logr.LogLine(logr.Lfatal, ltagsrc, err.Error())
	}
	defer file.Close()

	config := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// remove leading space
		s := strings.TrimSpace(scanner.Text())
		// if starts with # -- skip
		// test if contains = sign
		if !strings.HasPrefix(s, "#") && strings.Index(s, "=") != -1 {
			// split on = sign
			slist := strings.Split(s, "=")
			// set map[k] = v
			if len(slist) == 2 {
				config[slist[0]] = slist[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logr.LogLine(logr.Lfatal, ltagsrc, err.Error())
	}
	return config
}

func main() {
	config := ParseConfig("daemon.cfg")
	logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprint("%v", config))
}
