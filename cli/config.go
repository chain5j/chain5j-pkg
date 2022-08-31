// Package cli cli命令
//
// @author: xwc1125
// @date: 2020/10/11
package cli

import (
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// AppInfo 应用信息
type AppInfo struct {
	App     string `json:"app"`
	Version string `json:"version"`
	Welcome string `json:"welcome"`
}

var (
	eventMapLock sync.Mutex
	eventMap     = make(map[string][]func(e fsnotify.Event), 0)
)

// LoadConfig 加载本地配置文件
func LoadConfig(configFile string, key string, result interface{}) error {
	viper := viper.New()
	// 初始化配置文件
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		// 添加读取的配置文件路径
		viper.AddConfigPath(".")
		viper.AddConfigPath("./conf")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if len(key) == 0 {
		if err := viper.Unmarshal(&result); err != nil {
			return err
		}
	} else {
		if err := viper.UnmarshalKey(key, &result); err != nil {
			return err
		}
	}

	return nil
}
func LoadConfigWithEvent(configFile string, key string, result interface{}, event func(e fsnotify.Event)) (err error) {
	// 初始化配置文件
	if configFile != "" {
		configFile, err = filepath.Abs(configFile)
		if err != nil {
			return err
		}
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		// 添加读取的配置文件路径
		viper.AddConfigPath(".")
		viper.AddConfigPath("./conf")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if len(key) == 0 {
		if err := viper.Unmarshal(&result); err != nil {
			return err
		}
	} else {
		if err := viper.UnmarshalKey(key, &result); err != nil {
			return err
		}
	}
	eventMapLock.Lock()
	configFile = filepath.Clean(configFile)
	if _, ok := eventMap[configFile]; !ok {
		viper.WatchConfig()
	}
	eventMapLock.Unlock()
	if event != nil {
		eventMapLock.Lock()
		if e, ok := eventMap[configFile]; ok {
			e = append(e, event)
			eventMap[configFile] = e
		} else {
			eventMap[configFile] = append(eventMap[configFile], event)
		}
		eventMapLock.Unlock()
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		eventMapLock.Lock()
		if events, ok := eventMap[e.Name]; ok {
			for _, f := range events {
				f(e)
			}
		}
		eventMapLock.Unlock()
	})

	return nil
}
