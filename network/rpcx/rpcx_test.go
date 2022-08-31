// Package rpcx
//
// @author: xwc1125
package rpcx

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"golang.org/x/net/context"
)

type AA struct {
	AA string
}

func (a *AA) Version() (string, error) {
	return "1.0.0", nil
}
func (a *AA) Add(aa, bb int) (int, error) {
	return aa + bb, nil
}

type BB struct {
}

func (a *BB) Version(ctx context.Context) (string, error) {
	return "1.0.0", nil
}
func (a *BB) Sub(aa, bb int64) (int64, error) {
	return aa - bb, nil
}
func (a *BB) InfoInt(a1 int, a2 int8, a3 int16, a4 int32, a5 int64) {
	fmt.Println(a1, a2, a3, a4, a5)
}
func (a *BB) InfoUint(a1 uint, a2 uint8, a3 uint16, a4 uint32, a5 uint64) {
	fmt.Println(a1, a2, a3, a4, a5)
}
func (a *BB) InfoFloat(a1 float32, a2 float64) {
	fmt.Println(a1, a2)
}
func (a *BB) InfoOther(a1 bool, a2 byte, a3 []string, a4 [1]int) {
	fmt.Println(a1, a2, a3, a4)
}
func (a *BB) InfoM(a1 map[string]int, a2 interface{}, a3 []interface{}, a4 AA) {
	fmt.Println(a1, a2, a3, a4)
}
func (a *BB) InfoM1(a1 *big.Int, a2 *hexutil.Big) {
	fmt.Println(a1, a2)
}

func TestRPC2(t *testing.T) {
	rpc := New()
	rpc.Register("aa", new(AA))
	rpc.Register("bb", new(BB))

	{
		fmt.Println("========aa.version()========")
		request := rpc.Exec(context.Background(), "aa", "version", "123", nil)
		fmt.Println(request)
	}
	{
		fmt.Println("========aa.add()========")
		request := rpc.Exec(context.Background(), "aa", "add", "123", 1, 2)
		fmt.Println(request)
		// string= int ,name= int ,Kind= int
		// string= int ,name= int ,Kind= int
	}
	{
		fmt.Println("========bb.Version()========")
		request := rpc.Exec(context.Background(), "bb", "Version", "123", context.Background())
		fmt.Println(request)
	}
	{
		fmt.Println("========bb.Sub()========")
		request := rpc.Exec(context.Background(), "bb", "Sub", "123", "1", "2")
		fmt.Println(request)
		// string= int64 ,name= int64 ,Kind= int64
		// string= int64 ,name= int64 ,Kind= int64
	}
	{
		fmt.Println("========bb.InfoInt()========")
		// InfoInt(a1 int, a2 int8, a3 int16, a4 int32, a5 int64)
		request := rpc.Exec(context.Background(), "bb", "InfoInt", "123", 1, int8(2), int16(3), int32(4), int64(5))
		fmt.Println(request)
		// string= int ,name= int ,Kind= int
		// string= int8 ,name= int8 ,Kind= int8
		// string= int16 ,name= int16 ,Kind= int16
		// string= int32 ,name= int32 ,Kind= int32
		// string= int64 ,name= int64 ,Kind= int64
	}
	{
		fmt.Println("========bb.InfoUint()========")
		// InfoUint(a1 uint, a2 uint8, a3 uint16, a4 uint32, a5 uint64)
		request := rpc.Exec(context.Background(), "bb", "InfoUint", "123", uint(1), uint8(2), uint16(3), uint32(4), uint64(5))
		fmt.Println(request)
		// string= uint ,name= uint ,Kind= uint
		// string= uint8 ,name= uint8 ,Kind= uint8
		// string= uint16 ,name= uint16 ,Kind= uint16
		// string= uint32 ,name= uint32 ,Kind= uint32
		// string= uint64 ,name= uint64 ,Kind= uint64
	}
	{
		fmt.Println("========bb.InfoFloat()========")
		// InfoFloat(a1 float32, a2 float64)
		request := rpc.Exec(context.Background(), "bb", "InfoFloat", "123", float32(1.0), 3.14)
		fmt.Println(request)
		// string= float32 ,name= float32 ,Kind= float32
		// string= float64 ,name= float64 ,Kind= float64
	}
	{
		fmt.Println("========bb.InfoOther()========")
		// InfoOther(a1 bool, a2 byte, a3 []string, a4 [1]int)
		request := rpc.Exec(context.Background(), "bb", "InfoOther", "123", true, byte(1), []string{"aa"}, [1]int{1})
		fmt.Println(request)
		// string= bool ,name= bool ,Kind= bool
		// string= uint8 ,name= uint8 ,Kind= uint8
		// string= []string ,name=  ,Kind= slice
		// string= [1]int ,name=  ,Kind= array
	}
	{
		fmt.Println("========bb.InfoM()========")
		// InfoM(a1 map[string]int, a2 interface{}, a3 []interface{}, a4 AA)
		request := rpc.Exec(context.Background(), "bb", "InfoM", "123", map[string]int{
			"11": 11,
			"22": 22,
		}, 1, []interface{}{"aa"}, "{\n    \"AA\": \"aa\"\n}")
		fmt.Println(request)
		// string= map[string]int ,name=  ,Kind= map
		// string= interface {} ,name=  ,Kind= interface
		// string= []interface {} ,name=  ,Kind= slice
		// string= rpcx.AA ,name= AA ,Kind= struct
	}
	{
		fmt.Println("========bb.InfoM1()========")
		// InfoM(a1 map[string]int, a2 interface{}, a3 []interface{}, a4 AA)
		request := rpc.Exec(context.Background(), "bb", "InfoM1", "123", "10", "\"0xa\"")
		fmt.Println(request)
		// string= map[string]int ,name=  ,Kind= map
		// string= interface {} ,name=  ,Kind= interface
		// string= []interface {} ,name=  ,Kind= slice
		// string= rpcx.AA ,name= AA ,Kind= struct
	}
	fmt.Println("==========================")
	namespaces := rpc.Namespaces(context.Background())
	bytes, _ := json.Marshal(namespaces)
	fmt.Println(string(bytes))
}

func TestBigInt(t *testing.T) {
	newInt := big.NewInt(10)
	bytes, _ := json.Marshal(newInt)
	fmt.Println(string(bytes))
	var newInt1 big.Int
	json.Unmarshal(bytes, &newInt1)
	fmt.Println(newInt1)

	var hexBig hexutil.Big = hexutil.Big(*big.NewInt(10))
	bytes1, _ := json.Marshal(hexBig)
	fmt.Println(string(bytes1))
	var hexBig1 hexutil.Big
	err := json.Unmarshal(bytes1, &hexBig1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hexBig1)
}

type CC struct {
}

func (c *CC) Addr(addr Addr) Addr {
	return addr
}

type Addr string

func TestCC(t *testing.T) {
	rpc := New()
	rpc.Register("cc", new(CC))
	{
		fmt.Println("========cc.addr()========")
		result := rpc.Exec(context.Background(), "cc", "addr", "12345", "11")
		fmt.Println(result)
	}
}
