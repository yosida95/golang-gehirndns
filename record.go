package gehirndns

import (
	"time"
)

type (
	IPv4        string
	IPv6        string
	HostName    string
	MailAddress string
	RecordType  string
	Priority    uint64
)

type Seconds int64

func (s Seconds) Duration() time.Duration {
	return time.Duration(s) * time.Second
}

type RecordId string

func (id RecordId) String() string {
	return string(id)
}

type IRecord interface {
	GetId() RecordId
	GetHostName() HostName
	SetHostName(HostName)
	GetName() HostName
	clearName()
	GetType() RecordType
	GetTTL() Seconds
	SetTTL(Seconds)
}

type record struct {
	Id       RecordId `json:"ID,omitempty"`
	Name     HostName `json:",omitempty"`
	HostName HostName
	Type     RecordType `json:"Type"`
	TTL      Seconds    `json:"TTL"`
}

func (r *record) GetId() RecordId {
	return r.Id
}

func (r *record) GetHostName() HostName {
	return r.HostName
}

func (r *record) SetHostName(name HostName) {
	r.HostName = name
}

func (r *record) GetName() HostName {
	return r.Name
}

func (r *record) clearName() {
	r.Name = ""
}

func (r *record) GetType() RecordType {
	return r.Type
}

func (r *record) GetTTL() Seconds {
	return r.TTL
}

func (r *record) SetTTL(ttl Seconds) {
	r.TTL = ttl
}

type SOARecord struct {
	Mname            HostName    `json:"MNAME"`
	Rname            MailAddress `json:"RNAME"`
	Serial           uint
	Refresh          Seconds
	Retry            Seconds
	Expire           Seconds
	NegativeCacheTTL Seconds
	record
}

type NSRecord struct {
	NameServer HostName
	record
}

type ARecord struct {
	IPAddress IPv4 `json:"IPAddress"`
	record
}

type AAAARecord struct {
	IPAddress IPv6 `json:"IPAddress"`
	record
}

type CNAMERecord struct {
	AliasTo HostName `json:"AliasTo"`
	record
}

type MXRecord struct {
	MailServer HostName
	Priority   Priority
	record
}

type TXTRecord struct {
	Value string
	record
}

type SRVRecord struct {
	Priority Priority
	Weight   uint
	Port     uint
	Target   HostName
	record
}
