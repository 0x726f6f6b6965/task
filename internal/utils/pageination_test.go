package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeToken(t *testing.T) {
	token := NewPageToken("test-id", 225)
	tokenStr := token.GetToken()
	newToken, err := GetPageTokenByString(tokenStr)
	assert.Nil(t, err)
	assert.Equal(t, token.GetID(), newToken.GetID())
	assert.Equal(t, token.GetSize(), newToken.GetSize())
}

func TestSetToken(t *testing.T) {
	token := NewPageToken("test-id", 225)
	tokenStr := token.GetToken()
	token.SetID("test-id2")
	token.SetSize(60)
	newTokenStr := token.GetToken()
	assert.NotEqual(t, newTokenStr, tokenStr)
	assert.NotEqual(t, "test-id", token.GetID())
	assert.NotEqual(t, 225, token.GetSize())
}
