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
	client HttpClient
	URL    string //api url
	token  string //access token in order to login
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewDigitaloceanClient(client HttpClient, URL, token string) (Digitalocean, error) {

	return Digitalocean{client, URL, token}, nil
}

func (d Digitalocean) Servers(serversModel interface{}) error {

	uri := "/v2/droplets?page=1&per_page=20"

	p, err := url.Parse(d.URL)

	realURL := p.Scheme + "://" + path.Join(p.Host, uri)

	if err != nil {
		return fmt.Errorf("Couldnt parse URL: %s", realURL)
	}

	r, err := d.request("GET", realURL)

	if err != nil {
		return err
	}

	err = json.Unmarshal(r, &serversModel)

	if err != nil {
		return err
	}

	return err
}

func (d Digitalocean) request(method, url string) ([]byte, error) {

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Authorization", "Bearer "+d.token)

	res, err := d.client.Do(req)
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
