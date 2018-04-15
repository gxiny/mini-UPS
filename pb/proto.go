// Package pb contains packages generated from Google Protocol Buffer
// and provides convenience functions for reading/writing protobuf.
package pb

//go:generate protoc --go_out=ups ups.proto

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)

// WriteProto writes a proto.Message to an Writer,
// with a Varint32-prefixed length.
func WriteProto(w io.Writer, message proto.Message) (n int, err error) {
	buf := proto.NewBuffer(nil)
	err = buf.EncodeMessage(message)
	if err != nil {
		return
	}
	n, err = w.Write(buf.Bytes())
	return
}

type Reader interface {
	io.Reader
	io.ByteReader
}

// ReadProto reads a proto.Message from an Reader.
// The message is prefixed by a Varint32-prefixed length
// (as written by WriteProto) so that this function can know
// how much to read from a stream.
func ReadProto(r Reader, message proto.Message) (n int, err error) {
	size, err := binary.ReadUvarint(r)
	if err != nil {
		return
	}
	b := make([]byte, size)
	n, err = io.ReadFull(r, b)
	n += int(size)
	if err != nil {
		return
	}
	err = proto.Unmarshal(b, message)
	return
}
