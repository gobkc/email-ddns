FROM alpine:3.14

RUN mkdir /app

WORKDIR /app

ADD email-ddns /app

ADD config.json /app

# run service
CMD ["./email-ddns"]
