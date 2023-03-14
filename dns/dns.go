package dns

import (
	"reflect"
	"sync"
)

type DNS interface {
	GetListByDomain(domain string) (records []*DomainRecord, err error)
	UpsertDomain(domain, ip string) error
}

type DomainRecord struct {
	Prefix string
	Ip     string
	Id     string
}

var cache = make(map[any]any)
var l sync.Mutex

func Factory(dnsImp DNS) DNS {
	target := reflect.TypeOf(dnsImp)
	l.Lock()
	defer l.Unlock()
	find, ok := cache[target]
	findValue := reflect.ValueOf(find)
	if ok {
		typeOf := reflect.TypeOf(dnsImp)
		valueOf := reflect.ValueOf(dnsImp)
		if typeOf.Kind() == reflect.Pointer {
			typeOf = typeOf.Elem()
			valueOf = valueOf.Elem()
			findValue = findValue.Elem()
		}
		valueOf.Set(findValue)
		return dnsImp
	}
	cache[target] = dnsImp
	return dnsImp
}
