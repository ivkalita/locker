package locker_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/fnv"
	"runtime"
	"testing"

	. "github.com/kamilsk/locker"
)

var keys = [...]string{runtime.GOOS, runtime.GOARCH}

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
