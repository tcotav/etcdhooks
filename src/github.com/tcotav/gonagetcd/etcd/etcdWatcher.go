package etcdWatcher

/*
Script that watched etcd and rewrites configuration files on change in etcd
*/

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/
import (
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strconv"
	"strings"
)

type HostData struct {
	Name, IP string
	Status   int
}
type KVPair struct {
	Key, Value string
}

var hostMap = make(map[string]*HostData)

func Map() map[string]*HostData {
	return hostMap
}

// ClientGet gets data from etcd sending in an url and receiving a etcd.Response object
func ClientGet(client *etcd.Client, url string) *etcd.Response {
	resp, err := client.Get(url, true, true)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

// InitDataMap initializes a local map of hostnames and their respective metadata as
// struct: ip, status, name
func InitDataMap(client *etcd.Client) {
	baseStr := "/site"
	resp := ClientGet(client, baseStr)
	// get the list of host type
	for _, n := range resp.Node.Nodes {
		resp1 := ClientGet(client, n.Key)
		for _, n1 := range resp1.Node.Nodes {
			resp2 := ClientGet(client, n1.Key)
			for _, n2 := range resp2.Node.Nodes {
				// key format is /site/web/001 -- we want site-web-001
				hostName := strings.Replace(n1.Key[1:], "/", "-", -1)
				log.Print(hostName)
				if _, ok := hostMap[hostName]; !ok {
					var m = new(HostData)
					m.Name = hostName
					hostMap[hostName] = m
					log.Printf("Creating new entry for %s", hostName)
				}
				// key in format:  /site/extapi/500/status
				//log.Printf("n2.Key is %s", n2.Key)
				// want just the last part of url
				urlArray := strings.Split(n2.Key, "/")
				switch urlArray[len(urlArray)-1] {
				case "ip":
					hostMap[hostName].IP = n2.Value
				case "status":
					i, err := strconv.Atoi(n2.Value)
					if err != nil {
						// handle error
						log.Fatal(err)
					}
					hostMap[hostName].Status = i
				default:
					log.Fatal("Unknown key: " + n2.Key)
				}
			}
		}
	}
}

// DumpServices is a utility method that dumps all contents of etcd that match
// a specified base string
func DumpServices(client *etcd.Client, baseStr string) {
	//baseStr := "/site"
	resp := ClientGet(client, baseStr)
	// get the list of host type
	for _, n := range resp.Node.Nodes {
		resp1 := ClientGet(client, n.Key)
		for _, n1 := range resp1.Node.Nodes {
			resp2 := ClientGet(client, n1.Key)
			for _, n2 := range resp2.Node.Nodes {
				log.Printf("%s: %s\n", n2.Key, n2.Value)
			}
		}
	}
}

func DumpMap() {
	for k, v := range hostMap {
		log.Printf("%s: %+v\n", k, v)
	}
}

// UpdateMap updates the local hostmap with the given KV pair
func UpdateMap(k string, v string) {
	// format of key is /site/web/502/status
	keyArray := strings.Split(k[1:], "/")
	hostName := fmt.Sprintf("%s-%s-%s", keyArray[0], keyArray[1], keyArray[2])
	log.Printf("UpdateMap hostname: %s", hostName)
	param := keyArray[3]
	switch param {
	case "ip":
		hostMap[hostName].IP = v
	case "status":
		i, err := strconv.Atoi(v)
		if err != nil {
			// handle error
			log.Fatal(err)
		}
		hostMap[hostName].Status = i
	default:
		log.Fatal("Unknown key: " + k)
	}
}

func main() {
	/*
	  TODO:
	  - write code for generating our nagios config and our hostfile
	*/
	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	//hostMap =map[string]*HostData
	InitDataMap(client)
	watchChan := make(chan *etcd.Response)
	go client.Watch("/site/", 0, true, watchChan, nil)
	log.Println("Waiting for an update...")
	for {
		select {
		case r := <-watchChan:
			// do something with it here
			log.Printf("Updated KV: %s: %s\n", r.Node.Key, r.Node.Value)
			kvp := new(KVPair)
			kvp.Key = r.Node.Key
			kvp.Value = r.Node.Value
			go UpdateMap(kvp.Key, kvp.Value)
		}
	}
	// we don't really care what changed in this case so...
	//DumpServices(client)
}
