// Package world encapsulates the communication with the world simulator.
package world

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"reflect"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb"
)

// Sim represents a world simulator.
type Sim struct {
	r *bufio.Reader
	w io.WriteCloser
}

// Connected is the type of message that the world simulator will respond after sending a "connect" message to it.
// It is a generated protobuf type and should have an "error" field.
type Connected interface {
	proto.Message
	GetError() string
}

// Connect connects to the given address, sends the connect message
// to the simulator and receives the connected message.
func (w *Sim) Connect(addr string, connect proto.Message, connected Connected) (err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	w.r = bufio.NewReader(conn)
	w.w = conn
	defer func() {
		if err != nil {
			w.r = nil
			w.w = nil
			conn.Close()
		}
	}()
	err = w.WriteProto(connect)
	if err != nil {
		return
	}
	err = w.ReadProto(connected)
	if e := connected.GetError(); e != "" {
		err = errors.New(e)
	}
	if err != nil {
		return
	}
	return
}

// Close closes the connection to the simulator.
// To properly disconnect from the simulator (not crash it),
// send a message which contains "disconnect:true"
// and make sure you have received "finished:true" before calling Close.
func (w *Sim) Close() error {
	return w.w.Close()
}

// ReadProto receives a single protobuf message.
func (w *Sim) ReadProto(msg proto.Message) error {
	_, err := pb.ReadProto(w.r, msg)
	return err
}

// WriteProto writes a single protobuf message.
func (w *Sim) WriteProto(msg proto.Message) error {
	_, err := pb.WriteProto(w.w, msg)
	return err
}

// Responses is the type of message that the world simulator will respond after sending a "commands" message to it.
// It is a generated protobuf type and should have a "finished" field.
type Responses interface {
	proto.Message
	GetFinished() bool
}

var responsesType = reflect.TypeOf((*Responses)(nil)).Elem()

// ReadChan returns a channel by reading which one can receive
// all messages from the world simulator.
// The channel is closed after a message with "finished:true"
// is received.  Therefore, one can safely close the connection
// after reading all messages from the channel.
// Don't call ReadChan more than once on a Sim object.
func (w *Sim) ReadChan(typ reflect.Type) <-chan proto.Message {
	if !typ.Implements(responsesType) {
		panic("response type does not implement Responses")
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem() // should also implment that interface
	}
	ch := make(chan proto.Message)
	go func() {
		defer close(ch)
		for {
			val := reflect.New(typ).Interface().(Responses)
			err := w.ReadProto(val)
			if err != nil {
				log.Println(err)
				if err == io.EOF {
					break
				}
			} else {
				ch <- val
				if val.GetFinished() {
					break
				}
			}
		}
	}()
	return ch
}
