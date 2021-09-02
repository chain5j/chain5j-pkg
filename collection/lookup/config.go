// Package lookup
//
// @author: xwc1125
// @date: 2019/11/20
package lookup

import "time"

type PoolConfig struct {
	MaxTxSize     uint64 // Max size of tx pool
	PriceBump     int    // Price bump to decide whether to replace tx or not
	BatchTimeout  time.Duration
	BatchCapacity int

	TxLifeTime time.Duration // 分钟
	TxTaskTime time.Duration // 刷新时间（秒）
}

func DefaultConfig() *PoolConfig {
	return &PoolConfig{
		MaxTxSize:     4096,
		PriceBump:     10,
		BatchTimeout:  10000 * time.Millisecond,
		BatchCapacity: 1000,
		TxLifeTime:    1,
		TxTaskTime:    10,
	}
}
