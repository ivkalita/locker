package internal

import (
	"encoding/hex"
	"math"
	"math/big"
)

func ShardNumberNaive(in []byte, size uint64) uint64 {
	base, success := big.NewInt(0).SetString(hex.EncodeToString(in), 16)
	if !success {
		panic("invalid input")
	}
	return big.NewInt(0).Mod(base, big.NewInt(int64(size))).Uint64()
}

func ShardNumberSimple(in []byte, size uint64) uint64 {
	return big.NewInt(0).Mod(big.NewInt(0).SetBytes(in), big.NewInt(int64(size))).Uint64()
}

func ShardNumberFast(in []byte, size uint64) uint64 {
	var shard uint64
	for f, i := uint64(1), len(in)-1; i >= 0; i-- {
		shard = (shard + uint64(in[i])*f) % size
		f = (f * (math.MaxUint8 + 1)) % size
	}
	return shard
}
