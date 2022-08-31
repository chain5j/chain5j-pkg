// Package reflectutil
//
// @author: xwc1125
package reflectutil

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

const name = "常量" // 常量doc

var age = 10 // 参数doc

// 学生表
type Student struct {
	Name  string  `chain5j:"stu_name,nil"` // 姓名
	Age   int     // 年龄
	Score float32 // 分数
	sex   int     // 性别
}

// 打印内容
func (s Student) Print() {
	fmt.Println(s)
}

// 内部打印内容
func (s Student) print() {
	fmt.Println(s)
}

func TestTag(t *testing.T) {
	var a = Student{
		Name:  "stu01",
		Age:   18,
		Score: 92.8,
	}
	TagStruct(&a)

	valueOf := reflect.ValueOf(&a)
	s := valueOf.Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		// 获取字段内容
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, // 字段名称
			f.Type(),              // 字段类型
			f.Interface())         // 字段值
	}
	// 设置
	s.Field(0).SetString("stu")
	s.Field(1).SetInt(77)
	s.Field(2).SetFloat(77.0)
	fmt.Println("a is now", a)

	v := valueOf.MethodByName("Print")
	v.Call([]reflect.Value{})
}

func TagStruct(a interface{}) {
	typ := reflect.TypeOf(a)

	tag := typ.Elem().Field(0).Tag.Get("chain5j")
	fmt.Printf("Tag:%s\n", tag)
}

func TestDoc(t *testing.T) {
	fset := token.NewFileSet() // positions are relative to fset
	d, err := parser.ParseDir(fset, "./", nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, f := range d {
		fmt.Println("package", k)
		for n, f := range f.Files {
			// 文件列表
			fmt.Println(fmt.Sprintf("文件名称: %q", n))
			// f.Doc 文件注释
			// f.Scope 当前文件下的直接调用的方法/struct
			// f.Imports 引入的文件
			for i, c := range f.Comments {
				// 所有的注释
				fmt.Println(fmt.Sprintf("Comment Group %d", i))
				for _, c1 := range c.List {
					// fmt.Println(fmt.Sprintf("Comment %d: Position: %d, Text: %q", i2, c1.Slash, c1.Text))
					fmt.Println(fmt.Sprintf("Text: %q", c1.Text))
				}
			}
		}

		p := doc.New(f, "./", doc.AllDecls) // doc.AllDecls显示私有和公有，AllMethods：显示公有
		// p.Name 包名
		// p.ImportPath 引入地址
		// p.Imports 引入
		// p.Filenames 包下的文件
		// p.Filenames 包下的文件
		// p.Consts 常量的声明
		// p.Types type的声明
		// p.Vars var的声明
		// p.Funcs func的声明

		for _, t := range p.Types {
			fmt.Println("type=", t.Name, "docs=", t.Doc)
			for _, m := range t.Methods {
				fmt.Println("type=", m.Name, "docs=", m.Doc)
			}
		}

		for _, v := range p.Vars {
			fmt.Println("type", v.Names)
			fmt.Println("docs:", v.Doc)
		}

		for _, f := range p.Funcs {
			fmt.Println("type", f.Name)
			fmt.Println("docs:", f.Doc)
		}

		for _, n := range p.Notes {
			fmt.Println("body", n[0].Body)
		}
	}
}
