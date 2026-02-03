package config

type EventBusType string

const (
	RocketEventBusType EventBusType = "rocket"
	LocalEventBusType  EventBusType = "local"
)

type RocketEventBusConfig struct {
	Namesrv    []string `mapstructure:"namesrv"`
	RetryTimes int      `mapstructure:"retry-times"`
}
