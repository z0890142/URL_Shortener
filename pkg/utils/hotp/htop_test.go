package hotp

import (
	"URL_Shortener/pkg/utils/algorithm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHotpGenerateCode(t *testing.T) {
	type args struct {
		secret  string
		counter uint64
		opts    GenerateOptions
	}

	code, err := HotpGenerateCode("127.0.0.1", 1, GenerateOptions{
		Algorithm: algorithm.AlgorithmSHA1,
		Digits:    6,
	})
	assert.Nil(t, err)
	assert.NotNil(t, code)

	code, err = HotpGenerateCode("", 1, GenerateOptions{
		Algorithm: algorithm.AlgorithmSHA1,
		Digits:    6,
	})
	assert.NotNil(t, err)

}
