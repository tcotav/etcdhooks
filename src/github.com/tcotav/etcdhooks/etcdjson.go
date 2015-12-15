package etcdhooks

import (
	"encoding/json"
	"flag"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/tcotav/etcdhooks/config"
	"github.com/tcotav/etcdhooks/web"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func jsonFromFile(fileName string) []webservice.HostState {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var stateList []webservice.HostState
	json.Unmarshal(file, &stateList)
	return stateList
}

func LoadJson(config map[string]string, sourceFile string) {
	// get the server list
	etcd_server_list := strings.Split(config["etcd_server_list"], ",")
	cfg := client.Config{
		Endpoints: etcd_server_list,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	kapi := client.NewKeysAPI(c)
	stateList := jsonFromFile(sourceFile)
	for _, hostState := range stateList {
		//log.Print(hostState.Name)
		resp, err := kapi.Set(context.Background(), strings.Replace(hostState.Name, "-", "/", -1), hostState.State, nil)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Print(resp)
		}
	}
}

/*
Take a json file from the file system, parse, and tick through it one by one
to populate/SET values in etcd
*/
func main() {
	var configFile string
	flag.StringVar(&configFile, "cfg", "daemon.cfg", "full path to daemon config")

	var jsonFile string
	flag.StringVar(&jsonFile, "dump", "dump.json", "full path to json host dump")
	flag.Parse()

	config, err := config.ParseConfig(configFile)
	if err != nil {
		log.Fatal("couldn't open config file", err)
	}
	LoadJson(config, jsonFile)
}
