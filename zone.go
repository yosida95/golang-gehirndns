package gehirndns

import (
	"io"
	"net/http"
	"path"
	"reflect"
	"strconv"
)

type ZoneId int

func (id ZoneId) Int() int {
	return int(id)
}

func (id ZoneId) String() string {
	return strconv.Itoa(id.Int())
}

type HostName string

func (host HostName) String() string {
	return string(host)
}

type Zone struct {
	id     ZoneId
	client *Client

	SOA   *SOARecord
	NS    []*NSRecord
	A     []*ARecord
	AAAA  []*AAAARecord
	CNAME []*CNAMERecord
	MX    []*MXRecord
	TXT   []*TXTRecord
	SRV   []*SRVRecord
}

func (z *Zone) postGetRecord(record IRecord, domain HostName) {
	record.setHostNameByName(domain)
}

func (z *Zone) postGetRecords(records interface{}, domain HostName) {
	func(recordsV reflect.Value) {
		for i := 0; i < recordsV.Len(); i++ {
			recordV := recordsV.Index(i)
			record := recordV.Interface().(IRecord)
			z.postGetRecord(record, domain)
		}
	}(reflect.ValueOf(records))
}

func (z *Zone) postGet(domain HostName) {
	z.postGetRecord(z.SOA, domain)
	z.postGetRecords(z.NS, domain)
	z.postGetRecords(z.A, domain)
	z.postGetRecords(z.AAAA, domain)
	z.postGetRecords(z.CNAME, domain)
	z.postGetRecords(z.MX, domain)
	z.postGetRecords(z.TXT, domain)
	z.postGetRecords(z.SRV, domain)
}

func (z *Zone) makeRequest(method string, id RecordId, body io.Reader) (*http.Request, error) {
	endpoint := path.Join("resource", z.id.String(), id.String())
	return z.client.makeRequest(method, endpoint, body)
}

func (z *Zone) GetId() ZoneId {
	return z.id
}

func (z *Zone) AddNS(name, ns HostName, ttl Seconds) (resource *NSRecord, err error) {
	const (
		recordType RecordType = "NS"
	)

	resource = &NSRecord{
		NameServer: ns,
		record: record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = z.AddResource(resource)
	return
}

func (z *Zone) AddA(name HostName, addr IPv4, ttl Seconds) (resource *ARecord, err error) {
	const (
		recordType RecordType = "A"
	)

	resource = &ARecord{
		IPAddress: addr,
		record: record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = z.AddResource(resource)
	return
}

func (z *Zone) AddAAAA(name HostName, addr IPv6, ttl Seconds) (resource *AAAARecord, err error) {
	const (
		recordType RecordType = "AAAA"
	)

	resource = &AAAARecord{
		IPAddress: addr,
		record: record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = z.AddResource(resource)
	return
}

func (z *Zone) AddCNAME(name, to HostName, ttl Seconds) (resource *CNAMERecord, err error) {
	const (
		recordType RecordType = "CNAME"
	)

	resource = &CNAMERecord{
		AliasTo: to,
		record: record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = z.AddResource(resource)
	return
}

func (z *Zone) AddMX(name, mailServer HostName, priority Priority, ttl Seconds) (resource *MXRecord, err error) {
	const (
		recordType RecordType = "MX"
	)

	resource = &MXRecord{
		MailServer: mailServer,
		Priority:   priority,
		record: record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = z.AddResource(resource)
	return
}

func (z *Zone) AddTXT(name HostName, value string, ttl Seconds) (resource *TXTRecord, err error) {
	const (
		recordType RecordType = "TXT"
	)

	resource = &TXTRecord{
		Value: value,
		record: record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = z.AddResource(resource)
	return
}

func (z *Zone) AddSRV(name, target HostName, port, weight uint, priority Priority, ttl Seconds) (resource *SRVRecord, err error) {
	const (
		recordType RecordType = "SRV"
	)

	resource = &SRVRecord{
		Target:   target,
		Port:     port,
		Weight:   weight,
		Priority: priority,
		record: record{
			HostName: name,
			Type:     recordType,
			TTL:      ttl,
		},
	}

	err = z.AddResource(resource)
	return
}

func (z *Zone) AddResource(record IRecord) (err error) {
	if record.GetId() != "" {
		return ErrMaybeRegistered
	}

	bodyObject := struct {
		Resource IRecord
	}{
		Resource: record,
	}
	body, err := z.client.encodeJSON(bodyObject)
	if err != nil {
		return
	}

	request, err := z.makeRequest("POST", "", body)
	if err != nil {
		return
	}

	err = z.client.request(request, &bodyObject)
	return
}

func (z *Zone) UpdateResource(record IRecord) (err error) {
	if record.GetId() == "" {
		return ErrIdUnset
	}
	record.clearName()

	bodyObject := struct {
		Resource IRecord
	}{
		Resource: record,
	}
	body, err := z.client.encodeJSON(bodyObject)
	if err != nil {
		return
	}

	request, err := z.makeRequest("PUT", record.GetId(), body)
	if err != nil {
		return
	}

	err = z.client.request(request, &bodyObject)
	return
}

func (z *Zone) DeleteResource(record IRecord) (err error) {
	if record.GetId() == "" {
		return ErrIdUnset
	}

	request, err := z.makeRequest("DELETE", record.GetId(), nil)
	if err != nil {
		return
	}

	responseBody := struct {
		Resource IRecord
	}{
		Resource: record,
	}
	err = z.client.request(request, &responseBody)
	return
}
