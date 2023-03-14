package dns

import (
	"errors"
	ali "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

var ErrKeyOrSecretNotSetting = errors.New("ErrKeyOrSecretNotSetting")

type AliDNS struct {
	AccessKey    string
	AccessSecret string
	client       *ali.Client
}

func (a *AliDNS) getClient() (client *ali.Client, err error) {
	if a.AccessKey == "" || a.AccessSecret == "" {
		return nil, ErrKeyOrSecretNotSetting
	}
	if a.client == nil {
		client, err = ali.NewClientWithAccessKey("cn-hangzhou", a.AccessKey, a.AccessSecret)
		a.client = client
		return
	}
	return a.client, nil
}

func (a *AliDNS) GetListByDomain(domain string) (records []*DomainRecord, err error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}
	request := ali.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = domain
	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return nil, err
	}
	if response.TotalCount > 0 {
		for _, record := range response.DomainRecords.Record {
			if record.Type != "A" {
				continue
			}
			records = append(records, &DomainRecord{
				Prefix: record.RR,
				Ip:     record.Value,
				Id:     record.RecordId,
			})
		}
	}
	return
}

type isUpdate = bool
type domainPrefix = string

func (a *AliDNS) UpsertDomain(domain, ip string) error {
	records, err := a.GetListByDomain(domain)
	if err != nil {
		return err
	}
	var prefixMap = make(map[domainPrefix]*DomainRecord)
	for _, record := range records {
		if record.Prefix == `@` || record.Prefix == `www` {
			prefixMap[record.Prefix] = record
		}
	}
	process := map[isUpdate]func(id, prefix, domain, ip string) error{
		true: func(id, prefix, domain, ip string) error {
			request := ali.CreateUpdateDomainRecordRequest()
			request.Scheme = "https"
			request.RecordId = id
			request.RR = prefix //one of = @ or www
			request.Type = "A"
			request.Value = ip
			_, err = a.client.UpdateDomainRecord(request)
			return err
		},
		false: func(id, prefix, domain, ip string) error {
			request := ali.CreateAddDomainRecordRequest()
			request.Scheme = "https"
			request.DomainName = domain
			request.RR = prefix //one of = @ or www
			request.Type = "A"
			request.Value = ip
			_, err = a.client.AddDomainRecord(request)
			return err
		},
	}
	for _, prefix := range []string{`@`, `www`} {
		if find, ok := prefixMap[prefix]; ok {
			if find.Ip != ip && ip != `` {
				if err := process[true](find.Id, find.Prefix, domain, ip); err != nil {
					return err
				}
			}
		} else {
			if err := process[false](find.Id, find.Prefix, domain, ip); err != nil {
				return err
			}
		}
	}
	return nil
}
