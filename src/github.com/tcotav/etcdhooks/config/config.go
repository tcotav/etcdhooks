package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// ParseConfig parse a simple K=V pair based config file and return
// a map equivalent.
func ParseConfig(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Print(err)
		return nil, err
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
			if len(slist) >= 2 {
				// bit of a hack in case equals sign appears as part of the value
				config[slist[0]] = strings.Join(slist[1:], "=")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return config, nil
}

func main() {
	config, err := ParseConfig("daemon.cfg")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", config)
}
