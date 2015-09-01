package main

/*
Script that watched etcd and rewrites configuration files on change in etcd
*/

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/
import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/tcotav/gonagetcd/etcd"
	//	"github.com/tcotav/gonagetcd/nagios"
	"log"
)

// updateHost wrapper containing async function calls to update the internal map
// as well as the config files
func updateHost(k string, v string) {
	go etcdWatcher.UpdateMap(k, v)
	// run the updateNagios command
	go nagios.GenerateFiles(etcdWatcher.Map(), k, v)
}

func removeHost(k string) {
	go etcdWatcher.DeleteFromMap(k)
	// remove from map
	// run the updateNagios command
	go nagios.GenerateFiles(etcdWatcher.Map(), k, v)
	log.Printf("in delete for key:%s\n", k)
}

func main() {
	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	etcdWatcher.InitDataMap(client)
	log.Println("Dumping map contents for verification")
	etcdWatcher.DumpMap()
	watchChan := make(chan *etcd.Response)
	go client.Watch("/site/", 0, true, watchChan, nil)
	log.Println("Waiting for an update...")
	for {
		select {
		case r := <-watchChan:
			// do something with it here
			log.Printf("Changed KV: %+v\n", r)
			log.Printf("Updated KV: %s: %s\n", r.Node.Key, r.Node.Value)
			action := r.Action
			k := r.Node.Key
			v := r.Node.Value
			switch action {
			case "delete":
				log.Printf("delete of key: %s", k)
				go removeHost(k)
			case "set":
				go updateHost(k, v)
			}
		}
	}
	// we don't really care what changed in this case so...
	//DumpServices(client)
}
