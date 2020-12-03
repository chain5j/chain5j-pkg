// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"crypto/tls"
	"crypto/x509"
	log "github.com/chain5j/log15"
	"io/ioutil"
	"net"
)

// StartHTTPEndpoint starts the HTTP RPC endpoint, configured with cors/vhosts/modules
func StartHTTPEndpoint(httpConfig HttpConfig, tlsConfig TlsConfig, apis []API) (net.Listener, *Server, error) {
	log := log.New("rpc_http")
	// Generate the whitelist based on the allowed modules
	whitelist := make(map[string]bool)
	for _, module := range httpConfig.Modules {
		whitelist[module] = true
	}
	// Register all the APIs exposed by the services
	handler := NewServer()
	for _, api := range apis {
		if whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
			if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
				return nil, nil, err
			}
			log.Debug("HTTP registered", "namespace", api.Namespace)
		}
	}
	// All APIs registered, start the HTTP listener
	var (
		listener net.Listener
		err      error
	)

	endpoint := httpConfig.Endpoint()
	if listener, err = net.Listen("tcp", endpoint); err != nil {
		return nil, nil, err
	}

	httpServer := NewHTTPServer(httpConfig.Cors, httpConfig.VirtualHosts, httpConfig.Timeouts, handler)
	if tlsConfig.Mod == Disable {
		go httpServer.Serve(listener)
	} else {
		tlsConf, err := getTls(true, tlsConfig)
		if err != nil {
			return nil, nil, err
		}
		httpServer.TLSConfig = tlsConf
		go httpServer.ServeTLS(listener, tlsConfig.CrtFile, tlsConfig.PrvkeyFile)
	}

	return listener, handler, err
}

func getTls(isServer bool, tlsConfig TlsConfig) (*tls.Config, error) {
	switch tlsConfig.Mod {
	case OneWay:
		if isServer {
			// 添加证书
			cert, err := tls.LoadX509KeyPair(tlsConfig.CrtFile, tlsConfig.PrvkeyFile)
			if err != nil {
				log.Error("tls.LoadX509KeyPair err", "err", err)
				return nil, err
			}
			return &tls.Config{
				Certificates: []tls.Certificate{cert},
			}, err
		} else {
			conf := &tls.Config{
				InsecureSkipVerify: true, // 用来控制客户端是否证书和服务器主机名。如果设置为true,则不会校验证书以及证书中的主机名和服务器主机名是否一致。
			}
			return conf, nil
		}
	case TwoWay:
		cert, err := tls.LoadX509KeyPair(tlsConfig.CrtFile, tlsConfig.PrvkeyFile)
		if err != nil {
			log.Error("tls.LoadX509KeyPair err", "err", err)
			return nil, err
		}

		// CA证书池
		certPool := x509.NewCertPool()
		for _, root := range tlsConfig.CaRoots {
			certBytes, err := ioutil.ReadFile(root)
			if err != nil {
				log.Error("unable to read ca.pem", "err", err)
				return nil, err
			}
			ok := certPool.AppendCertsFromPEM(certBytes)
			if !ok {
				log.Error("failed to parse root certificate", "err", err)
				return nil, err
			}
		}
		if isServer {
			return &tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    certPool,
				MinVersion:   tls.VersionTLS12,
			}, nil
		} else {
			return &tls.Config{
				InsecureSkipVerify: true, // 用来控制客户端是否证书和服务器主机名。如果设置为true,则不会校验证书以及证书中的主机名和服务器主机名是否一致。
				Certificates:       []tls.Certificate{cert},
				ClientCAs:          certPool,
			}, nil
		}
	case Disable:
	}
	return nil, nil
}

// StartWSEndpoint starts a websocket endpoint
func StartWSEndpoint(wsConfig WSConfig, tlsConfig TlsConfig, apis []API, ) (net.Listener, *Server, error) {
	log := log.New("ws")
	// Generate the whitelist based on the allowed modules
	whitelist := make(map[string]bool)
	for _, module := range wsConfig.Modules {
		whitelist[module] = true
	}
	// Register all the APIs exposed by the services
	handler := NewServer()
	for _, api := range apis {
		if wsConfig.ExposeAll || whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
			if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
				return nil, nil, err
			}
			log.Debug("WebSocket registered", "service", api.Service, "namespace", api.Namespace)
		}
	}
	// All APIs registered, start the HTTP listener
	var (
		listener net.Listener
		err      error
	)
	if listener, err = net.Listen("tcp", wsConfig.Endpoint()); err != nil {
		return nil, nil, err
	}

	wsServer := NewWSServer(wsConfig.Origins, handler)

	if tlsConfig.Mod == Disable {
		go wsServer.Serve(listener)
	} else {
		tlsConf, err := getTls(true, tlsConfig)
		if err != nil {
			return nil, nil, err
		}
		wsServer.TLSConfig = tlsConf
		go wsServer.ServeTLS(listener, tlsConfig.CrtFile, tlsConfig.PrvkeyFile)
	}
	return listener, handler, err

}

// StartIPCEndpoint starts an IPC endpoint.
func StartIPCEndpoint(ipcEndpoint string, apis []API) (net.Listener, *Server, error) {
	log := log.New("ipc")
	// Register all the APIs exposed by the services.
	handler := NewServer()
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			return nil, nil, err
		}
		log.Debug("IPC registered", "namespace", api.Namespace)
	}
	// All APIs registered, start the IPC listener.
	listener, err := ipcListen(ipcEndpoint)
	if err != nil {
		return nil, nil, err
	}
	go handler.ServeListener(listener)
	return listener, handler, nil
}
