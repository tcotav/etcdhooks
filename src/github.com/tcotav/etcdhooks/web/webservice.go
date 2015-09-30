package webservice

import (
	"encoding/json"
	"fmt"
	"github.com/tcotav/etcdhooks/etcd"
	"github.com/tcotav/etcdhooks/logr"
	"net/http"
)

const ltagsrc = "etcdweb"

var serviceMap = map[string]string{
	"/getall": "Returns json dict in form of hostname:state",
}

type HostState struct {
	Name  string
	State string
}

func services(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(serviceMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func mapToJson(mmap map[string]string) ([]byte, error) {

	mapList := []HostState{}
	for k, v := range mmap {
		mapList = append(mapList, HostState{k, v})
	}

	bData, err := json.Marshal(mapList)
	if err != nil {
		return nil, err
	}
	return bData, nil
}

func getAll(w http.ResponseWriter, r *http.Request) {
	js, err := mapToJson(etcdWatcher.Map())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func load(w http.ResponseWriter, r *http.Request) {
	// we want to take in a json  map
	// and replace the internal data structure with
	// that

	// also, update etcd with the data?
	/*
		js, err := json.Marshal(etcdWatcher.Map())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)*/
}

func StartWebService(listenPort string) {
	http.HandleFunc("/getall", getAll)
	//http.HandleFunc("/loadHosts", loadHosts)
	http.HandleFunc("/", services)
	logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("Starting webservice on port: %s", listenPort))
	http.ListenAndServe(fmt.Sprintf(":%s", listenPort), nil)
}

/*func main() {
	config := config.ParseConfig("daemon.cfg")
	// expect this to be csv or single entry
	etcd_server_list := strings.Split(config["etcd_server_list"], ",")
	client := etcd.NewClient(etcd_server_list)
	etcdWatcher.InitDataMap(client)

	listenPort := config["web_listen_port"]
	http.HandleFunc("/", dump)
	http.ListenAndServe(fmt.Sprintf(":%s", listenPort), nil)
}*/
