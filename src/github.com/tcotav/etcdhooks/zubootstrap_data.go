package main

import (
	"github.com/coreos/go-etcd/etcd"
	//"http"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func formatHostToList(s string) []string {
	if strings.Index(s, ".") != -1 { // we have fqdn
		s = strings.Split(s, ".")[0]
	}
	if strings.Index(s, "-") != -1 { // we have fqdn
		return strings.Split(s, "-")
	}
	return []string{}
}

func BootstrapData(fileName string, client *etcd.Client) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// remove leading space
		s := strings.TrimSpace(scanner.Text())
		// if starts with # -- skip
		if !strings.HasPrefix(s, "#") {
			slist := formatHostToList(s)
			if len(slist) == 3 {
				if _, err := client.Set(fmt.Sprintf("%s/%s", slist[0], slist[1]), slist[2], 0); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	config := config.ParseConfig("daemon.cfg")
	etcd_server_list := strings.Split(config["etcd_server_list"], ",")
	client := etcd.NewClient(etcd_server_list)
	BootstrapData("zubootstrap.txt", client)
	/*
	   spin through a list of hosts?
	   curl create all the hosts
	   use the created hostlist format -- one host one line
	*/
}
