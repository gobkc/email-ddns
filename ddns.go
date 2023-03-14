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
	lastSendMail := time.Now()
	subject := fmt.Sprintf(emailReader.DomainFilter, conf.Domain)
	for {
		process := func(sleep time.Duration) {
			defer time.Sleep(sleep * time.Second)
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			if lastSendMail.Unix() <= time.Now().Add(10*time.Minute).Unix() {
				lastSendMail = time.Now()
				if err := mailSender.SendToMail(emailConf.User, subject, currentTime); err != nil {
					log.Printf("[ERR]\tSEND IP:%s", err.Error())
					return
				}
			}
			ip, err := emailReader.GetIp(emailConf.Imap, emailConf.User, emailConf.Pass, conf.Domain)
			if err != nil {
				log.Printf("[ERR]\tGET IP:\t%s", err.Error())
			}
			log.Printf("[INFO]\tGET IP:%s\t%s", ip, currentTime)
			if currentIP != ip || err != nil {
				currentIP = ip
				lastSendMail = time.Now()
				if err := mailSender.SendToMail(emailConf.User, subject, currentTime); err != nil {
					log.Printf("[ERR]\tSEND IP:%s", err.Error())
					return
				}
				log.Printf("[INFO]\tSEND IP:%s\t%s\tTO %s", ip, currentTime, subject)
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
