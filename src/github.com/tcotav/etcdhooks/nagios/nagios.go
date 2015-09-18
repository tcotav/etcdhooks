package nagios

import (
	"fmt"
	"github.com/tcotav/etcdhooks/logr"
	"os"
	"os/exec"
	"sort"
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

const ltagsrc = "etcdnagios"

// TODO: convert this to params
var nagiosCheckCmd = "/usr/sbin/nagios3"
var nagiosCheckArgs = []string{"-v", "/etc/nagios3/nagios.cfg"}
var nagiosPIDCmd = "pgrep"
var nagiosPIDArgs = []string{"nagios3"}
var nagiosHUPCmd = "kill"
var nagiosHUPArgs = []string{"-HUP"}

func execCmdOutput(cmdName string, cmdArgs []string) (string, error) {
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, fmt.Sprintf("cmd.exec:%s -- %s", cmdName, err))
		return "", err
	}
	return strings.TrimSpace(string(cmdOut)), nil
}

func execCmd(cmdName string, cmdArgs []string) error {
	cmd := exec.Command(cmdName, cmdArgs...)
	if err := cmd.Start(); err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, fmt.Sprintf("cmd.Start:%s -- %s", cmdName, err))
		return err
	}
	// check for non-zero exit code
	if err := cmd.Wait(); err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, fmt.Sprintf("cmd.Wait:%s -- %s", cmdName, err))
		return err
	}
	return nil
}

func RestartNagios() {
	if _, err := execCmdOutput(nagiosCheckCmd, nagiosCheckArgs); err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, "check nagios config failed")
	}
	logr.LogLine(logr.Linfo, ltagsrc, "check nagios succeeded")

	pid, err := execCmdOutput(nagiosPIDCmd, nagiosPIDArgs)
	if err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, "get nagios PID failed")
	}
	logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("got nagios pid: %s", pid))
	useArgs := append(nagiosHUPArgs, pid)
	if _, err := execCmdOutput(nagiosHUPCmd, useArgs); err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, "HUP nagios failed")
	}
	logr.LogLine(logr.Linfo, ltagsrc, "nagios restarted")
}

func extractGroup(s string) string {
	slist := strings.Split(s, "-")
	if len(slist) != 3 {
		logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("Invalid format: %s", s))
	}
	return slist[1]
}

// GenerateFiles takes the source host map and writes out a host and group nagios config file
// to the path passed to the function.
func GenerateFiles(hdMap map[string]string, hostPath string, groupPath string) {
	f, err := os.Create(hostPath)
	if err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, err.Error())
	}
	defer f.Close()
	hostGroups := make(map[string][]string)
	hostlist := make([]string, 0, len(hdMap))
	for host := range hdMap {
		hostlist = append(hostlist, host)
	}
	sort.Strings(hostlist)

	for _, h := range hostlist {
		// for each hostey in the map
		// write out a hostdef
		// append hostname to a group list
		f.WriteString(fmt.Sprintf(HostDef, h, h, h))
		group := extractGroup(h)
		hostGroups[group] = append(hostGroups[group], h)
	}
	// at the end, write out the group file using the group list

	f1, err := os.Create(groupPath)
	if err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, err.Error())
	}
	defer f1.Close()

	// now print out the group file
	for k := range hostGroups {
		sHosts := strings.Join(hostGroups[k], ",")
		f1.WriteString(fmt.Sprintf(GroupDef, k, k, sHosts))
	}

	go RestartNagios()
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
		logr.LogLine(logr.Lerror, ltagsrc, err.Error())
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
	logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("%v", hostGroups))

	f1, err := os.Create("/tmp/groups.cfg")
	if err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, err.Error())
	}
	defer f1.Close()

	// now print out the group file
	for k := range hostGroups {
		sHosts := strings.Join(hostGroups[k], ",")
		logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("group: %s, hosts: %s\n", k, sHosts))
		f1.WriteString(fmt.Sprintf(GroupDef, k, k, sHosts))
	}
}
