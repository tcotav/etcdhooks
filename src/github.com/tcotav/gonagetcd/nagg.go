package main

import (
	"fmt"
	"log"
	"os"
)

var HostDef = `define host {
         use             site-host
         host_name       %s
         alias           %s
         address         %s
         }

 `

var GroupDef = `define hostgroup {
         hostgroup_name  %s
         alias           %s
         members         %s
         }

`

func main() {

	var fakeMap = map[string]string{
		"site-web-100": "site-web-100",
		"site-web-200": "site-web-200",
		"site-web-300": "site-web-300",
		"site-db-100":  "site-db-100",
	}

	f, err := os.Create("/tmp/host.cfg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	/*
		take a canned map of hosts, generate host and hostgroup files
	*/
	for k := range fakeMap {
		// for each key in the map
		// write out a hostdef
		// append hostname to a group list
		f.WriteString(fmt.Sprintf(HostDef, k, k, k))
	}
	// at the end, write out the group file using the group list
}
