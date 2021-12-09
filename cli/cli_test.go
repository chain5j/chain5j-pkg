// Package cli
//
// @author: xwc1125
package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestNewCli(t *testing.T) {
	rootCli := NewCli(&AppInfo{
		App:     "cli",
		Version: "v0.0.1",
		Welcome: "Helle world",
	})
	err := rootCli.InitFlags(true, func(flags *pflag.FlagSet) {}, func(viper *viper.Viper) {
		// 这一步时,rootFlags已经获取到了值
		// 当config不为空时,才启用
	})
	if err != nil {
		t.Log("initCli", "err", err)
		return
	}

	//rootCli.AddCommand(GenesisCommand)

	rootCli.Execute()
}

func TestNewCliNoDefault(t *testing.T) {
	rootCli := NewCli(&AppInfo{
		App:     "cli",
		Version: "v0.0.1",
		Welcome: "Helle world",
	})
	err := rootCli.InitFlags(true, func(flags *pflag.FlagSet) {
		var configFile string
		flags.StringVar(&configFile, "config", "./conf/config.yaml", "config file (default is ./conf/config.yaml)")
		// 初始化配置文件
		if configFile != "" {
			rootCli.Viper().SetConfigFile(configFile)
		} else {
			rootCli.Viper().SetConfigName("config")
			// 添加读取的配置文件路径
			rootCli.Viper().AddConfigPath(".")
			rootCli.Viper().AddConfigPath("./conf")
		}
		if err := rootCli.Viper().ReadInConfig(); err != nil {
			t.Log("viper.ReadInConfig err", err)
			os.Exit(1)
		}
	}, func(viper *viper.Viper) {
		// 这一步时,rootFlags已经获取到了值
		// 当config不为空时,才启用
	})
	if err != nil {
		t.Log("initCli", "err", err)
		return
	}

	//rootCli.AddCommand(GenesisCommand)

	rootCli.Execute()
}
