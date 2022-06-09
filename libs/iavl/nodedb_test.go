package iavl

import (
	"encoding/binary"
	"math/rand"
	"testing"

	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/stretchr/testify/require"
)

func BenchmarkNodeKey(b *testing.B) {
	ndb := &nodeDB{}
	hashes := makeHashes(b, 2432325)
	for i := 0; i < b.N; i++ {
		ndb.nodeKey(hashes[i])
	}
}

func BenchmarkOrphanKey(b *testing.B) {
	ndb := &nodeDB{}
	hashes := makeHashes(b, 2432325)
	for i := 0; i < b.N; i++ {
		ndb.orphanKey(1234, 1239, hashes[i])
	}
}

func makeHashes(b *testing.B, seed int64) [][]byte {
	b.StopTimer()
	rnd := rand.NewSource(seed)
	hashes := make([][]byte, b.N)
	hashBytes := 8 * ((hashSize + 7) / 8)
	for i := 0; i < b.N; i++ {
		hashes[i] = make([]byte, hashBytes)
		for b := 0; b < hashBytes; b += 8 {
			binary.BigEndian.PutUint64(hashes[i][b:b+8], uint64(rnd.Int63()))
		}
		hashes[i] = hashes[i][:hashSize]
	}
	b.StartTimer()
	return hashes
}

// sink is kept as a global to ensure that value checks and assignments to it can't be
// optimized away, and this will help us ensure that benchmarks successfully run.
var sink interface{}

func BenchmarkConvertLeafOp(b *testing.B) {
	var versions = []int64{
		0,
		1,
		100,
		127,
		128,
		1 << 29,
		-0,
		-1,
		-100,
		-127,
		-128,
		-1 << 29,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, version := range versions {
			sink = convertLeafOp(version)
		}
	}
	if sink == nil {
		b.Fatal("Benchmark wasn't run")
	}
	sink = nil
}

func BenchmarkTreeString(b *testing.B) {
	tree := makeAndPopulateMutableTree(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = tree.String()
	}

	if sink == nil {
		b.Fatal("Benchmark did not run")
	}
	sink = (interface{})(nil)
}

func makeAndPopulateMutableTree(tb testing.TB) *MutableTree {
	memDB := dbm.NewMemDB()
	tree, err := NewMutableTreeWithOpts(memDB, 0, &Options{InitialVersion: 9})
	require.NoError(tb, err)

	for i := 0; i < 1e4; i++ {
		buf := make([]byte, 0, (i/255)+1)
		for j := 0; 1<<j <= i; j++ {
			buf = append(buf, byte((i>>j)&0xff))
		}
		tree.Set(buf, buf)
	}
	_, _, _, err = tree.SaveVersion(false)
	require.Nil(tb, err, "Expected .SaveVersion to succeed")
	return tree
}
