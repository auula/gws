package sessionx

import (
	"bytes"
	"encoding/gob"
)

func Encoder(v interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	ecoder := gob.NewEncoder(&buffer)
	if err := ecoder.Encode(v); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Decoder(v []byte, obj interface{}) (interface{}, error) {
	reader := bytes.NewReader(v)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(obj); err != nil {
		return nil, err
	}
	return obj, nil
}
