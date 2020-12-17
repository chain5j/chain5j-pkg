// description: chain5j 
// 
// @author: xwc1125
// @date: 2020/11/18
package rpc

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/util/ioutil"
	log "github.com/chain5j/log15"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	logger = log.New("rpc")
)

type HttpConfig struct {
	Host         string       `json:"host" mapstructure:"host"`                   // ip
	Port         int          `json:"port" mapstructure:"port"`                   // 端口号
	Cors         []string     `json:"cors" mapstructure:"cors"`                   // 跨域
	VirtualHosts []string     `json:"virtual_hosts" mapstructure:"virtual_hosts"` // 虚拟host
	Modules      []string     `json:"modules" mapstructure:"modules"`             // 需要加载的模块
	Timeouts     HTTPTimeouts `json:"timeouts" mapstructure:"timeouts"`           // 超时配置
}

func (c HttpConfig) Endpoint() string {
	if c.Host == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type TlsMod string

var (
	Disable TlsMod = "disable" //不启用TLS
	OneWay  TlsMod = "oneway"  //单向认证
	TwoWay  TlsMod = "twoway"  //双向认证
)

type TlsConfig struct {
	Mod        TlsMod   `json:"mod" mapstructure:"mod"`                 // 模式
	PrvkeyFile string   `json:"prvkey_file" mapstructure:"prvkey_file"` // 服务私钥
	CrtFile    string   `json:"crt_file" mapstructure:"crt_file"`       // 服务证书
	CaRoots    []string `json:"ca_roots" mapstructure:"ca_roots"`       // ca根证书
}

type WSConfig struct {
	Host      string   `json:"host" mapstructure:"host"` // ip
	Port      int      `json:"port" mapstructure:"port"` // 端口号
	Origins   []string `json:"origins" mapstructure:"origins"`
	Modules   []string `json:"modules" mapstructure:"modules"` // 需要加载的模块
	ExposeAll bool     `json:"expose_all" mapstructure:"expose_all"`
}

func (c WSConfig) Endpoint() string {
	if c.Host == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type IpcConfig struct {
	Path string `json:"path" mapstructure:"path"`
}

func (c IpcConfig) Endpoint() string {
	// Short circuit if IPC has not been enabled
	if c.Path == "" {
		return ""
	}
	// On windows we can only use plain top-level pipes
	if runtime.GOOS == "windows" {
		if strings.HasPrefix(c.Path, `\\.\pipe\`) {
			return c.Path
		}
		return `\\.\pipe\` + c.Path
	}

	dir, _ := ioutil.GetProjectPath()
	if filepath.Base(c.Path) == c.Path {
		return filepath.Join(dir, c.Path)
	}
	return c.Path
}
