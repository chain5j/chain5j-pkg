// description: chain5j-pkg
// 
// @author: xwc1125
// @date: 2020/10/11
package ioutil

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// 路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 获取项目的路径
func GetProjectPath() (dir string, err error) {
	return os.Getwd()
}

// DefaultDataDir
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Pangu")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Pangu")
		} else {
			return filepath.Join(home, ".pangu")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

// homeDir Home路径
func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func MakeDirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
