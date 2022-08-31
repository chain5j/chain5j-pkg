// Package rpc
//
// @author: xwc1125
package rpc

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chain5j/chain5j-pkg/util/ioutil"
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
