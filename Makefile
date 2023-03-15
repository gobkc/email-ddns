build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"
	upx email-ddns
	docker build -t registry.cn-shenzhen.aliyuncs.com/core-image/email-ddns .
	sudo docker push registry.cn-shenzhen.aliyuncs.com/core-image/email-ddns:latest

