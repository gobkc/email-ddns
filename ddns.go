package main

import (
	"encoding/json"
	"fmt"
	"github.com/gobkc/email-ddns/configurator"
	"github.com/gobkc/email-ddns/dns"
	emailReader "github.com/gobkc/email-ddns/email-reader"
	"github.com/gobkc/email-ddns/internal"
	"github.com/gobkc/mails"
	"log"
	"net"
	"time"
)

func main() {
	conf := internal.Conf{}
	confJson := configurator.Factory[configurator.Json]()
	confJson.UnMarshal("./config.json", &conf)
	var currentIP string

	log.Printf("[INFO]\tStart Service")
	dnsUtil := dns.Factory(&dns.AliDNS{
		AccessKey:    conf.Dns.Key,
		AccessSecret: conf.Dns.Secret,
	})
	emailConf := conf.Email
	mailSender := mails.Factory(&mails.DefaultEmail{})
	if err := mailSender.SetEnv(emailConf.User, emailConf.Pass, emailConf.Pop3); err != nil {
		log.Printf("[ERR]\tSetEnv:\t%s", err.Error())
		return
	}
	data, _ := json.MarshalIndent(conf, "", "\t")
	log.Printf("[INFO]\tRead Configure:\n%s\n", string(data))
	subject := fmt.Sprintf(emailReader.DomainFilter, conf.Domain)
	waitUpsert := func(oldIp string) {
		log.Printf("[INFO]\tNETWORK CHANGED WAITING...\n")
		for {
			_, err := net.DialTimeout("ip4:icmp", conf.Ping, time.Duration(1000*1000*1000))
			if err != nil {
				log.Printf("[INFO]\tNETWORK FAILED,SLEEP 1 SECOND\n")
				time.Sleep(1 * time.Second)
				continue
			}
			log.Printf("[INFO]\tNETWORK CONNECT SUCCESS!\n")
			now := time.Now().Format("2006-01-02 15:04:05")
			if err := mailSender.SendToMail(emailConf.User, subject, now); err != nil {
				log.Printf("[ERR]\tSEND IP:%s", err.Error())
				time.Sleep(1 * time.Second)
				continue
			}
			log.Printf("[INFO]\tSEND IP\n")
			ip, err := emailReader.GetIp(emailConf.Imap, emailConf.User, emailConf.Pass, conf.Domain)
			if err != nil {
				log.Printf("[ERR]\tGET IP:\t%s", err.Error())
				time.Sleep(1 * time.Second)
				continue
			}
			if ip == oldIp {
				log.Printf("[ERR]\tOLD IP:%s\tNEW IP:%s\n", oldIp, ip)
				for i := 0; i < 10; i++ {
					ip, _ = emailReader.GetIp(emailConf.Imap, emailConf.User, emailConf.Pass, conf.Domain)
					if ip == oldIp {
						log.Printf("[ERR]\tOLD IP:%s\tNEW IP:%s\n", oldIp, ip)
						time.Sleep(15 * time.Second)
						continue
					}
				}
			}
			if err := dnsUtil.UpsertDomain(conf.Domain, ip); err != nil {
				log.Printf("[ERR]\tUPSERT DOMAIN:%s", err.Error())
			}
			break
		}
	}
	for {
		process := func(sleep time.Duration) {
			defer time.Sleep(sleep * time.Second)
			_, err := net.DialTimeout("ip4:icmp", conf.Ping, time.Duration(1000*1000*1000))
			if err != nil {
				waitUpsert(currentIP)
			}
			ip, err := emailReader.GetIp(emailConf.Imap, emailConf.User, emailConf.Pass, conf.Domain)
			if err != nil {
				log.Printf("[ERR]\tGET IP:\t%s", err.Error())
			}
			log.Printf("[INFO]\tGET IP:%s\n", ip)
			if currentIP != ip || err != nil {
				currentIP = ip
				if err := dnsUtil.UpsertDomain(conf.Domain, currentIP); err != nil {
					log.Printf("[ERR]\tUPSERT DOMAIN:%s", err.Error())
				}
			}
		}
		if conf.Interval == 0 {
			conf.Interval = 15
		}
		process(time.Duration(conf.Interval))
	}
}
