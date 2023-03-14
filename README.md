# About Email DDNS
Use email for dynamic IP auto-writing to domain name A records

### Note
Currently only aliyun DDNS are supported

### How to get it?
````
go clone github.com/gobkc/email-ddns.git
go build
````

### Configuration
````
{
  "domain": "kafs.cn",
  "interval": 15,
  "email": {
    "imap": "imap.163.com:993",
    "pop3": "smtp.163.com:25",
    "user": "xxx@163.com",
    "pass": "your email password"
  },
  "dns": {
    "key": "your ali ddns key",
    "secret": "your ali ddns secret"
  }
}
````

### License
© Gobkc, 2023~time.Now

Released under the Apache License