// Package color
//
// @author: xwc1125
package color

import (
	"fmt"
	"reflect"
	"strings"
)

// ShowType 显示方式
type ShowType int

const (
	ShowType_Terminal  ShowType = 0 // 终端默认设置
	ShowType_Highlight ShowType = 1 // 高亮显示
	ShowType_Underline ShowType = 4 // 使用下划线
	ShowType_Blink     ShowType = 5 // 闪烁
	ShowType_Reverse   ShowType = 7 // 反白显示
	ShowType_Invisible ShowType = 8 // 不可见
)

type ForegroundColor int

const (
	ForegroundColor_Black  ForegroundColor = iota + 30 // 30
	ForegroundColor_Red                                // 31
	ForegroundColor_Green                              // 32
	ForegroundColor_Yellow                             // 33
	ForegroundColor_Blue                               // 34
	ForegroundColor_Purple                             // 35
	ForegroundColor_Cyan                               // 36
	ForegroundColor_White                              // 37
)

type BackgroundColor int

const (
	BackgroundColor_Transparent BackgroundColor = iota + 39 // 39
	BackgroundColor_Black                                   // 40
	BackgroundColor_Red                                     // 41
	BackgroundColor_Green                                   // 42
	BackgroundColor_Yellow                                  // 43
	BackgroundColor_Blue                                    // 44
	BackgroundColor_Fuchsia                                 // 45
	BackgroundColor_Cyan                                    // 46
	BackgroundColor_White                                   // 47
)

// Black 黑色
func Black(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Black), 0, modifier...)
}

// DarkGray 深灰色
func DarkGray(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Black), 1, modifier...)
}

// Red 红字体
func Red(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Red), 0, modifier...)
}

// LightRed 淡红色
func LightRed(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Red), 1, modifier...)
}

// Green 绿色字体，modifier里，第一个控制闪烁，第二个控制下划线
func Green(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Green), 0, modifier...)
}

// LightGreen 淡绿
func LightGreen(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Green), 1, modifier...)
}

// Yellow 黄色字体
func Yellow(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Yellow), 0, modifier...)
}

// Brown 棕色
func Brown(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Yellow), 0, modifier...)
}

// Blue 蓝色
func Blue(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Blue), 0, modifier...)
}

// LightBlue 淡蓝
func LightBlue(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Blue), 1, modifier...)
}

// Cyan 青色/蓝绿色
func Cyan(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Cyan), 0, modifier...)
}

// LightCyan 淡青色
func LightCyan(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Cyan), 1, modifier...)
}

// White 白色
func White(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_White), 1, modifier...)
}

// LightGray 浅灰色
func LightGray(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_White), 0, modifier...)
}

// Purple 紫色
func Purple(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Purple), 0, modifier...)
}

// LightPurple 淡紫色
func LightPurple(str string, modifier ...interface{}) string {
	return cliColorRender(str, int(ForegroundColor_Purple), 1, modifier...)
}

// cliColorRender color颜色处理
// str 文本
// color 颜色值
// weight 权重
// extraArgs [0]闪烁 [1]下划线
func cliColorRender(str string, color int, weight int, extraArgs ...interface{}) string {
	// 闪烁效果
	isBlink := int64(0)
	if len(extraArgs) > 0 {
		isBlink = reflect.ValueOf(extraArgs[0]).Int()
	}
	// 下划线效果
	isUnderLine := int64(0)
	if len(extraArgs) > 1 {
		isUnderLine = reflect.ValueOf(extraArgs[1]).Int()
	}
	var mo []string
	if isBlink > 0 {
		mo = append(mo, "05")
	}
	if isUnderLine > 0 {
		mo = append(mo, "04")
	}
	if weight > 0 {
		mo = append(mo, fmt.Sprintf("%d", weight))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	return fmt.Sprintf("\033[%s;%dm"+str+"\033[0m", strings.Join(mo, ";"), color)
}

// ColorRender color颜色处理
// str 文本
// foregroundColor 前景色
// backgroundColor 背景色
// extraArgs 扩展的显示方式
func ColorRender(str string, foregroundColor ForegroundColor, backgroundColor BackgroundColor, extraArgs ...ShowType) string {
	var mo []string
	if len(extraArgs) > 0 {
		for _, arg := range extraArgs {
			mo = append(mo, fmt.Sprintf("%d", arg))
		}
	}
	if backgroundColor >= BackgroundColor_Transparent && backgroundColor <= BackgroundColor_White {
		mo = append(mo, fmt.Sprintf("%d", backgroundColor))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	return fmt.Sprintf("\033[%s;%dm"+str+"\033[0m", strings.Join(mo, ";"), foregroundColor)
}
