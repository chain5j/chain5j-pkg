// Package color
//
// @author: xwc1125
package color

import (
	"fmt"
	"strings"
	"testing"
)

func TestColor(t *testing.T) {
	// 默认的不带任何效果的字体显示
	fmt.Println(Green("字体：Green"))
	fmt.Println(LightGreen("字体：LightGreen"))
	fmt.Println(Cyan("字体：Cyan"))
	fmt.Println(LightCyan("字体：LightCyan"))
	fmt.Println(Red("字体：Red"))
	fmt.Println(LightRed("字体：LightRed"))
	fmt.Println(Yellow("字体：Yellow"))
	fmt.Println(Black("字体：Black"))
	fmt.Println(DarkGray("字体：DarkGray"))
	fmt.Println(LightGray("字体：LightGray"))
	fmt.Println(White("字体：White"))
	fmt.Println(Blue("字体：Blue"))
	fmt.Println(LightBlue("字体：LightBlue"))
	fmt.Println(Purple("字体：Purple"))
	fmt.Println(LightPurple("字体：LightPurple"))
	fmt.Println(Brown("字体：Brown"))
	fmt.Println(Blue("字体：Blue", 1, 1))

	// 带闪烁效果的彩色字体显示
	fmt.Println(Green("闪烁字体：Green", 1, 1))
	fmt.Println(LightGreen("闪烁字体：LightGreen", 1))
	fmt.Println(Cyan("闪烁字体：Cyan", 1))
	fmt.Println(LightCyan("闪烁字体：LightCyan", 1))
	fmt.Println(Red("闪烁字体：Red", 1))
	fmt.Println(LightRed("闪烁字体：LightRed", 1))
	fmt.Println(Yellow("闪烁字体：Yellow", 1))
	fmt.Println(Black("闪烁字体：Black", 1))
	fmt.Println(DarkGray("闪烁字体：DarkGray", 1))
	fmt.Println(LightGray("闪烁字体：LightGray", 1))
	fmt.Println(White("闪烁字体：White", 1))
	fmt.Println(Blue("闪烁字体：Blue", 1))
	fmt.Println(LightBlue("闪烁字体：LightBlue", 1))
	fmt.Println(Purple("闪烁字体：Purple", 1))
	fmt.Println(LightPurple("闪烁字体：LightPurple", 1))
	fmt.Println(Brown("闪烁字体：Brown", 1))

	// 带下划线效果的字体显示
	fmt.Println(Green("闪烁且带下划线字体：Green", 1, 1, 1))
	fmt.Println(LightGreen("闪烁且带下划线字体：LightGreen", 1, 1))
	fmt.Println(Cyan("闪烁且带下划线字体：Cyan", 1, 1))
	fmt.Println(LightCyan("闪烁且带下划线字体：LightCyan", 1, 1))
	fmt.Println(Red("闪烁且带下划线字体：Red", 1, 1))
	fmt.Println(LightRed("闪烁且带下划线字体：LightRed", 1, 1))
	fmt.Println(Yellow("闪烁且带下划线字体：Yellow", 1, 1))
	fmt.Println(Black("闪烁且带下划线字体：Black", 1, 1))
	fmt.Println(DarkGray("闪烁且带下划线字体：DarkGray", 1, 1))
	fmt.Println(LightGray("闪烁且带下划线字体：LightGray", 1, 1))
	fmt.Println(White("闪烁且带下划线字体：White", 1, 1))
	fmt.Println(Blue("闪烁且带下划线字体：Blue", 1, 1))
	fmt.Println(LightBlue("闪烁且带下划线字体：LightBlue", 1, 1))
	fmt.Println(Purple("闪烁且带下划线字体：Purple", 1, 1))
	fmt.Println(LightPurple("闪烁且带下划线字体：LightPurple", 1, 1))
	fmt.Println(Brown("闪烁且带下划线字体：Brown", 1, 1))
}

// Go语言要打印彩色字符与Linux终端输出彩色字符类似，以黑色背景高亮绿色字体为例：
// fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
// 其中0x1B是标记，[开始定义颜色，1代表高亮，40代表黑色背景，32代表绿色前景，0代表恢复默认颜色。显示效果为：
func TestColor2(t *testing.T) {

	fmt.Println(Green("字体：Green"))
	fmt.Println(ColorRender("字体：Green", ForegroundColor_Green, BackgroundColor_Transparent))
	fmt.Println(Green("闪烁字体：Green", 1))
	fmt.Println(ColorRender("闪烁字体：Green", ForegroundColor_Green, BackgroundColor_Transparent, []ShowType{
		ShowType_Blink,
	}...))
	fmt.Println(Green("闪烁且带下划线字体：Green", 1, 1))
	fmt.Println(ColorRender("闪烁且带下划线字体：Green", ForegroundColor_Green, BackgroundColor_Transparent, []ShowType{
		ShowType_Blink,
		ShowType_Underline,
	}...))

	var foreground = 30
	var background = 41
	var typ1 = 1
	msg := "字体：Green"
	sprintf := fmt.Sprintf("%c[%d;%d;%dm%s%s%c[0m", 0x1B, typ1, background, foreground, "", msg, 0x1B)
	fmt.Println(sprintf)
	//  0  终端默认设置
	//  1  高亮显示
	//  4  使用下划线
	//  5  闪烁
	//  7  反白显示
	//  8  不可见
	var mo = []string{
		"05", // 闪烁
		// "04", // 使用下划线
		// "07", // 反白显示
		fmt.Sprintf("%d", 1), // 高亮显示
		// fmt.Sprintf("%d", 0),// 终端默认设置
		// fmt.Sprintf("%d", 8), // 不可见
		fmt.Sprintf("%d", background), // 不可见
	}
	sprintf2 := fmt.Sprintf("\033[%s;%dm"+msg+"\033[0m", strings.Join(mo, ";"), foreground)
	fmt.Println(sprintf2)
	test1()
	// test2()
}

func test1() {
	fmt.Println("")
	// 前景 背景 颜色
	// ---------------------------------------
	//  -  39  透明
	// 30  40  黑色
	// 31  41  红色
	// 32  42  绿色
	// 33  43  黄色
	// 34  44  蓝色
	// 35  45  紫红色
	// 36  46  青蓝色
	// 37  47  白色
	//
	// 代码 意义
	// -------------------------
	//  0  终端默认设置
	//  1  高亮显示
	//  4  使用下划线
	//  5  闪烁
	//  7  反白显示
	//  8  不可见

	for background := 39; background <= 47; background++ { // 背景色彩 = 39-47
		for foreground := 30; foreground <= 37; foreground++ { // 前景色彩 = 30-37
			for _, t := range []int{0, 1, 4, 5, 7, 8} { // 显示方式 = 0,1,4,5,7,8
				msg := fmt.Sprintf("(f=%d,b=%d,t=%d)", foreground, background, t)
				sprintf := fmt.Sprintf(" %c[%d;%d;%dm%s"+msg+"%c[0m ", 0x1B, t, background, foreground, "", 0x1B)
				fmt.Print(sprintf)
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func test2() {
	// 0x1B是标记
	// 1代表高亮，40代表黑色背景，32代表绿色前景，0代表恢复默认颜色
	fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
	fmt.Printf("\n %c[1;34;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
}
