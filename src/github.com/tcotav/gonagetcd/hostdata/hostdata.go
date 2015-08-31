package hostdata

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
