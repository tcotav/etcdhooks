package main

/*
Script that watched etcd and rewrites configuration files on change in etcd
*/

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/
import (
	"flag"
	"fmt"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/tcotav/etcdhooks/config"
	"github.com/tcotav/etcdhooks/etcd"
	"github.com/tcotav/etcdhooks/logr"
	"github.com/tcotav/etcdhooks/nagios"
	"github.com/tcotav/etcdhooks/web"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// think we want to dump a lot of this into a config
// stuff like the etcd info
//
var nagios_host_file = "/tmp/hosts.cfg"
var nagios_group_file = "/tmp/groups.cfg"
var host_list_file = "/tmp/host_list.cfg"

// updateHost wrapper containing async function calls to update the internal map
// as well as the config files
func updateHost(k string, v string) {
	hostMap := etcdWatcher.Map()
	_, containsHost := hostMap[k]
	go etcdWatcher.UpdateMap(k, v)
	// regenerate these files ONLY if it is a new host
	if !containsHost {
		go fixHostKey(k)
		regenHosts()
	}
}

const ltagsrc = "etcmain"

// TODO -- make this configurable
const sshkey_clean_path = "/opt/site-scripts/key_clean.sh"
const sshkey_clean_user = "nagios"

func fixHostKey(hostName string) {
	logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("fixhostkey for host: %s", hostName))
	hostName = strings.Replace(hostName[1:], "/", "-", -1)
	_, err := exec.Command(sshkey_clean_path, sshkey_clean_user, hostName).Output()
	if err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, fmt.Sprintf("key_clean for host %s failed", hostName))
	}
}

// dump hostmap out to a file
func writeHostMap(hostMap map[string]string) {
	if host_list_file == "" {
		return
	}
	f, err := os.Create(host_list_file)
	if err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, err.Error())
	}
	defer f.Close()

	for host := range hostMap {
		f.WriteString(fmt.Sprintf("%s\n", host))
	}
}

var limiterOn = false

var fileRewriteInterval = 15

var lastFileWrite = time.Now()

// regenHostFiles utility function that calls regen methods for files/persistence that contain only
// host data.  We pass along up/down and in/out service info too -- that should be handled with a different
// method.  Currently limited so that we don't write More than fileRewriteInterval seconds.
func regenHosts() {
	if limiterOn { // we're already waiting on a file rewrite
		logr.LogLine(logr.Linfo, ltagsrc, "limiter already on")
		return
	}

	// do some date math here -- have we waited long enough to write our file?
	// now < lastfilewrite + fileRewriteInterval
	if time.Now().Before(lastFileWrite.Add(time.Duration(fileRewriteInterval) * time.Second)) {
		logr.LogLine(logr.Linfo, ltagsrc, "limiter kicked in")
		limiterOn = true
		// these statements cause us to wait fileRewriteInterval seconds before continuing
		limiter := time.Tick(time.Duration(fileRewriteInterval) * time.Second)
		<-limiter
	}

	// flip back our counters
	limiterOn = false
	lastFileWrite = time.Now()

	logr.LogLine(logr.Linfo, ltagsrc, "generating files")
	// do the work
	etcdWatcher.BuildMap()
	hostMap := etcdWatcher.Map()
	nagios.GenerateFiles(hostMap, nagios_host_file, nagios_group_file)
	writeHostMap(hostMap)
}

func removeHost(k string) {
	logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("removeHost in daemon.go -- k:%s", k))
	regenHosts()
}

func GetEtcdKapi(serverList []string) (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints: serverList,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		logr.LogLine(logr.Lfatal, ltagsrc, err.Error())
		return nil, err
	}
	return client.NewKeysAPI(c), nil
}

func main() {

	// handle command line args
	var configFile string
	flag.StringVar(&configFile, "cfg", "daemon.cfg", "full path to daemon config")
	var logConfigPath string
	flag.StringVar(&logConfigPath, "logcfg", "none", "full path to log config")
	flag.Parse()

	if logConfigPath != "none" {
		logr.SetConfig(logConfigPath)
	}

	config, err := config.ParseConfig(configFile)
	if err != nil {
		log.Fatal("couldn't open config file %s", configFile, err)
	}
	nagios_host_file = config["nagios_host_file"]
	nagios_group_file = config["nagios_groups_file"]
	host_list_file = config["host_list_file"]
	watch_root := config["etcd_watch_root_url"]

	s := config["file_rewrite_interval"]
	if s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("invalid file rewrite val in config: ", err))
		}

		if i != 0 {
			fileRewriteInterval = i
		}
	}
	logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("file rewrite interval set to: %d", fileRewriteInterval))
	// expect this to be csv or single entry
	etcd_server_list := strings.Split(config["etcd_server_list"], ",")
	kapi, err := GetEtcdKapi(etcd_server_list)

	if err != nil {
		logr.LogLine(logr.Lfatal, ltagsrc, fmt.Sprintf("Error watching etcd", err.Error()))
		os.Exit(2)
	}
	logr.LogLine(logr.Linfo, ltagsrc, "got client")
	etcdWatcher.InitDataMap(kapi, watch_root)
	logr.LogLine(logr.Linfo, ltagsrc, "Dumping map contents for verification")
	etcdWatcher.DumpMap()
	logr.LogLine(logr.Linfo, ltagsrc, "Generating initial config files")
	regenHosts()
	//
	// spin up the web server
	//
	go webservice.StartWebService(config["web_listen_port"])
	watcherOpts := client.WatcherOptions{AfterIndex: 0, Recursive: true}
	w := kapi.Watcher(watch_root, &watcherOpts)
	logr.LogLine(logr.Linfo, ltagsrc, "Waiting for an update...")
	restartCount := 0
	for {
		r, err := w.Next(context.Background())
		if err != nil {
			logr.LogLine(logr.Lfatal, ltagsrc, fmt.Sprintf("Error watching etcd", err.Error()))
			// has etcd gone away?
			restartCount++
			switch {
			case restartCount < 10:
				logr.LogLine(logr.Lfatal, ltagsrc, "Sleeping for 10 seconds then retrying")
				time.Sleep(10 * time.Second)
			case restartCount < 20:
				time.Sleep(30 * time.Second)
				logr.LogLine(logr.Lfatal, ltagsrc, "Sleeping for 30 seconds then retrying.")
			default:
				time.Sleep(60 * time.Second * 5) // default sleep 5 minutes before retry
				logr.LogLine(logr.Lfatal, ltagsrc, "Sleeping for 5 minutes then retrying.")
			}
			kapi, err := GetEtcdKapi(etcd_server_list)
			if err != nil {
				logr.LogLine(logr.Lfatal, ltagsrc, fmt.Sprintf("Error getting etcdKAPI", err.Error()))
			}
			w = kapi.Watcher(watch_root, &watcherOpts)
			continue
		}
		// do something with it here
		action := r.Action
		k := r.Node.Key
		v := r.Node.Value
		switch action {
		case "delete":
			logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("delete of key: %s", k))
			go removeHost(k)
		case "set":
			logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("update of key: %s, value: %s", k, v))
			go updateHost(k, v)
		}
	}
}
