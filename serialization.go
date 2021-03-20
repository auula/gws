package sessionx

import (
	"bytes"
	"encoding/gob"
)

func encoder(s *Session) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(s); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func decoder(v []byte, s *Session) error {
	reader := bytes.NewReader(v)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(s); err != nil {
		return err
	}
	return nil
}
