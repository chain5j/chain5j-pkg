// Package log15
//
// @author: xwc1125
// @date: 2021/8/30
package log15

import (
	"github.com/chain5j/chain5j-pkg/logger"
	log "github.com/chain5j/log15"
)

var (
	_ logger.Logger = new(log15)
)

type log15 struct {
	log.Logger
}

func (l *log15) New(module string, ctx ...interface{}) logger.Logger {
	return &log15{
		log.New(module, ctx),
	}
}

func (l *log15) Trace(msg string, ctx ...interface{}) {
	l.Logger.Trace(msg, ctx)
}

func (l *log15) Debug(msg string, ctx ...interface{}) {
	l.Logger.Debug(msg, ctx)
}

func (l *log15) Info(msg string, ctx ...interface{}) {
	l.Logger.Info(msg, ctx)
}

func (l *log15) Warn(msg string, ctx ...interface{}) {
	l.Logger.Warn(msg, ctx)
}

func (l *log15) Error(msg string, ctx ...interface{}) {
	l.Logger.Error(msg, ctx)
}

func (l *log15) Crit(msg string, ctx ...interface{}) {
	l.Logger.Crit(msg, ctx)
}

func (l *log15) Printf(format string, v ...interface{}) {
	l.Logger.Printf(format, v)
}

func (l *log15) Print(v ...interface{}) {
	l.Logger.Print(v)
}

func (l *log15) Println(v ...interface{}) {
	l.Logger.Println(v)
}

func (l *log15) Fatal(v ...interface{}) {
	l.Logger.Fatal(v)
}

func (l *log15) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(format, v)
}

func (l *log15) Fatalln(v ...interface{}) {
	l.Logger.Fatalln(v)
}

func (l *log15) Panic(v ...interface{}) {
	l.Logger.Panic(v)
}

func (l *log15) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v)
}

func (l *log15) Panicln(v ...interface{}) {
	l.Logger.Panicln(v)
}
