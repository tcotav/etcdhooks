package nagios

import (
	"fmt"
	"log"
	"os"
	"strings"
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

func extractGroup(s string) string {
	slist := strings.Split(s, "-")
	if len(slist) != 3 {
		log.Fatal(fmt.Sprintf("Invalid format: %s", s))
	}
	return slist[1]
}

// GenerateFiles takes the source host map and writes out a host and group nagios config file
// to the path passed to the function.
func GenerateFiles(hdMap map[string]int, hostPath string, groupPath string) {
	f, err := os.Create(hostPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hostGroups := make(map[string][]string)
	//take a canned map of hosts, generate host and hostgroup files
	for host := range hdMap {
		// for each hostey in the map
		// write out a hostdef
		// append hostname to a group list
		f.WriteString(fmt.Sprintf(HostDef, host, host, host))

		group := extractGroup(host)
		hostGroups[group] = append(hostGroups[group], host)
	}
	// at the end, write out the group file using the group list

	f1, err := os.Create(groupPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	// now print out the group file
	for k := range hostGroups {
		sHosts := strings.Join(hostGroups[k], ",")
		f1.WriteString(fmt.Sprintf(GroupDef, k, k, sHosts))
	}
}

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

	hostGroups := make(map[string][]string)
	//take a canned map of hosts, generate host and hostgroup files
	for host := range fakeMap {
		// for each hostey in the map
		// write out a hostdef
		// append hostname to a group list
		f.WriteString(fmt.Sprintf(HostDef, host, host, host))

		group := extractGroup(host)
		fmt.Printf("%s\n", group)
		hostGroups[group] = append(hostGroups[group], host)
	}
	// at the end, write out the group file using the group list
	log.Printf("%v", hostGroups)

	f1, err := os.Create("/tmp/groups.cfg")
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	// now print out the group file
	for k := range hostGroups {
		sHosts := strings.Join(hostGroups[k], ",")
		log.Printf("group: %s, hosts: %s\n", k, sHosts)
		f1.WriteString(fmt.Sprintf(GroupDef, k, k, sHosts))
	}
}
