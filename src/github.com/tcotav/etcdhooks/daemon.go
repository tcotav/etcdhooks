package main

/*
Script that watched etcd and rewrites configuration files on change in etcd
*/

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/
import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/tcotav/etcdhooks/config"
	"github.com/tcotav/etcdhooks/etcd"
	"github.com/tcotav/etcdhooks/nagios"
	"github.com/tcotav/etcdhooks/web"
	"log"
	"strings"
  "os"
  "fmt"
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
	go etcdWatcher.UpdateMap(k, v)
	// run the updateNagios command
	regenFiles()
}

func writeHostMap(hostMap map[string]int) {
  f, err := os.Create(host_list_file)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  for host:= range hostMap {
    f.WriteString(fmt.Sprintf("%s\n",host))
    }
}

// regenFiles utility function that calls ALL of the file regen methods.
// Currently only handles nagios
func regenFiles() {
  hostMap :=etcdWatcher.Map()
	go nagios.GenerateFiles(hostMap, nagios_host_file, nagios_group_file)
  go writeHostMap(hostMap)
}

func removeHost(k string) {
	go etcdWatcher.DeleteFromMap(k)
	// remove from map
	// run the updateNagios command
	regenFiles()
}

func main() {
	config := config.ParseConfig("daemon.cfg")
	nagios_host_file = config["nagios_host_file"]
	nagios_group_file = config["nagios_groups_file"]
  host_list_file = config["host_list_file"]

	// expect this to be csv or single entry
	etcd_server_list := strings.Split(config["etcd_server_list"], ",")
	client := etcd.NewClient(etcd_server_list)
	etcdWatcher.InitDataMap(client)
	//log.Println("Dumping map contents for verification")
	//etcdWatcher.DumpMap()
	log.Println("Generating initial config files")
	regenFiles()
  //
  // spin up the web server
  //
  go webservice.StartWebService(config["web_listen_port"])
	watchChan := make(chan *etcd.Response)
	go client.Watch(config["base_etcd_url"], 0, true, watchChan, nil)
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
