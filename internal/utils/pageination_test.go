package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeToken(t *testing.T) {
	token := &PageToken{
		Id:   "test-id",
		Size: 225,
	}
	tokenStr := token.String()
	newToken := &PageToken{}
	err := DecodePageTokenStruct(tokenStr, newToken)
	assert.Nil(t, err)
	assert.Equal(t, token.Id, newToken.Id)
	assert.Equal(t, token.Size, newToken.Size)
}
