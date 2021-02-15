package logredis

import (
	"time"

	"github.com/sirupsen/logrus"
)

// LogMetaConfig stores Log message meta info
type LogMetaConfig struct {
	Channel       string
	Application   string
	Hostname      string
	Origin        string
	Source        string
	Extras        map[string]interface{}
	MessageFormat int
}

// Message Formats
const (
	V1 = iota
	V2
	AccessLog
	Custom
)

// EncodeV1 is used to assemble the message that storages meta information in @fields
func (l LogMetaConfig) EncodeV1(entry *logrus.Entry) map[string]interface{} {
	msg := make(map[string]interface{})

	msg["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	msg["@source_host"] = l.Hostname
	msg["@message"] = entry.Message

	fields := make(map[string]interface{})
	fields["level"] = entry.Level.String()
	fields["application"] = l.Application

	for k, v := range entry.Data {
		fields[k] = v
	}
	msg["@fields"] = fields

	return msg
}

// EncodeV2 is used to assemble the message that storages meta information flatly
func (l LogMetaConfig) EncodeV2(entry *logrus.Entry) map[string]interface{} {
	msg := make(map[string]interface{})

	msg["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	msg["host"] = l.Hostname
	msg["message"] = entry.Message
	msg["level"] = entry.Level.String()
	msg["application"] = l.Application
	for k, v := range entry.Data {
		msg[k] = v
	}

	return msg
}
