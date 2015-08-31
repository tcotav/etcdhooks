package main

/*
Script that watched etcd and rewrites configuration files on change in etcd
*/

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/
import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/tcotav/gonagetcd/etcdWatcher"
	//	"github.com/tcotav/gonagetcd/nagios"
	"log"
)

func main() {
	//  TODO: write code for generating our nagios config and our hostfile
	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	//hostMap =map[string]*HostData
	etcdWatcher.InitDataMap(client)
	watchChan := make(chan *etcd.Response)
	go client.Watch("/site/", 0, true, watchChan, nil)
	log.Println("Waiting for an update...")
	for {
		select {
		case r := <-watchChan:
			// do something with it here
			log.Printf("Updated KV: %s: %s\n", r.Node.Key, r.Node.Value)
			kvp := new(etcdWatcher.KVPair)
			kvp.Key = r.Node.Key
			kvp.Value = r.Node.Value
			go etcdWatcher.UpdateMap(kvp.Key, kvp.Value)
		}
	}
	// we don't really care what changed in this case so...
	//DumpServices(client)
}
