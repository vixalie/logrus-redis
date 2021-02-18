# Redis Hook for logrus

This library provides a logrus hook base on both Redis PUB/SUB function and LIST. Useful for centralized logging, espaecially for RELK Stack (Redis, Elasticsearch, Logstash, Kibana). 

When this hook instance was added to logrus, all log message will be sent to redis by PUBLISH command in JSON string form, and by RPUSH command in list mode.

## Install

```bash
$ go get github.com/vixalie/logrus-redis
```

## Usage

To build a hook, there are two configurations needed. One is for connecting to Redis service.

```go
// This library was imported as "logredis"
hookRedisConfig := logredis.RedisConfig{
    Host:     "127.0.0.1",
    Port:     6379,
    DB:       0,
    Password: "", // A empty string means there is no password.
}
```

The another is for construct log message sent to Redis service. We call it log meta configuration.

```go
hookMetaConfig := logredis.LogMetaConfig{
    Channel:         "log_topic",
    Application:     "test_app",
    Hostname:        "library.archknight.xyz",
    TrailingNewLine: false, // set true to add a trailing new line character '\n'.
    ListMode:        false, // set true to use RPUSH command instead of PUBLISH.
    MessageFormat:   logredis.V1, // Use logredis.Custom when you need custom formatter.
    Source:          "test_demo",
}
```

There are three message formatters can be used directly.

- `logredis.V1`, stores log message in `@message` key, and extra information in `@field` key.
- `logredis.V2`, stores every information flatly.
- `logredis.AccessLog`, stores log message in `message` key, and extra information in `@field` key. This format does NOT store time.

There is also a fourth message formatter (`logredis.Custom`) can be used, but you need to define the formatter you need in `HookConfig.Formatter` field. This field accept a function signed with `func(*logrus.Entry, *LogMetaConfig) map[string]interface{}`. The return value of that formatter will be encoded in JSON and sent to Redis.

When everything is completed, there will be a hook.

```go
logger := logrus.New()
hookConfig := logredis.HookConfig{
    Redis: hookRedisConfig,
    Meta: hookMetaConfig,
    // Formatter: customFormatter,  // When you need a custom formatter, set it.
}
hook := logredis.NewHook(hookConfig)
logger.Hooks.Add(hook)
```