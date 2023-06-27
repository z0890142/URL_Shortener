package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase62Encode(t *testing.T) {
	code := Base62Encode(2778191738)
	assert.Equal(t, "3210Aa", code)
}
