package main

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"github.com/tcotav/etcdhooks/config"
	"github.com/tcotav/etcdhooks/etcd"
	"net/http"
	"strings"
)

func dump(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(etcdWatcher.Map())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	config := config.ParseConfig("daemon.cfg")
	// expect this to be csv or single entry
	etcd_server_list := strings.Split(config["etcd_server_list"], ",")
	client := etcd.NewClient(etcd_server_list)
	etcdWatcher.InitDataMap(client)

	listenPort := config["web_listen_port"]
	http.HandleFunc("/", dump)
	http.ListenAndServe(fmt.Sprintf(":%s", listenPort), nil)
}
