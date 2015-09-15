package main

/*
Script that watched etcd and rewrites configuration files on change in etcd
*/

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/
import (
	"fmt"
  "github.com/coreos/etcd/client"
	"github.com/tcotav/etcdhooks/config"
	"github.com/tcotav/etcdhooks/etcd"
	"github.com/tcotav/etcdhooks/nagios"
	"github.com/tcotav/etcdhooks/web"
	"log"
	"os"
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
		regenHosts()
	}
}

func writeHostMap(hostMap map[string]int) {
	f, err := os.Create(host_list_file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for host := range hostMap {
		f.WriteString(fmt.Sprintf("%s\n", host))
	}
}

var limiterOn = false

const fileRewriteInterval = 30

var lastFileWrite = time.Now().Add(time.Second * 1000 * -1) // initialize to some point in the past

// regenHostFiles utility function that calls regen methods for files/persistence that contain only
// host data.  We pass along up/down and in/out service info too -- that should be handled with a different
// method.  Currently limited so that we don't write More than fileRewriteInterval seconds.
func regenHosts() {
	if limiterOn { // we're already waiting on a file rewrite
		return
	}

	// do some date math here -- have we waited long enough to write our file?
	if time.Now().Before(lastFileWrite.Add(time.Second * fileRewriteInterval)) {
		log.Println("limiter kicked in")
		limiterOn = true
		// these statements cause us to wait fileRewriteInterval seconds before continuing
		limiter := time.Tick(time.Second * fileRewriteInterval)
		<-limiter
	}

	// flip back our counters
	limiterOn = false
	lastFileWrite = time.Now()

	log.Println("generating files")
	// do the work
	hostMap := etcdWatcher.Map()
	go nagios.GenerateFiles(hostMap, nagios_host_file, nagios_group_file)
	go writeHostMap(hostMap)
}

func removeHost(k string) {
	etcdWatcher.DeleteFromMap(k)
	// remove from map
	// run the updateNagios command
	regenHosts()
}

func main() {
	config := config.ParseConfig("daemon.cfg")
	nagios_host_file = config["nagios_host_file"]
	nagios_group_file = config["nagios_groups_file"]
	host_list_file = config["host_list_file"]

	// expect this to be csv or single entry
	etcd_server_list := strings.Split(config["etcd_server_list"], ",")
  cfg := client.Config{
      Endpoints:               etcd_server_list,
      Transport:               client.DefaultTransport,
      // set timeout per request to fail fast when the target endpoint is unavailable
      HeaderTimeoutPerRequest: time.Second,
  }
  c, err := client.New(cfg)
  if err != nil {
      log.Fatal(err)
  }
  kapi := client.NewKeysAPI(c)

	log.Println("got client")
	etcdWatcher.InitDataMap(kapi)
	log.Println("Dumping map contents for verification")
	etcdWatcher.DumpMap()
	log.Println("Generating initial config files")
	regenHosts()
	//
	// spin up the web server
	//
	go webservice.StartWebService(config["web_listen_port"])
	watchChan := make(chan *client.Response)
	watcherOpts := client.WatcherOptions{AfterIndex: 0, Recursive: true}
	go kapi.Watcher(config["base_etcd_url"], 0, true, watchChan, nil)
	log.Println("Waiting for an update...")
	for {
		select {
		case r := <-watchChan:
			// do something with it here
			action := r.Action
			k := r.Node.Key
			v := r.Node.Value
			switch action {
			case "delete":
				log.Printf("delete of key: %s", k)
				go removeHost(k)
			case "set":
				log.Printf("update of key: %s, value: %s", k, v)
				go updateHost(k, v)
			}
		}
	}
	// we don't really care what changed in this case so...
	//DumpServices(client)
}
