package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenAppKey(t *testing.T) {
	appKey := GenAppKey()
	t.Log(appKey)
	assert.Equal(t, 16, len(GenAppKey()), "length must be 16")
}

func TestGenSecretKey(t *testing.T) {
	secretKey := GenSecretKey()
	t.Log(secretKey)
	assert.Equal(t, 32, len(secretKey), "length must be 16")
}
