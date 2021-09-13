package digitalocean

type RemoteEntities struct {
	Droplets []Droplets `json:"droplets"`
}
type Droplets struct {
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
	IP string `json:"ip_address"`
}
type Server struct {
	IP, Port, Name string
}
