// Package pool
//
// @author: xwc1125
package pool

import "github.com/panjf2000/ants/v2"

type Pool struct {
	*ants.Pool
}

type PoolWithFunc struct {
	*ants.PoolWithFunc
}

func NewPool(size int, options ...ants.Option) (*Pool, error) {
	pool, err := ants.NewPool(size, options...)
	if err != nil {
		return nil, err
	}
	return &Pool{
		pool,
	}, nil
}

func NewPoolWithFunc(size int, pf func(interface{}), options ...ants.Option) (*PoolWithFunc, error) {
	pool, err := ants.NewPoolWithFunc(size, pf, options...)
	if err != nil {
		return nil, err
	}
	return &PoolWithFunc{
		pool,
	}, nil
}
