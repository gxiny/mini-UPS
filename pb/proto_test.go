package pb

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

var refStream = []byte{0x2, 0x10, 0x1, 0x2, 0x20, 0x1}

func TestWriteProto(t *testing.T) {
	w := bytes.NewBuffer(nil)
	con := ups.Connect{
		NumTrucksInit: proto.Int32(1),
	}
	com := ups.Commands{
		Disconnect: proto.Bool(true),
	}
	var err error
	_, err = WriteProto(w, &con)
	if err != nil {
		t.Error(err)
	}
	_, err = WriteProto(w, &com)
	if err != nil {
		t.Error(err)
	}
	if b := w.Bytes(); !bytes.Equal(b, refStream) {
		t.Errorf("bytes written = %v", b)
	}
}

func TestReadProto(t *testing.T) {
	r := bytes.NewReader(refStream)
	var (
		con ups.Connect
		com ups.Commands
	)
	_, err := ReadProto(r, &con)
	if err != nil {
		t.Error(err)
	}
	if x := con.NumTrucksInit; *x != 1 {
		t.Errorf("NumTrucksInit = %d", *x)
	}
	_, err = ReadProto(r, &com)
	if err != nil {
		t.Error(err)
	}
	if x := com.Disconnect; *x != true {
		t.Errorf("Disconnect = %v", *x)
	}
}
