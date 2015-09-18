package etcdWatcher

/*
Script that watched etcd and rewrites configuration files on change in etcd
*/

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/
import (
	"fmt"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/tcotav/etcdhooks/logr"
	"strings"
)

var hostMap map[string]string

// Map returns the hostmap
func Map() map[string]string {
	return hostMap
}

var clientGetOpts = client.GetOptions{Recursive: true, Sort: true}

const ltagsrc = "etcwatc"

// ClientGet gets data from etcd sending in an url and receiving a etcd.Response object
func ClientGet(kapi client.KeysAPI, url string) *client.Response {
	resp, err := kapi.Get(context.Background(), url, &clientGetOpts)
	if err != nil {
		logr.LogLine(logr.Lerror, ltagsrc, err.Error())
	}
	return resp
}

func DeleteFromMap(k string) {
	logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("deleting key: %s", k))
	delete(hostMap, k)
}

var kapiClient client.KeysAPI
var rootUrl string

// InitDataMap initializes a local map of hostnames and their respective metadata as
// struct: ip, status, name
func InitDataMap(kapi client.KeysAPI, baseStr string) {
	kapiClient = kapi
	rootUrl = baseStr
	BuildMap()
}

func BuildMap() {
	resp := ClientGet(kapiClient, rootUrl)
	hostMap = make(map[string]string)
	// get the list of host type
	for _, n := range resp.Node.Nodes {
		resp1 := ClientGet(kapiClient, n.Key)
		for _, n1 := range resp1.Node.Nodes {
			// key format is /site/web/001 -- we want site-web-001
			hostName := strings.Replace(n1.Key[1:], "/", "-", -1)
			hostMap[hostName] = n1.Value
		}
	}
}

// DumpServices is a utility method that dumps all contents of etcd that match
// a specified base string
func DumpServices(kapi client.KeysAPI, baseStr string) {
	//baseStr := "/site"
	resp := ClientGet(kapi, baseStr)
	// get the list of host type
	for _, n := range resp.Node.Nodes {
		resp1 := ClientGet(kapi, n.Key)
		for _, n1 := range resp1.Node.Nodes {
			logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("%s: %s", n1.Key, n1.Value))
		}
	}
}

// DumpMap walks the host map and dumps out key-value pairs
func DumpMap() {
	for k, v := range hostMap {
		logr.LogLine(logr.Linfo, ltagsrc, fmt.Sprintf("%s: %+v\n", k, v))
	}
}

// UpdateMap updates the local hostmap with the given KV pair
func UpdateMap(k string, v string) {
	// format of key is /site/web/502/status
	keyArray := strings.Split(k[1:], "/")
	hostName := fmt.Sprintf("%s-%s-%s", keyArray[0], keyArray[1], keyArray[2])
	hostMap[hostName] = v
}

func main() {
	/*
		  TODO:
		  - write code for generating our nagios config and our hostfile
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
				}
		}
		// we don't really care what changed in this case so...
		//DumpServices(client)
	*/
}
