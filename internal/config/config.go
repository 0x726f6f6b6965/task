package config

type RedisCfg struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	MaxRetries int    `yaml:"max-retries" default:"3"`
	DB         int    `yaml:"db"`
}

type Log struct {
	// Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	Level            int    `default:"1" yaml:"level" help:"the application log level"`
	TimeFormat       string `default:"2006-01-02T15:04:05Z07:00" yaml:"time-format" help:"the application log time format"`
	TimestampEnabled bool   `yaml:"timestamp-enabled" default:"false"`
	ServiceName      string `yaml:"service-name" help:"the application service name"`
}

type Rest struct {
	Host string `yaml:"host" help:"the host to bind for REST server"`
	Port int    `yaml:"port" help:"the port to bind for REST server"`
}

type Config struct {
	Name   string   `yaml:"name" help:"the application name"`
	Rest   Rest     `yaml:"rest" help:"the application rest information"`
	Redis  RedisCfg `yaml:"redis" help:"the application redis option"`
	NodeID uint64   `yaml:"node-id"`
	Log    Log      `yaml:"log" help:"the application log"`
}
