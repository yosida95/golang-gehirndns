package gehirndns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/kr/pretty"
)

const (
	APIENDPOINT = "https://cp.gehirn.jp/api/dns/"
)

type (
	ZoneId      uint
	IPv4        string
	IPv6        string
	HostName    string
	MailAddress string
	RecordId    string
	RecordType  string
	Seconds     int64
	Priority    uint64
)

func (s Seconds) Duration() time.Duration {
	return time.Duration(s) * time.Second
}

type IRecord interface {
	GetId() RecordId
	GetHostName() HostName
	GetType() RecordType
	GetTTL() Seconds
}

type Record struct {
	Id       RecordId `json:"ID,omitempty"`
	Name     HostName `json:",omitempty"`
	HostName HostName
	Type     RecordType `json:"Type"`
	TTL      Seconds    `json:"TTL"`
}

func (r *Record) GetId() RecordId {
	return r.Id
}

func (r *Record) GetHostName() HostName {
	if len(r.Name) > 0 {
		return r.Name
	}

	return r.HostName
}

func (r *Record) GetType() RecordType {
	return r.Type
}

func (r *Record) GetTTL() Seconds {
	return r.TTL
}

type SOARecord struct {
	Mname            HostName    `json:"MNAME"`
	Rname            MailAddress `json:"RNAME"`
	Serial           uint
	Refresh          Seconds
	Retry            Seconds
	Expire           Seconds
	NegativeCacheTTL Seconds
	Record
}

type NSRecord struct {
	NameServer HostName
	Record
}

type ARecord struct {
	IPAddress IPv4 `json:"IPAddress"`
	Record
}

type AAAARecord struct {
	IPAddress IPv6 `json:"IPAddress"`
	Record
}

type CNAMERecord struct {
	AliasTo HostName `json:"AliasTo"`
	Record
}

type MXRecord struct {
	MailServer HostName
	Priority   Priority
	Record
}

type TXTRecord struct {
	Value string
	Record
}

type SRVRecord struct {
	Priority Priority
	Weight   uint
	Port     uint
	Target   HostName
	Record
}

type Client struct {
	endpoint  *url.URL
	apiToken  string
	apiSecret string
}

func NewClient(zoneId ZoneId, apiToken, apiSecret string) *Client {
	endpoint, err := url.Parse(APIENDPOINT)
	if err != nil {
		panic(err)
	}

	endpoint.Path = path.Join(
		endpoint.Path,
		"resource",
		strconv.Itoa(int(zoneId)))

	return &Client{
		endpoint:  endpoint,
		apiToken:  apiToken,
		apiSecret: apiSecret,
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

	req.SetBasicAuth(c.apiToken, c.apiSecret)
	return
}

func (c *Client) request(req *http.Request, body interface{}) (err error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
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
		Successful bool `json:"is_success"`
	}{}

	if err = c.request(req, &body); err != nil {
		return
	} else if !body.Successful {
		return fmt.Errorf("unsuccessful")
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

	c.AddResource(record)
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

	c.AddResource(record)
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

	c.AddResource(record)
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

	c.AddResource(record)
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

	c.AddResource(record)
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

	c.AddResource(record)
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

	c.AddResource(record)
	return
}

func (c *Client) AddResource(record IRecord) (err error) {
	request := struct {
		Resource IRecord
	}{
		Resource: record,
	}

	body, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		return
	}

	fmt.Println(string(body))
	return
}
