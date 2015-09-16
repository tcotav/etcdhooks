package main

import (
	"github.com/coreos/etcd/client"
	 "github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"log"
	"time"
)

func main() {
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:4001"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)
	/*resp, err := kapi.Set(context.Background(), "foo", "bar", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", resp)*/
	watcherOpts := client.WatcherOptions{AfterIndex: 0, Recursive: true}
	w := kapi.Watcher("/site", &watcherOpts)
	for {
		r,err := w.Next(context.Background()) 
		if err != nil {
			log.Fatal("Error occurred", err)
		}
		action := r.Action
	  log.Printf("action %s", action)
	}
}
