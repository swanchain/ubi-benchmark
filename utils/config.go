package utils

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

var config *Config

type Config struct {
	MCS MCS
	HUB HUB
}

type MCS struct {
	ApiKey     string
	BucketName string
	Network    string
}

type HUB struct {
	HubUrl           string `toml:"HUB_URL"`
	TaskUrl          string `toml:"TASK_URL"`
	CheckInterval    int64  `toml:"CHECK_INTERVAL"`
	BatchNum         int    `toml:"BATCH_NUM"`
	ENABLE_TITAN     int    `toml:"ENABLE_TITAN"`
	TITAN_KEY        string `toml:"TITAN_KEY"`
	TITAN_FOLDER_512 int    `toml:"TITAN_FOLDER_512"`
	TITAN_FOLDER_32  int    `toml:"TITAN_FOLDER_32"`
}

func InitConfig() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	configFile := filepath.Join(dir, "config.toml")

	if metaData, err := toml.DecodeFile(configFile, &config); err != nil {
		return fmt.Errorf("failed load config file, path: %s, error: %w", configFile, err)
	} else {
		if !requiredFieldsAreGiven(metaData) {
			log.Fatal("Required fields not given")
		}
	}
	return nil
}

func GetConfig() *Config {
	return config
}

func requiredFieldsAreGiven(metaData toml.MetaData) bool {
	requiredFields := [][]string{
		{"MCS"},
		{"MCS", "ApiKey"},
		{"MCS", "BucketName"},
		{"MCS", "Network"},
	}

	for _, v := range requiredFields {
		if !metaData.IsDefined(v...) {
			log.Fatal("Required fields ", v)
		}
	}

	return true
}
