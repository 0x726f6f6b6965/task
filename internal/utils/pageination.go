package utils

import (
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"strings"
)

type PageToken struct {
	// Id
	Id string
	// Size
	Size int64
}

// String returns a string representation of the page token.
func (p *PageToken) String() string {
	return encodePageTokenStruct(&p)
}

// encodePageTokenStruct encodes an arbitrary struct as a page token.
func encodePageTokenStruct(v interface{}) string {
	var b strings.Builder
	base64Encoder := base64.NewEncoder(base64.URLEncoding, &b)
	gobEncoder := gob.NewEncoder(base64Encoder)
	_ = gobEncoder.Encode(v)
	_ = base64Encoder.Close()
	return b.String()
}

// DecodePageTokenStruct decodes an encoded page token into an arbitrary struct.
func DecodePageTokenStruct(s string, v interface{}) error {
	dec := gob.NewDecoder(base64.NewDecoder(base64.URLEncoding, strings.NewReader(s)))
	if err := dec.Decode(v); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode page token struct: %w", err)
	}
	return nil
}
