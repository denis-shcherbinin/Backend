package config

import (
	"github.com/spf13/viper"
	"os"
)

const (
	storageAccessKey = "ACCESS_KEY"
	storageSecretKey = "SECRET_KEY"
	storageEndpoint  = "storage.endpoint"
	storageRegion    = "storage.region"
	storageName      = "storage.name"
)

type StorageConfig struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Region    string
	Name      string
}

func storageConfigInit() *StorageConfig {
	return &StorageConfig{
		AccessKey: os.Getenv(storageAccessKey),
		SecretKey: os.Getenv(storageSecretKey),
		Endpoint:  viper.GetString(storageEndpoint),
		Region:    viper.GetString(storageRegion),
		Name:      viper.GetString(storageName),
	}
}
