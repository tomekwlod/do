package digitalocean

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type Digitalocean struct {
	URL   string
	token string
}

func NewDigitaloceanClient(URL, token string) (Digitalocean, error) {

	return Digitalocean{URL, token}, nil
}

func (d Digitalocean) Servers(serversModel interface{}) error {

	r, err := d.request("GET", "/v2/droplets?page=1&per_page=20")

	if err != nil {
		return err
	}

	err = json.Unmarshal(r, &serversModel)

	if err != nil {
		return err
	}

	return err
}

func (d Digitalocean) request(method, uri string) ([]byte, error) {

	// check method here
	//
	//

	p, err := url.Parse(d.URL)

	realURL := p.Scheme + "://" + path.Join(p.Host, uri)

	if err != nil {
		return nil, fmt.Errorf("Couldnt parse URL: %s", realURL)
	}

	req, err := http.NewRequest(method, realURL, nil)

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(body))

	return body, nil
}
