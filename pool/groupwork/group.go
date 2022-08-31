// Package groupwork
//
// @author: xwc1125
package groupwork

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/chain5j/chain5j-pkg/pool/pool"
	"github.com/chain5j/logger"
)

type Group struct {
	pool *pool.Pool

	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error

	okOnce sync.Once
	ok     interface{}
}

func WithContext(ctx context.Context, poolSize int) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	newPool, _ := pool.NewPool(poolSize)
	return &Group{
		cancel: cancel,
		pool:   newPool,
	}, ctx
}

func (g *Group) WaitErr() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
func (g *Group) WaitOk() interface{} {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.ok
}

func (g *Group) OnceErr(f func() error) {
	g.wg.Add(1)
	g.pool.Submit(func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	})
}
func (g *Group) OnceOk(f func() (interface{}, error)) {
	g.wg.Add(1)
	g.pool.Submit(func() {
		defer g.wg.Done()

		if data, err := f(); err == nil {
			g.okOnce.Do(func() {
				g.ok = data
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	})
}

func GoAndWaitErr(handlers ...func() error) (err error) {
	var wg sync.WaitGroup
	var once sync.Once
	for _, f := range handlers {
		wg.Add(1)
		go func(handler func() error) {
			defer func() {
				if e := recover(); e != nil {
					buf := make([]byte, 1024)
					buf = buf[:runtime.Stack(buf, false)]
					logger.Error("[PANIC]%v\n%s\n", e, buf)
					once.Do(func() {
						err = errors.New("panic found in call handlers")
					})
				}
				wg.Done()
			}()
			if e := handler(); e != nil {
				once.Do(func() {
					err = e
				})
			}
		}(f)
	}
	wg.Wait()
	return err
}
func GoAndWaitOk(handlers ...func() (interface{}, error)) (result interface{}, err error) {
	var wg sync.WaitGroup
	var once sync.Once
	for _, f := range handlers {
		wg.Add(1)
		go func(handler func() (interface{}, error)) {
			defer func() {
				if e := recover(); e != nil {
					buf := make([]byte, 1024)
					buf = buf[:runtime.Stack(buf, false)]
					logger.Error("[PANIC]%v\n%s\n", e, buf)
					err = errors.New("panic found in call handlers")
				}
				wg.Done()
			}()
			if f, e := handler(); e == nil {
				once.Do(func() {
					result = f
				})
			}
		}(f)
	}
	wg.Wait()
	return result, err
}
