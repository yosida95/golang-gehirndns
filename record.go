package gehirndns

import (
	"time"
)

type (
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
