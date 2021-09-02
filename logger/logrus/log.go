// Package logrus
//
// @author: xwc1125
// @date: 2021/8/30
package logrus

import (
	"github.com/chain5j/chain5j-pkg/logger"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"time"
)

var (
	_    logger.Logger = new(log)
	root               = logrus.New()
)

type log struct {
	*logrus.Logger
	entry *logrus.Entry
}

func InitWithConfig(config *logger.LogConfig) logger.Logger {
	root.SetReportCaller(config.Console.ShowPath)
	root.SetLevel(logrus.Level(config.Console.Level))

	root.SetFormatter(&Formatter{
		UseColor:        config.Console.UseColor,
		LocationEnabled: config.Console.ShowPath,
		Modules:         strings.Split(config.Console.Modules, ","),
	})

	if config.File.Save {
		// WithMaxAge和WithRotationCount 二者只能设置一个，
		// WithMaxAge 设置文件清理前的最长保存时间，
		// WithRotationCount 设置文件清理前最多保存的个数。
		writer, _ := rotatelogs.New(
			config.File.GetLogFile()+".%Y%m%d%H",
			rotatelogs.WithLinkName(config.File.GetLogFile()), // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(time.Hour*24*30),            // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Hour),            // 日志切割时间间隔
		)
		var writerMap = make(map[logrus.Level]io.Writer, 0)
		if config.File.Level >= 5 {
			writerMap[logrus.TraceLevel] = writer
		}
		if config.File.Level >= 4 {
			writerMap[logrus.DebugLevel] = writer
		}
		if config.File.Level >= 3 {
			writerMap[logrus.InfoLevel] = writer
		}
		if config.File.Level >= 2 {
			writerMap[logrus.WarnLevel] = writer
		}
		if config.File.Level >= 1 {
			writerMap[logrus.ErrorLevel] = writer
		}
		if config.File.Level >= 0 {
			writerMap[logrus.FatalLevel] = writer
			writerMap[logrus.PanicLevel] = writer
		}
		root.Hooks.Add(lfshook.NewHook(
			lfshook.WriterMap(writerMap),
			&logrus.JSONFormatter{},
		))
	} else if config.Console.Console {
		root.SetOutput(os.Stdout)
	}

	return &log{
		Logger: root,
	}
}

func (l *log) New(module string, ctx ...interface{}) logger.Logger {
	return New(module, ctx)
}

func New(module string, ctx ...interface{}) logger.Logger {
	var fields = make(map[string]interface{}, 0)
	fields["module"] = module
	if ctx != nil && len(ctx) > 0 {
		if len(ctx)%2 != 0 {
			ctx = append(ctx, "NULL")
		}
	}
	for i := 0; i < len(ctx); i = i + 2 {
		fields[ToString(ctx[i])] = ctx[i+1]
	}
	return &log{
		Logger: root,
		entry:  logrus.WithField("module", module),
	}
}

func (l *log) Trace(msg string, ctx ...interface{}) {
	fields := make(map[string]interface{})
	for i := 0; i < len(ctx); i = i + 2 {
		if i+2 <= len(ctx) {
			fields[ToString(ctx[i])] = ctx[i+1]
		} else {
			fields[ToString(ctx[i])] = nil
		}
	}
	if l.entry != nil {
		l.entry.WithFields(fields).Trace(msg)
	} else {
		l.Logger.WithFields(fields).Trace(msg)
	}
}

func (l *log) Debug(msg string, ctx ...interface{}) {
	fields := make(map[string]interface{})
	for i := 0; i < len(ctx); i = i + 2 {
		if i+2 <= len(ctx) {
			fields[ToString(ctx[i])] = ctx[i+1]
		} else {
			fields[ToString(ctx[i])] = nil
		}
	}
	if l.entry != nil {
		l.entry.WithFields(fields).Debug(msg)
	} else {
		l.Logger.WithFields(fields).Debug(msg)
	}
}

func (l *log) Info(msg string, ctx ...interface{}) {
	fields := make(map[string]interface{})
	for i := 0; i < len(ctx); i = i + 2 {
		if i+2 <= len(ctx) {
			fields[ToString(ctx[i])] = ctx[i+1]
		} else {
			fields[ToString(ctx[i])] = nil
		}
	}
	if l.entry != nil {
		l.entry.WithFields(fields).Info(msg)
	} else {
		l.Logger.WithFields(fields).Info(msg)
	}
}

func (l *log) Warn(msg string, ctx ...interface{}) {
	fields := make(map[string]interface{})
	for i := 0; i < len(ctx); i = i + 2 {
		if i+2 <= len(ctx) {
			fields[ToString(ctx[i])] = ctx[i+1]
		} else {
			fields[ToString(ctx[i])] = nil
		}
	}
	if l.entry != nil {
		l.entry.WithFields(fields).Warn(msg)
	} else {
		l.Logger.WithFields(fields).Warn(msg)
	}
}

func (l *log) Error(msg string, ctx ...interface{}) {
	fields := make(map[string]interface{})
	for i := 0; i < len(ctx); i = i + 2 {
		if i+2 <= len(ctx) {
			fields[ToString(ctx[i])] = ctx[i+1]
		} else {
			fields[ToString(ctx[i])] = nil
		}
	}
	if l.entry != nil {
		l.entry.WithFields(fields).Error(msg)
	} else {
		l.Logger.WithFields(fields).Error(msg)
	}
}

func (l *log) Crit(msg string, ctx ...interface{}) {
	fields := make(map[string]interface{})
	for i := 0; i < len(ctx); i = i + 2 {
		if i+2 <= len(ctx) {
			fields[ToString(ctx[i])] = ctx[i+1]
		} else {
			fields[ToString(ctx[i])] = nil
		}
	}
	if l.entry != nil {
		l.entry.WithFields(fields).Fatal(msg)
	} else {
		l.Logger.WithFields(fields).Fatal(msg)
	}
}
