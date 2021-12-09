// Package linkedHashMap
//
// @author: xwc1125
package linkedHashMap

import (
	"crypto/rand"
	"fmt"
	"github.com/chain5j/chain5j-pkg/util/dateutil"
	"log"
	"math/big"
	"strconv"
	"sync"
	"testing"
)

func TestLinkedHashMap(t *testing.T) {
	var wg sync.WaitGroup
	linkedHashMap := NewLinkedHashMap()
	startTime := dateutil.CurrentTime()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go add(linkedHashMap, i, 1000, &wg)
	}
	wg.Wait()
	endTime := dateutil.CurrentTime()
	log.Println(fmt.Sprintf("【添加】end-start=%d,len=%d", endTime-startTime, linkedHashMap.Len()))

	ids := make(map[string]string, 0)
	for i := 0; i < 100; i++ {
		ids2 := gene(i, 100)
		for _, id := range ids2 {
			ids[id] = id
		}
	}
	log.Println(fmt.Sprintf("【ids】len=%d", len(ids)))
}

func TestLinkedHashMapBatch(t *testing.T) {
	var wg sync.WaitGroup
	linkedHashMap := NewLinkedHashMap()
	startTime := dateutil.CurrentTime()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go addBatch(linkedHashMap, i, 1000, &wg)
	}
	wg.Wait()
	endTime := dateutil.CurrentTime()
	log.Println(fmt.Sprintf("【添加】end-start=%d,len=%d", endTime-startTime, linkedHashMap.Len()))

	var count = 0
	list := linkedHashMap.GetLinkList()
	if list != nil && list.Len() > 0 {
		head := list.Front()
		for head != nil && count < 10 {
			fmt.Println(head.Value)
			count++
			head = head.Next()
		}
	}

	ids := make(map[string]string, 0)
	for i := 0; i < 1000; i++ {
		ids2 := gene(i, 1000)
		for _, id := range ids2 {
			ids[id] = id
		}
	}
	log.Println(fmt.Sprintf("【ids】len=%d", len(ids)))
}

func add(linkedHashMap *LinkedHashMap, index int, count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		linkedHashMap.Add(strconv.Itoa(index)+"_"+strconv.Itoa(i), strconv.Itoa(index)+"_"+strconv.Itoa(i))
	}
}

func addBatch(linkedHashMap *LinkedHashMap, index int, count int, wg *sync.WaitGroup) {
	defer wg.Done()
	kvs := make([]KV, count)
	for i := 0; i < count; i++ {
		kvs[i] = KV{
			Key: strconv.Itoa(index) + "_" + strconv.Itoa(i),
			Val: strconv.Itoa(index) + "_" + strconv.Itoa(i),
		}
	}
	linkedHashMap.BatchAdd(kvs...)
}

func gene(index int, count int) []string {
	ids := make([]string, count)
	for j := 0; j < count; j++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(count)))
		ids[j] = strconv.Itoa(index) + "_" + strconv.Itoa(int(n.Int64()))
	}
	return ids
}
