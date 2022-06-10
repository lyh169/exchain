package types

import (
	"sync"
)

var prefetchAddressPool = sync.Pool{
	New: func() interface{} {
		return make([][]byte, 0)
	},
}

func getAddresses() (addr [][]byte) {
	addr = prefetchAddressPool.Get().([][]byte)
	return addr[:0]
}

func putAddresses(addr [][]byte) {
	prefetchAddressPool.Put(addr)
}
