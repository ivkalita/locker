package locker_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/fnv"
	"math/big"
	"math/rand"
	"runtime"
	"testing"
	"time"

	. "github.com/kamilsk/locker"
)

var keys = [...]string{runtime.GOOS, runtime.GOARCH}

func naive(in []byte, div uint64) uint64 {
	base, _ := big.NewInt(0).SetString(hex.EncodeToString(in), 16)
	return big.NewInt(0).Mod(base, big.NewInt(int64(div))).Uint64()
}

func simple(in []byte, div uint64) uint64 {
	return big.NewInt(0).Mod(big.NewInt(0).SetBytes(in), big.NewInt(int64(div))).Uint64()
}

func optimized(in []byte, div uint64) uint64 {
	var r uint64
	for i, m := len(in)-1, uint64(1); i >= 0; i-- {
		r = (r + uint64(in[i])*m) % div
		m = (m * 256) % div
	}
	return r
}

func TestCalculation(t *testing.T) {
	tests := map[string]hash.Hash{
		"md5":     md5.New(),
		"sha1":    sha1.New(),
		"sha256":  sha256.New(),
		"sha512":  sha512.New(),
		"sum32":   fnv.New32(),
		"sum32a":  fnv.New32a(),
		"sum64":   fnv.New64(),
		"sum64a":  fnv.New64a(),
		"sum128":  fnv.New128(),
		"sum128a": fnv.New128a(),
	}
	div := rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
	div = 112
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for _, key := range keys {
				test.Write([]byte(key))
				x, y, z := naive(test.Sum(nil), div), simple(test.Sum(nil), div), optimized(test.Sum(nil), div)
				test.Reset()
				fmt.Println(x, y, z)
			}
		})
	}
}

func BenchmarkCalculation(b *testing.B) {
	hashes := map[string]hash.Hash{
		"md5":     md5.New(),
		"sha1":    sha1.New(),
		"sha256":  sha256.New(),
		"sha512":  sha512.New(),
		"sum32":   fnv.New32(),
		"sum32a":  fnv.New32a(),
		"sum64":   fnv.New64(),
		"sum64a":  fnv.New64a(),
		"sum128":  fnv.New128(),
		"sum128a": fnv.New128a(),
	}
	benchmarks := map[string]func(in []byte, div uint64) uint64{
		"naive":     naive,
		"simple":    simple,
		"optimized": optimized,
	}
	div := rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
	div = 112
	for name, bm := range benchmarks {
		for fn, algorithm := range hashes {
			b.Run(name+", "+fn, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					for _, key := range keys {
						algorithm.Write([]byte(key))
						_ = bm(algorithm.Sum(nil), div)
						algorithm.Reset()
					}
				}
			})
		}
	}
}

func TestSet(t *testing.T) {
	set := Set(3, md5.New())
	set.ByKey("test").Lock()
	set.ByKey("another").Lock()
	defer set.ByKey("test").Unlock()
	defer set.ByKey("another").Unlock()
	go func() { set.ByKey("test").Unlock() }()
	go func() { set.ByKey("another").Unlock() }()
	set.ByKey("test").Lock()
	set.ByKey("another").Lock()
}

// BenchmarkSet/md5-4         	 3252902	       378 ns/op	      56 B/op	       7 allocs/op
// BenchmarkSet/sha1-4        	 2727148	       397 ns/op	      56 B/op	       7 allocs/op
// BenchmarkSet/sha256-4      	 2518604	       448 ns/op	      56 B/op	       7 allocs/op
// BenchmarkSet/sha512-4      	 2896508	       380 ns/op	      56 B/op	       7 allocs/op
// BenchmarkSet/sum32-4       	 3305530	       351 ns/op	      56 B/op	       7 allocs/op
func BenchmarkSet(b *testing.B) {
	benchmarks := []struct {
		name string
		hash hash.Hash
	}{
		{name: "md5", hash: md5.New()},
		{name: "sha1", hash: sha1.New()},
		{name: "sha256", hash: sha256.New()},
		{name: "sha512", hash: sha512.New()},
		{name: "sum32", hash: fnv.New32()},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()

			set := Set(10, bm.hash)
			for i := 0; i < b.N; i++ {
				for _, key := range keys {
					_ = set.ByKey(key)
				}
			}
		})
	}
}
