// Package cli
//
// @author: xwc1125
// @date: 2020/10/11
package cli

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Command *cobra.Command

// Cli cobra 在使用中，如果执行过中任何的Run或RunE没有执行过，那么cobra.OnInitialize(func func1)
// 中的func1 就不会执行
type Cli struct {
	rootCmd *cobra.Command
	viper   *viper.Viper

	subCmds []Command

	readConfigFunc func(viper *viper.Viper)
}

// NewCli 创建新的命令对象
func NewCli(a *AppInfo) *Cli {
	return NewCliWithViper(a, nil)
}

// NewCliWithViper 添加viper创建新的命令对象
func NewCliWithViper(a *AppInfo, _viper *viper.Viper) *Cli {
	if _viper == nil {
		_viper = viper.New()
	}
	rootCmd := &cobra.Command{
		Use:     a.App,
		Version: a.Version,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(a.Welcome)
		},
	}
	return &Cli{
		rootCmd: rootCmd,
		viper:   _viper,
		subCmds: make([]Command, 0),
	}
}

// InitFlags 初始化flag
// viper 按照如下顺序查找flag key:
// - pflag里面的被命令行显式设置的key
// - 环境变量显式设置的
// - 配置文件显式设置的
// - KV存储的
// - 通过viper设置的default flag
// - 如果前面都没有变化，最后使用pflag的默认值
//
// 所以在Unmarshal的时候命令行里面显式设置的flag会覆盖配置文件里面的flag
// 如果配置文件没有这个flag，会用pflag的默认值
// @params useDefaultFlags 是否使用默认的flag：config，env
// @params flagSetFunc flag设置的回调函数，函数中为viper及rootFlags
// @params readConfigFunc 读取配置的回调函数
func (cli *Cli) InitFlags(useDefaultFlags bool, flagSetFunc func(rootFlags *pflag.FlagSet), readConfigFunc func(viper *viper.Viper)) (err error) {
	var (
		configFile string
		configEnv  string
	)

	// 获取当前命令行
	{
		rootFlags := cli.rootCmd.PersistentFlags()
		if useDefaultFlags {
			rootFlags.StringVar(&configFile, "config", "./conf/config.yaml", "config file (default is ./conf/config.yaml)")
			rootFlags.StringVar(&configEnv, "env", "", "config env")
		}
		if flagSetFunc != nil {
			flagSetFunc(rootFlags)
		}
		// 将完整的命令绑定到viper上
		cli.viper.BindPFlags(rootFlags)
	}

	// 进行config初始化
	cobra.OnInitialize(func() {
		if useDefaultFlags {
			// 初始化配置文件
			if configFile != "" {
				cli.viper.SetConfigFile(configFile)
			} else {
				// 如果含有环境类型，那么使用config_{env}
				if configEnv != "" {
					cli.viper.SetConfigName("config_" + configEnv)
				} else {
					cli.viper.SetConfigName("config")
				}
				// 添加读取的配置文件路径
				cli.viper.AddConfigPath(".")
				cli.viper.AddConfigPath("./conf")
			}
			// viper加载配置
			if err = cli.viper.ReadInConfig(); err != nil {
				return
			}
		}
		cli.viper.AutomaticEnv()

		if readConfigFunc != nil {
			cli.readConfigFunc = readConfigFunc
			cli.readConfigFunc(cli.viper)
		}

		// 观察配置变更
		cli.watchConfig()
	})
	return nil
}

// 监听配置文件是否改变,用于热更新
func (cli *Cli) watchConfig() {
	cli.viper.WatchConfig()
	cli.viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
		if cli.readConfigFunc != nil {
			cli.readConfigFunc(cli.viper)
		}
	})
}

// RootCmd 获取cobra.Command
func (cli *Cli) RootCmd() *cobra.Command {
	return cli.rootCmd
}

// Viper 获取viper
func (cli *Cli) Viper() *viper.Viper {
	return cli.viper
}

// AddCommands 往cobra.Command添加多个命令
func (cli *Cli) AddCommands(cmds ...Command) {
	if cmds == nil {
		return
	}
	cli.subCmds = cmds
	for _, cmd := range cmds {
		cli.rootCmd.AddCommand(cmd)
	}
}

// GetCommands 获取所有的命令command
func (cli *Cli) GetCommands() []Command {
	return cli.subCmds
}

// Execute 执行命令
func (cli *Cli) Execute() error {
	return cli.rootCmd.Execute()
}

// Get 从配置文件中获取参数
func (cli *Cli) Get(key string) interface{} {
	return cli.viper.Get(key)
}

// GetString 获取字符串参数
func (cli *Cli) GetString(key string) string {
	return cli.viper.GetString(key)
}

// GetInt 获取int参数
func (cli *Cli) GetInt(key string) int {
	return cli.viper.GetInt(key)
}

// GetInt64 获取int64参数
func (cli *Cli) GetInt64(key string) int64 {
	return cli.viper.GetInt64(key)
}

// GetBool 获取bool参数
func (cli *Cli) GetBool(key string) bool {
	return cli.viper.GetBool(key)
}
