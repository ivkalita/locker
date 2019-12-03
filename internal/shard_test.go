package internal_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/fnv"
	"math"
	"runtime"
	"testing"

	. "github.com/kamilsk/locker/internal"
)

var keys = [...]string{runtime.GOOS, runtime.GOARCH}

func TestShardNumberCalculation(t *testing.T) {
	sources := map[string]hash.Hash{
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
	size := uint64(math.MaxUint32)
	for name, checksum := range sources {
		t.Run(name, func(t *testing.T) {
			for _, key := range keys {
				checksum.Write([]byte(key))
				in := checksum.Sum(nil)
				checksum.Reset()

				x, y := ShardNumberNaive(in, size), ShardNumberSimple(in, size)
				if x != y {
					t.Errorf("%d != %d", x, y)
					t.FailNow()
				}

				z := ShardNumberFast(in, size)
				if x != z {
					t.Errorf("%d != %d", x, z)
					t.FailNow()
				}
			}
		})
	}
}

// BenchmarkShardNumberCalculation/naive,md5:darwin-4         	 1377596	       818 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,md5:amd64-4          	 1392057	       845 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha1:darwin-4        	 1225707	       957 ns/op	     264 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha1:amd64-4         	 1253956	       950 ns/op	     264 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha256:darwin-4      	  846073	      1284 ns/op	     296 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha256:amd64-4       	  819980	      1295 ns/op	     296 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sha512:darwin-4      	  520620	      2729 ns/op	     552 B/op	       9 allocs/op
// BenchmarkShardNumberCalculation/naive,sha512:amd64-4       	  449974	      2840 ns/op	     552 B/op	       9 allocs/op
// BenchmarkShardNumberCalculation/naive,sum32:darwin-4       	 2967585	       393 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive,sum32:amd64-4        	 3158434	       378 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive,sum32a:darwin-4      	 3147517	       375 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive,sum32a:amd64-4       	 3152670	       380 ns/op	      72 B/op	       6 allocs/op
// BenchmarkShardNumberCalculation/naive,sum64:darwin-4       	 1737346	       713 ns/op	     144 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum64:amd64-4        	 1649262	       799 ns/op	     144 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum64a:darwin-4      	 1982733	       554 ns/op	     144 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum64a:amd64-4       	 2150758	       554 ns/op	     144 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum128:darwin-4      	 1463527	       819 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum128:amd64-4       	 1326205	      1014 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum128a:darwin-4     	 1332231	       914 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/naive,sum128a:amd64-4      	 1346344	      1084 ns/op	     216 B/op	       8 allocs/op
// BenchmarkShardNumberCalculation/simple,md5:darwin-4        	 4558246	       397 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,md5:amd64-4         	 3068404	       396 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha1:darwin-4       	 4375424	       486 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha1:amd64-4        	 4332129	       401 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha256:darwin-4     	 3926182	       297 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha256:amd64-4      	 4037755	       284 ns/op	     144 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha512:darwin-4     	 2294641	       442 ns/op	     208 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sha512:amd64-4      	 2610360	       486 ns/op	     208 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum32:darwin-4      	 8728225	       158 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum32:amd64-4       	 7804992	       136 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum32a:darwin-4     	 9531890	       122 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum32a:amd64-4      	 9003367	       122 ns/op	      24 B/op	       3 allocs/op
// BenchmarkShardNumberCalculation/simple,sum64:darwin-4      	 7476349	       155 ns/op	      32 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum64:amd64-4       	 7506585	       174 ns/op	      32 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum64a:darwin-4     	 6965594	       159 ns/op	      32 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum64a:amd64-4      	 6831498	       157 ns/op	      32 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum128:darwin-4     	 5479737	       215 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum128:amd64-4      	 5423877	       269 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum128a:darwin-4    	 5247582	       338 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/simple,sum128a:amd64-4     	 5239660	       214 ns/op	     112 B/op	       4 allocs/op
// BenchmarkShardNumberCalculation/fast,md5:darwin-4          	 4624431	       256 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,md5:amd64-4           	 4640300	       255 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha1:darwin-4         	 3090070	       366 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha1:amd64-4          	 3723790	       339 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha256:darwin-4       	 2247442	       532 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha256:amd64-4        	 2209219	       581 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha512:darwin-4       	 1042747	      1158 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sha512:amd64-4        	 1000000	      1105 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum32:darwin-4        	15295339	        70.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum32:amd64-4         	16811192	        67.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum32a:darwin-4       	17580582	        66.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum32a:amd64-4        	17699610	        66.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum64:darwin-4        	 9106179	       129 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum64:amd64-4         	 9121912	       129 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum64a:darwin-4       	 9186740	       129 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum64a:amd64-4        	 9201374	       129 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum128:darwin-4       	 4655198	       256 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum128:amd64-4        	 4619257	       256 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum128a:darwin-4      	 4696797	       256 ns/op	       0 B/op	       0 allocs/op
// BenchmarkShardNumberCalculation/fast,sum128a:amd64-4       	 4683013	       257 ns/op	       0 B/op	       0 allocs/op
func BenchmarkShardNumberCalculation(b *testing.B) {
	hashes := []struct {
		name     string
		checksum hash.Hash
	}{
		{"md5", md5.New()},
		{"sha1", sha1.New()},
		{"sha256", sha256.New()},
		{"sha512", sha512.New()},
		{"sum32", fnv.New32()},
		{"sum32a", fnv.New32a()},
		{"sum64", fnv.New64()},
		{"sum64a", fnv.New64a()},
		{"sum128", fnv.New128()},
		{"sum128a", fnv.New128a()},
	}
	benchmarks := []struct {
		name           string
		implementation func([]byte, uint64) uint64
	}{
		{"naive", ShardNumberNaive},
		{"simple", ShardNumberSimple},
		{"fast", ShardNumberFast},
	}
	size := uint64(math.MaxUint32)
	for _, benchmark := range benchmarks {
		for _, algorithm := range hashes {
			for _, key := range keys {
				b.Run(benchmark.name+","+algorithm.name+":"+key, func(b *testing.B) {
					algorithm.checksum.Write([]byte(key))
					in := algorithm.checksum.Sum(nil)
					algorithm.checksum.Reset()

					b.ResetTimer()
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						_ = benchmark.implementation(in, size)
					}
				})
			}
		}
	}
}
