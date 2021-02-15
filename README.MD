# Redis Hook for logrus

This library provides a logrus hook base on Redis PUB/SUB function. Useful for centralized logging, espaecially for RELK Stack (Redis, Elasticsearch, Logstash, Kibana). 

When this hook instance was added to logrus, all log message will be sent to redis by PUBLISH command in JSON string form.

## Install

```bash
$ go get github.com/vixalie/logrus-redis
```

## Usage

