package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// think we want to dump a lot of this into a config
// stuff like the etcd info
//

func ParseConfig(fileName string) map[string]string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
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
				log.Printf("%s", slist[0])
				config[slist[0]] = slist[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return config
}

func main() {
	config := ParseConfig("daemon.cfg")
	log.Print("%v", config)
}
