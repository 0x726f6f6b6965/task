package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type PageToken interface {
	// GetToken - get token string
	GetToken() string
	// GetID - get ID
	GetID() string
	// GetSize - get page size
	GetSize() int64
	// SetID - set ID
	SetID(id string)
	// SetSize - set page size
	SetSize(size int64)
}

type token struct {
	// Id
	Id string `json:"id"`
	// Size
	Size int64 `json:"size"`
}

// GetID - get ID
func (t *token) GetID() string {
	return t.Id
}

// GetSize - get page size
func (t *token) GetSize() int64 {
	return t.Size
}

// GetToken - encodes an token struct as a page token string.
func (t *token) GetToken() string {
	b, _ := json.Marshal(t)
	encodedStr := base64.URLEncoding.EncodeToString(b)
	return encodedStr
}

// SetID - set ID
func (t *token) SetID(id string) {
	t.Id = id
}

// SetSize - set page size
func (t *token) SetSize(size int64) {
	t.Size = size
}

func NewPageToken(id string, size int64) PageToken {
	return &token{Id: id, Size: size}
}

// GetPageTokenByString - decodes an encoded page token into an token struct
func GetPageTokenByString(in string) (PageToken, error) {
	b, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return nil, fmt.Errorf("decode page token struct: %w", err)
	}
	t := new(token)
	if err = json.Unmarshal(b, t); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return t, nil
}
