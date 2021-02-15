package logredis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// RedisConfig stores Redis configuration need to setup Hook
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// HookConfig stores all Logrus Redis Hook needs
type HookConfig struct {
	Redis     RedisConfig
	Meta      LogMetaConfig
	Formatter func(*logrus.Entry, *LogMetaConfig) map[string]interface{}
}

// RedisHook to sends logs to Redis server
type RedisHook struct {
	ConnectionPool *redis.Client
	Config         *HookConfig
}

// NewHook creates a hook to be attached to logrus logger
func NewHook(config HookConfig) *RedisHook {
	rdsConnection := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	return &RedisHook{
		ConnectionPool: rdsConnection,
		Config:         &config,
	}
}

// Levels returns the available logging levels.
func (hook *RedisHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire is called when a log event is fired.
func (hook *RedisHook) Fire(entry *logrus.Entry) error {
	var msg interface{}

	switch hook.Config.Meta.MessageFormat {
	case V1:
		msg = hook.Config.Meta.EncodeV1(entry)
	case V2:
		msg = hook.Config.Meta.EncodeV2(entry)
	case AccessLog:

	case Custom:
		msg = hook.Config.Formatter(entry, &hook.Config.Meta)
	default:
		return errors.New("Invalid message formatter")
	}

	msgInJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to create message for REDIS: %s", err)
	}

	execution := hook.ConnectionPool.Publish(entry.Context, hook.Config.Meta.Channel, string(msgInJSON))
	if err := execution.Err(); err != nil {
		return fmt.Errorf("Publish message failed, %v", err)
	}

	return nil
}