// Package cli
//
// @author: xwc1125
package cli

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// AppInfo app Info
type AppInfo struct {
	App     string `json:"app"`
	Version string `json:"version"`
	Welcome string `json:"welcome"`
}

// LoadConfig load local profile
// @params configFile config file pat
// @params key if key=="",load all the config;else load the key config
// @params result the result must be ptr
// @params event the file change event
func LoadConfig(configFile string, key string, result interface{}, event func(e fsnotify.Event)) error {
	return LoadConfigWithViper(viper.GetViper(), configFile, key, result, event)
}

// LoadConfigWithViper load config with viper
func LoadConfigWithViper(_viper *viper.Viper, configFile string, key string, result interface{}, event func(e fsnotify.Event)) error {
	// init config
	if configFile != "" {
		_viper.SetConfigFile(configFile)
	} else {
		_viper.SetConfigName("config")
		// add config path
		_viper.AddConfigPath(".")
		_viper.AddConfigPath("./conf")
	}
	_viper.AutomaticEnv()

	if err := _viper.ReadInConfig(); err != nil {
		return err
	}
	if len(key) == 0 {
		if err := _viper.Unmarshal(&result); err != nil {
			return err
		}
	} else {
		if err := _viper.UnmarshalKey(key, &result); err != nil {
			return err
		}
	}

	_viper.WatchConfig()
	_viper.OnConfigChange(event)

	return nil
}
