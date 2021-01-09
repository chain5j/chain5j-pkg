// description: chain5j-pkg
// 
// @author: xwc1125
// @date: 2020/12/22
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
	// 查询
	startTime = dateutil.CurrentTime()
	var qCount = 0
	for i, id := range ids {
		v := linkedHashMap.Get(id)
		if v == nil {
			log.Println(fmt.Sprintf("【Get】no key,id=%s,i=%s", id, i))
			qCount++
		} else {
			//log.Println(fmt.Sprintf("key=%s,v=%s", id, v))
		}
	}
	endTime = dateutil.CurrentTime()
	log.Println(fmt.Sprintf("【查询】end-start=%d,qCount=%d", endTime-startTime, qCount))
	// 删除
	startTime = dateutil.CurrentTime()
	var dCount = 0
	for i, id := range ids {
		b, v := linkedHashMap.Remove(id)
		if !b {
			log.Println(fmt.Sprintf("【Remove】no key2,id=%s,i=%s", id, i))
			dCount++
			_ = v
		} else {
			//log.Println(fmt.Sprintf("[rm]id=%s,i=%s", id, v))
		}
	}
	endTime = dateutil.CurrentTime()
	log.Println(fmt.Sprintf("【删除】end-start=%d,dCount=%d", endTime-startTime, dCount))
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
	if list != nil && list.GetLength() > 0 {
		head := list.GetHead()
		for head != nil && count < 10 {
			fmt.Println(head.GetVal())
			count++
			head = head.GetNext()
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
	// 查询
	startTime = dateutil.CurrentTime()
	var qCount = 0
	keys := make([]string, 0)
	for _, id := range ids {
		keys = append(keys, id)
	}

	gets := linkedHashMap.BatchGet(keys...)
	endTime = dateutil.CurrentTime()
	log.Println(fmt.Sprintf("【查询】end-start=%d,qCount=%d,gets=%d", endTime-startTime, qCount, len(gets)))
	// 删除
	startTime = dateutil.CurrentTime()
	var dCount = 0

	removes := linkedHashMap.BatchRemove(keys...)
	dCount = len(removes)
	endTime = dateutil.CurrentTime()
	log.Println(fmt.Sprintf("【删除】end-start=%d,dCount=%d", endTime-startTime, dCount))
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
