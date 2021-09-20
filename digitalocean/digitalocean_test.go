package digitalocean

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type fakeClient struct{}

func (fc fakeClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{Body: io.NopCloser(strings.NewReader(`{"hello":"there"}`)), StatusCode: 200}, nil
}

func TestServers(t *testing.T) {

	client := fakeClient{}

	do, err := NewDigitaloceanClient(client, "https://fake.url.com", "")

	if err != nil {
		t.Error(err)
		return
	}

	var serversModel map[string]interface{}

	err = do.Servers(&serversModel)

	if err != nil {
		t.Error(err)
		return
	}
}
