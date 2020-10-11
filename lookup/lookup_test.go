// description: lookup
// 
// @author: xwc1125
// @date: 2020/3/20
package lookup

import (
	"fmt"
	"testing"
)

var (
	allKeyList map[string]struct{}
)

func Test1(t *testing.T) {
	allKeyList = make(map[string]struct{})
	hash := "0x57f09E3a3E66F4C70B19Be63dCa6e1bF72eA4Eb1"
	allKeyList[hash] = struct{}{}
	fmt.Println("pre", allKeyList)
	delete(allKeyList, hash)
	fmt.Println("post", allKeyList)
	if k, ok := allKeyList[hash]; ok {
		fmt.Println("Post", k, "ok=", ok)
	}
}
