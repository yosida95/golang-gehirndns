package gehirndns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/kr/pretty"
)

const (
	APIENDPOINT = "https://cp.gehirn.jp/api/dns/"
)

type ZoneId uint

type ApiKey struct {
	Token  string
	Secret string
}

type Client struct {
	endpoint *url.URL
	apiKey   *ApiKey
}

func NewClient(zoneId ZoneId, apiKey *ApiKey) *Client {
	endpoint, err := url.Parse(APIENDPOINT)
	if err != nil {
		panic(err)
	}

	endpoint.Path = path.Join(
		endpoint.Path,
		"resource",
		strconv.Itoa(int(zoneId)))

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
			Error struct {
				Code    uint   `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}{}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&body)
		if err != nil {
			return fmt.Errorf(resp.Status)
		}

		return fmt.Errorf(body.Error.Message)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&body); err != nil {
		return
	}

	return
}

func (c *Client) GetResources() (err error) {
	req, err := c.makeRequest("GET", "", nil)
	if err != nil {
		return
	}

	body := struct {
		Resource struct {
			SOA   SOARecord
			NS    []NSRecord
			A     []ARecord
			AAAA  []AAAARecord
			CNAME []CNAMERecord
			MX    []MXRecord
			TXT   []TXTRecord
			SRV   []SRVRecord
		}
	}{}

	if err = c.request(req, &body); err != nil {
		return
	}

	pretty.Println(body.Resource)
	return
}

func (c *Client) AddNS(name, ns HostName, ttl Seconds) (record *NSRecord, err error) {
	const (
		recordType RecordType = "NS"
	)

	record = &NSRecord{
		NameServer: ns,
		Record: Record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = c.AddResource(record)
	return
}

func (c *Client) AddA(name HostName, addr IPv4, ttl Seconds) (record *ARecord, err error) {
	const (
		recordType RecordType = "A"
	)

	record = &ARecord{
		IPAddress: addr,
		Record: Record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = c.AddResource(record)
	return
}

func (c *Client) AddAAAA(name HostName, addr IPv6, ttl Seconds) (record *AAAARecord, err error) {
	const (
		recordType RecordType = "AAAA"
	)

	record = &AAAARecord{
		IPAddress: addr,
		Record: Record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = c.AddResource(record)
	return
}

func (c *Client) AddCNAME(name, to HostName, ttl Seconds) (record *CNAMERecord, err error) {
	const (
		recordType RecordType = "CNAME"
	)

	record = &CNAMERecord{
		AliasTo: to,
		Record: Record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = c.AddResource(record)
	return
}

func (c *Client) AddMX(name, mailServer HostName, priority Priority, ttl Seconds) (record *MXRecord, err error) {
	const (
		recordType RecordType = "MX"
	)

	record = &MXRecord{
		MailServer: mailServer,
		Priority:   priority,
		Record: Record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = c.AddResource(record)
	return
}

func (c *Client) AddTXT(name HostName, value string, ttl Seconds) (record *TXTRecord, err error) {
	const (
		recordType RecordType = "TXT"
	)

	record = &TXTRecord{
		Value: value,
		Record: Record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = c.AddResource(record)
	return
}

func (c *Client) AddSRV(name, target HostName, port, weight uint, priority Priority, ttl Seconds) (record *SRVRecord, err error) {
	const (
		recordType RecordType = "SRV"
	)

	record = &SRVRecord{
		Target:   target,
		Port:     port,
		Weight:   weight,
		Priority: priority,
		Record: Record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = c.AddResource(record)
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

func (c *Client) AddResource(record IRecord) (err error) {
	bodyObject := struct {
		Resource IRecord
	}{
		Resource: record,
	}
	body, err := c.encodeJSON(bodyObject)
	if err != nil {
		return
	}

	request, err := c.makeRequest("POST", "", body)
	if err != nil {
		return
	}

	responseBody := struct {
		Resource IRecord
	}{
		Resource: record,
	}
	err = c.request(request, &responseBody)
	return
}
