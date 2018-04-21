package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb"
	"gitlab.oit.duke.edu/rz78/ups/pb/amz"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
	"gitlab.oit.duke.edu/rz78/ups/world"
)

var (
	worldAddr = flag.String("sim", ":23456", "world simulator address")
	upsAddr   = flag.String("ups", ":23333", "ups server address")
)

var responsesType = proto.MessageType("amz.Responses")

var (
	sim      *world.Sim
	simSpeed *uint32
)

func main() {
	flag.Parse()
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		msg := ParseProto(sc.Text())
		if msg == nil {
			continue
		}
		var err error
		if sim == nil { // not connected
			msg, ok := msg.(*amz.Connect)
			if ok {
				err = connect(msg)
			} else {
				err = fmt.Errorf("not connected to world")
			}
		} else {
			switch msg.(type) {
			case *amz.Commands:
				err = tellWorld(msg)
			case *bridge.UCommands:
				err = tellUPS(msg)
			default:
				err = fmt.Errorf("unknown message %T", msg)
			}
		}
		if err != nil {
			log.Println("error:", err)
		}
	}
}

func connect(msg proto.Message) error {
	var w world.Sim
	connected := new(amz.Connected)
	err := w.Connect(*worldAddr, msg, connected)
	if err != nil {
		return err
	}
	ch := w.ReadChan(responsesType)
	go func() {
		for msg := range ch {
			log.Println("world:", msg)
		}
	}()
	sim = &w
	return nil
}

func tellWorld(msg proto.Message) error {
	c := msg.(*amz.Commands)
	if c.SimSpeed != nil {
		simSpeed = c.SimSpeed
	} else {
		c.SimSpeed = simSpeed
		err := sim.WriteProto(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func tellUPS(msg proto.Message) error {
	conn, err := net.Dial("tcp", *upsAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = pb.WriteProto(conn, msg)
	if err != nil {
		return err
	}
	r := bufio.NewReader(conn)
	resp := new(bridge.UResponses)
	_, err = pb.ReadProto(r, resp)
	if err != nil {
		return err
	}
	log.Println(resp)
	return nil
}
