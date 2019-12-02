package locker

import (
	"hash"
	"math/big"
	"sync"
)

func Set(capacity int, hash hash.Hash) mxset {
	return mxset{hash: hash, set: make([]sync.Mutex, capacity), div: big.NewInt(int64(capacity))}
}

type mxset struct {
	hash hash.Hash
	set  []sync.Mutex
	div  *big.Int
}

func (mx mxset) ByFingerprint(fingerprint []byte) *sync.Mutex {
	_, _ = mx.hash.Write(fingerprint)
	shard := big.NewInt(0).Mod(big.NewInt(0).SetBytes(fingerprint), mx.div)
	mx.hash.Reset()
	return mx.ByVirtualShard(shard.Uint64())
}

func (mx mxset) ByKey(key string) *sync.Mutex {
	return mx.ByFingerprint([]byte(key))
}

func (mx mxset) ByVirtualShard(shard uint64) *sync.Mutex {
	return &mx.set[shard%uint64(len(mx.set))]
}
