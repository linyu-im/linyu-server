package config

type StorageType string

const (
	LocalStorageType StorageType = "local"
	OssStorageType   StorageType = "oss"
)

type LocalStorageConfig struct {
	BaseURL  string `mapstructure:"base-url"`
	FilePath string `mapstructure:"file-path"`
}

type OssStorageConfig struct {
	BaseURL         string `mapstructure:"base-url"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access-key-id" `
	AccessKeySecret string `mapstructure:"access-key-secret"`
	BucketName      string `mapstructure:"bucket-name"`
}
