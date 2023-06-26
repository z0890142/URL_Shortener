package common

import (
	"crypto/rand"
	"hash"

	"github.com/spaolacci/murmur3"
)

type HashPool struct {
	pool chan hash.Hash32
}

func NewHashPool(size int) *HashPool {
	pool := make(chan hash.Hash32, size)
	for i := 0; i < size; i++ {
		randomSeed, _ := generateRandomSeed()
		seed := uint32(randomSeed[0])<<24 | uint32(randomSeed[1])<<16 | uint32(randomSeed[2])<<8 | uint32(randomSeed[3])

		hash := murmur3.New32WithSeed(seed)
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

func generateRandomSeed() ([]byte, error) {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, err
	}
	return seed, nil
}
