// Package network
//
// @author: xwc1125
// @date: 2021/7/27
package network

type TlsMod string

var (
	Disable TlsMod = "disable" // 不启用TLS
	OneWay  TlsMod = "oneway"  // 单向认证
	TwoWay  TlsMod = "twoway"  // 双向认证
)

type TlsConfig struct {
	Mod                TlsMod   `json:"mod" mapstructure:"mod"`                                   // 模式
	KeyFile            string   `json:"key_file" mapstructure:"key_file"`                         // 服务私钥
	CertFile           string   `json:"cert_file" mapstructure:"cert_file"`                       // 服务证书
	CaRoots            []string `json:"ca_roots" mapstructure:"ca_roots"`                         // ca根证书
	CaRootPaths        []string `json:"ca_root_paths" mapstructure:"ca_root_paths"`               // ca根证书路径
	ServerName         string   `json:"server_name" mapstructure:"server_name"`                   // 服务name
	InsecureSkipVerify bool     `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify"` // 是否跳过非安全的校验
}
