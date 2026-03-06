package config

type EmbeddingConfig struct {
	BaseUrl    string `mapstructure:"base-url"`
	ApiKey     string `mapstructure:"api-key"`
	Model      string `mapstructure:"model"`
	Dimensions int    `mapstructure:"dimensions"`
}
