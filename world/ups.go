package world

import (
	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

// SimU provides an interface to the world simulator for UPS.
type SimU struct {
	Sim
}

// NewWorld creates a new world with the given number of trucks.
func (w *SimU) NewWorld(addr string, initTrucks int32) (id int64, err error) {
	connect := &ups.Connect{
		NumTrucksInit: &initTrucks,
	}
	connected := new(ups.Connected)
	err = w.Connect(addr, connect, connected)
	if err != nil {
		return
	}
	id = connected.GetWorldId()
	return
}

// ReconnectWorld reconnects to the world with the given ID.
func (w *SimU) ReconnectWorld(addr string, id int64) (err error) {
	connect := &ups.Connect{
		ReconnectId: &id,
	}
	connected := new(ups.Connected)
	err = w.Connect(addr, connect, connected)
	if err != nil {
		return
	}
	if connected.GetWorldId() != id {
		panic("world_id != reconnect_id")
	}
	return
}

// DisconnectWorld sends a disconnect message to the world.
func (w *SimU) DisconnectWorld() error {
	c := &ups.Commands{
		Disconnect: proto.Bool(true),
	}
	return w.WriteProto(c)
}
