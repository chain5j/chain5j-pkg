package dateutil

import (
	"fmt"
	log "github.com/chain5j/log15"
	"strings"
	"time"
)

/*
前面是含义，后面是 go 的表示值,多种表示,逗号","分割
年　 06,2006
月份 1,01,Jan,January
日　 2,02,_2
时　 3,03,15,PM,pm,AM,am
分　 4,04
秒　 5,05
周几 Mon,Monday
时区时差表示 -07,-0700,Z0700,Z07:00,-07:00,MST
时区字母缩写 MST
您看出规律了么！哦是的，你发现了，这里面没有一个是重复的，所有的值表示都唯一对应一个时间部分。
并且涵盖了很多格式组合。
*/
type DateFormat string

const (
	// 时间格式化字符串
	Default               DateFormat = "2006-01-02 15:04:05"
	YYYY                  DateFormat = "2006"
	YYYY_MM               DateFormat = "2006-01"
	YYYY_7MM              DateFormat = "2006/01"
	YYYY_0MM              DateFormat = "2006.01"
	YYYYMM                DateFormat = "200601"
	YYYY_MM_dd            DateFormat = "2006-01-02"
	YYYY_7MM_7dd          DateFormat = "2006/01/02"
	YYYY_0MM_0dd          DateFormat = "2006.01.02"
	YYYYMMdd              DateFormat = "20060102"
	YYYY_MM_dd8HH0mm      DateFormat = "2006-01-02 15:04"
	YYYY_7MM_7dd8HH0mm    DateFormat = "2006/01/02 15:04"
	YYYY_0MM_0dd8HH0mm    DateFormat = "2006.01.02 15:04"
	YYYYMMddHHmm          DateFormat = "200601021504"
	YYYY_MM_dd8HH0mm0ss   DateFormat = "2006-01-02 15:04:05"
	YYYY_7MM_7dd8HH0mm0ss DateFormat = "2006/01/02 15:04:05"
	YYYY_0MM_0dd8HH0mm0ss DateFormat = "2006.01.02 15:04:05"
	YYYYMMddHHmmss        DateFormat = "20060102150405"
	HH0mm0ss              DateFormat = "15:04:05"
)

func (d DateFormat) String() string {
	return string(d)
}

// 中国时区
var SysTimeLocation, _ = time.LoadLocation("Asia/Chongqing")

// 返回毫秒
func CurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}

// 纳秒转毫秒
func NanoToMillisecond(t int64) int64 {
	return t / 1e6
}

// 纳秒转秒
func NanoToSecond(t int64) int64 {
	return t / 1e9
}

// 秒转time
func SecondToTime(t int64) time.Time {
	return time.Unix(t, 0)
}

func MillisecondToTime(t int64) time.Time {
	return time.Unix(0, t*1e6)
}
func NanoToTime(t int64) time.Time {
	return time.Unix(0, t)
}

// 格式化输出
func Format(t time.Time, format DateFormat) string {
	return t.Format(format.String())
}

// 秒转成format
func SecondFormat(t int64, format DateFormat) string {
	return Format(SecondToTime(t), format)
}

// 时间转本地化
// s时间格式：如"2017-05-11 14:06:06"
// format：格式
// location：时区(Location)
func ParseInLocation(s string, format DateFormat, location *time.Location) {
	// func ParseInLocation(layout, value string, loc *Location) (Time, error) (layout已带时区时可直接用Parse)
	time.ParseInLocation(format.String(), s, location)
}

func LoadLocation(loc string) *time.Location {
	// 默认UTC
	//loc, err := time.LoadLocation("")
	//// 服务器设定的时区，一般为CST
	//loc, err := time.LoadLocation("Local")
	//// 美国洛杉矶PDT
	//loc, err := time.LoadLocation("America/Los_Angeles")
	//
	//// 获取指定时区的时间点
	//local, _ := time.LoadLocation("America/Los_Angeles")
	//fmt.Println(time.Date(2018,1,1,12,0,0,0, local))
	location, e := time.LoadLocation(loc)
	if e != nil {
		log.Error("loadLocation err", "err", e)
		location, _ = time.LoadLocation("")
		return location
	}
	return location
}

// GetDistanceTime 获取间隔时间错，传入的是毫秒
func GetDistanceTime(diffMS int64) string {
	isNegative := false
	if diffMS < 0 {
		isNegative = true
		diffMS = -diffMS
	}
	s := diffMS / 1000 // 秒
	m := s / 60        //分钟
	h := m / 60        //小时
	day := h / 24      //天
	hour := h - 24*day
	min := m - h*60
	sec := s - m*60
	ms := diffMS - s*1000
	var buff strings.Builder
	if isNegative {
		buff.WriteString("-")
	}
	if day > 0 {
		buff.WriteString(fmt.Sprintf("%dd", day))
	}
	if hour > 0 {
		buff.WriteString(fmt.Sprintf("%dh", hour))
	}
	if min > 0 {
		buff.WriteString(fmt.Sprintf("%dm", min))
	}
	if sec > 0 {
		buff.WriteString(fmt.Sprintf("%ds", sec))
	}
	if ms > 0 {
		buff.WriteString(fmt.Sprintf("%dms", ms))
	}
	return buff.String()
}

// GetDistanceTimeToCurrent 传入的是毫秒值
func GetDistanceTimeToCurrent(startTime int64) string {
	diff := CurrentTime() - startTime
	return GetDistanceTime(diff)
}

// GetAddDayTime 添加天数
// 如果day是负数,那么是往前算
func GetAddDayTime(day int) (days int64) {
	//获取两天前的时间
	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, day)
	startTime := time.Date(oldTime.Year(), oldTime.Month(), oldTime.Day(), 0, 0, 0, 0, oldTime.Location())
	return startTime.UnixNano()
}