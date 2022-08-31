// Package groupwork
//
// @author: xwc1125
package groupwork

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

// forever 持续打印数字,直到ctx结束
func forever(ctx context.Context, i int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("forever stop now")
			return
		default:
		}
		time.Sleep(time.Second)
		fmt.Printf("goroutine %d is runnint\n", i)
		runtime.Gosched()
	}
}

// delayError 5秒后报错
func delayError() error {
	time.Sleep(time.Second * 5)
	fmt.Println("dealyError return")
	return errors.New("should stop now")
}

func TestGroup_Go(t *testing.T) {
	g, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < 2; i++ {
		i := i
		if i == 0 {
			g.Go(delayError)
			continue
		}
		g.Go(func() error {
			forever(ctx, i)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
