package gehirndns

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	APIENDPOINT = "https://cp.gehirn.jp/api/dns/"
)

type ApiKey struct {
	Token  string
	Secret string
}

type Client struct {
	endpoint *url.URL
	apiKey   *ApiKey
}

func NewClient(apiKey *ApiKey) *Client {
	endpoint, err := url.Parse(APIENDPOINT)
	if err != nil {
		panic(err)
	}

	return &Client{
		endpoint: endpoint,
		apiKey:   apiKey,
	}
}

func (c *Client) buildURL(relativePath string) (endpoint *url.URL) {
	endpoint = new(url.URL)
	*endpoint = *c.endpoint
	endpoint.Path = path.Join(endpoint.Path, relativePath)
	return
}

func (c *Client) makeRequest(method, path string, body io.Reader) (req *http.Request, err error) {
	endpoint := c.buildURL(path)
	req, err = http.NewRequest(method, endpoint.String(), body)
	if err != nil {
		return
	}

	req.SetBasicAuth(c.apiKey.Token, c.apiKey.Secret)
	req.Header.Add("Content-Type", "text/json;charset=utf8")
	return
}

func (c *Client) request(req *http.Request, body interface{}) (err error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		body := struct {
			Error errorResponse `json:"error"`
		}{}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&body)
		if err != nil {
			return errorResponse{
				Code:    resp.StatusCode,
				Message: resp.Status,
			}
		}

		return body.Error
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	return
}

func (c *Client) encodeJSON(object interface{}) (reader io.Reader, err error) {
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	err = encoder.Encode(object)
	if err != nil {
		return
	}

	reader = buffer
	return
}

func (c *Client) GetZone(id ZoneId) (zone *Zone, err error) {
	req, err := c.makeRequest("GET", path.Join("resource", id.String()), nil)
	if err != nil {
		return
	}

	zone = &Zone{
		id:     id,
		client: c,
	}

	body := struct {
		Resource *Zone
		Domain   struct {
			Name HostName
		}
	}{
		Resource: zone,
	}

	err = c.request(req, &body)
	if err != nil {
		return
	}

	zone.postGet(body.Domain.Name)
	return
}
