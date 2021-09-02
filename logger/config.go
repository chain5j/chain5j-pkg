// Package logger
//
// @author: xwc1125
// @date: 2021/7/12
package logger

import (
	"os"
	"path/filepath"
	"strings"
)

// LogConfig log config
type LogConfig struct {
	Console ConsoleLogConfig `json:"console" mapstructure:"console"`
	File    FileLogConfig    `json:"file" mapstructure:"file"`
}

type ConsoleLogConfig struct {
	Level    int    `json:"level" mapstructure:"level"`         // log level
	Modules  string `json:"modules" mapstructure:"modules"`     // need to show modules。"*":all
	ShowPath bool   `json:"show_path" mapstructure:"show_path"` // show path
	Format   string `json:"format" mapstructure:"format"`       // format

	UseColor bool `json:"use_color" mapstructure:"use_color"` // show console color
	Console  bool `json:"console" mapstructure:"console"`     // show console
}

type FileLogConfig struct {
	Level   int    `json:"level" mapstructure:"level"`     // log level
	Modules string `json:"modules" mapstructure:"modules"` // need to show modules。"*":all
	Format  string `json:"format" mapstructure:"format"`   // format

	Save     bool   `json:"save" mapstructure:"save"`           // whether save file
	FilePath string `json:"file_path" mapstructure:"file_path"` // filepath
	FileName string `json:"file_name" mapstructure:"file_name"` // filename prefix
}

func (c *ConsoleLogConfig) GetModules() []string {
	if c.Modules == "" {
		return []string{"*"}
	}
	return strings.Split(c.Modules, ",")
}

func (c *FileLogConfig) GetModules() []string {
	if c.Modules == "" {
		return []string{"*"}
	}
	return strings.Split(c.Modules, ",")
}

func (c *FileLogConfig) GetLogFile() string {
	if !c.Save {
		return ""
	}
	file := c.FileName
	if file == "" {
		file = "log.json"
	}
	if c.FilePath != "" {
		os.MkdirAll(c.FilePath, os.ModePerm)
		file = filepath.Join(c.FilePath, file)
	} else {
		file = filepath.Join(".", file)
	}
	return file
}

func LvlFromInt(lvlInt int) string {
	switch lvlInt {
	case 0:
		return "crit"
	case 1:
		return "error"
	case 2:
		return "warn"
	case 3:
		return "info"
	case 4:
		return "debug"
	case 5:
		return "trace"
	default:
		return "info"
	}
}
