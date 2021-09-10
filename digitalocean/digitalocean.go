package digitalocean

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Digitalocean struct {
	URL   string
	token string
}

func NewDigitaloceanClient(URL, token string) (Digitalocean, error) {

	// url check here: u, err := url.Parse("http://foo")
	// auth

	return Digitalocean{URL, token}, nil
}

func (d Digitalocean) Servers() (io.ReadCloser, error) {

	return d.request("GET", "/v2/droplets?page=1&per_page=20")
}

func (d Digitalocean) request(method, uri string) (io.ReadCloser, error) {

	// check method here
	//
	//

	p, err := url.Parse(d.URL + "/" + uri)

	if err != nil {
		return nil, fmt.Errorf("Couldnt parse URL: %s", d.URL+"/"+uri)
	}

	req, err := http.NewRequest(method, p.String(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Authorization", "Bearer "+d.token)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return res.Body, nil
}
