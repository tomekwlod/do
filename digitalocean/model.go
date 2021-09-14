package digitalocean

type RemoteEntities struct {
	Droplet []Droplet `json:"droplets"`
}
type Droplet struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Networks Network  `json:"networks"`
	Tags     []string `json:"tags"`
}
type Network struct {
	Version []V4 `json:"v4"`
}
type V4 struct {
	IP   string `json:"ip_address"`
	Type string `json:"type"`
}
type Server struct {
	IP, Port, Name string
}

func (d Droplet) PublicIP() string {
	for _, v4 := range d.Networks.Version {
		if v4.Type == "public" {
			return v4.IP
		}
	}
	return ""
}
