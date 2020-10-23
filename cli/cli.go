package cli

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/util/ioutil"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	// 命令行
	RootCli *Cli
)

// CommandFunc 代表了一个子命令，用于往Cli注册子命令
type CommandFunc func(c *Cli) *cobra.Command

var (
	// commands 用于收集所有的子命令，在启动的时候统一往Cli注册
	commands []CommandFunc
)

// cobra 在使用中，如果执行过中任何的Run或RunE没有执行过，那么cobra.OnInitialize(func func1)
// 中的func1 就不会执行
type Cli struct {
	rootCmd *cobra.Command
	viper   *viper.Viper

	DataDir        string
	readConfigFunc func(viper *viper.Viper)
}

// 创建新的命令对象
func NewCli(a *AppInfo) *Cli {
	rootCmd := &cobra.Command{
		Use:     a.App,
		Version: a.Version,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(a.Welcome)
		},
	}
	RootCli = &Cli{
		rootCmd: rootCmd,
	}

	return RootCli
}

// 初始化
func (cli *Cli) Init() error {
	err := cli.initFlags()
	if err != nil {
		return err
	}
	return nil
}

// 初始化flag
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
func (cli *Cli) initFlags() error {
	var configFile string
	//var configEnv string
	// 获取当前命令行
	rootFlags := cli.rootCmd.PersistentFlags()

	// 设置数据目录
	rootFlags.StringVar(&cli.DataDir, "datadir", ioutil.DefaultDataDir(), "Data directory for the databases and keystore")
	// host
	rootFlags.StringP("host", "H", "127.0.0.1:9545", "server node ip:port")
	rootFlags.Int32P("rpcport", "p", 9545, "rpc port(default is 9545)")
	// configFile返回的值，config:参数名，value:默认值，usage:用法说明
	rootFlags.StringVar(&configFile, "config", "./conf/config.yaml", "config file (default is ./conf/config.yaml)")
	//rootFlags.StringVar(&configEnv, "env", "dev", "config env(default is dev)")

	// 将完整的命令绑定到viper上
	viper.BindPFlags(rootFlags)

	cobra.OnInitialize(func() {
		// 初始化配置文件
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			//viper.SetConfigName("config_" + configEnv)
			viper.SetConfigName("config")
			// 添加读取的配置文件路径
			viper.AddConfigPath(".")
			viper.AddConfigPath("./conf")
		}
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Print(err)
			os.Exit(1)
		}
		if cli.readConfigFunc != nil {
			cli.readConfigFunc(viper.GetViper())
		}

		cli.viper = viper.GetViper()
		//cli.watchConfig()
	})

	return nil
}

// 监听配置文件是否改变,用于热更新
func (cli *Cli) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
		if cli.readConfigFunc != nil {
			cli.readConfigFunc(cli.viper)
		}
	})
}

func (cli *Cli) GetCommands() []CommandFunc {
	return commands
}
func (cli *Cli) GetViper() *viper.Viper {
	return cli.viper
}

// 添加所有的命令
func (cli *Cli) AddCommands(cmds []CommandFunc) {
	for _, cmd := range cmds {
		cli.rootCmd.AddCommand(cmd(cli))
	}
}

// 添加子命令
func AddCommand(cmd CommandFunc) {
	commands = append(commands, cmd)
}

// 执行命令
func (cli *Cli) Execute(readConfigFunc func(viper *viper.Viper)) {
	cli.readConfigFunc = readConfigFunc
	err := cli.rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 从配置文件中获取参数
func (cli *Cli) Get(key string) interface{} {
	return cli.viper.Get(key)
}

func (cli *Cli) GetString(key string) string {
	return cli.viper.GetString(key)
}

func (cli *Cli) GetInt(key string) int {
	return cli.viper.GetInt(key)
}

func (cli *Cli) GetInt64(key string) int64 {
	return cli.viper.GetInt64(key)
}

func (cli *Cli) GetBool(key string) bool {
	return cli.viper.GetBool(key)
}
