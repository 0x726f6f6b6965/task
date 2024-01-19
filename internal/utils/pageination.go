package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type PageToken struct {
	// Id
	Id string `json:"id"`
	// Size
	Size int64 `json:"size"`
}

// String returns a string representation of the page token.
func (p *PageToken) String() string {
	return encodePageTokenStruct(&p)
}

// encodePageTokenStruct encodes an arbitrary struct as a page token.
func encodePageTokenStruct(v interface{}) string {
	b, _ := json.Marshal(v)
	encodedStr := base64.URLEncoding.EncodeToString(b)
	return encodedStr
}

// DecodePageTokenStruct decodes an encoded page token into an arbitrary struct.
func DecodePageTokenStruct(s string, v interface{}) error {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("decode page token struct: %w", err)
	}
	if err = json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	return nil
}
