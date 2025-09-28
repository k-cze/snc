package config

type Config struct {
	Source        string
	Target        string
	DeleteMissing bool
	LogLevel      string
	UpdateMethod  string
}

type ConfigProvider interface {
	Config() *Config
}
