package common

import (
	"hash"

	"github.com/spaolacci/murmur3"
)

type HashPool struct {
	pool chan hash.Hash32
}

func NewHashPool(size int) *HashPool {
	pool := make(chan hash.Hash32, size)
	for i := 0; i < size; i++ {
		hash := murmur3.New32()
		pool <- hash
	}
	return &HashPool{pool: pool}
}

func (p *HashPool) GetHash() hash.Hash32 {
	return <-p.pool
}

func (p *HashPool) ReleaseHash(hash hash.Hash32) {
	hash.Reset()
	p.pool <- hash
}

func (p *HashPool) Close() {
	close(p.pool)
}
