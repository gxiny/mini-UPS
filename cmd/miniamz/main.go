package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/world"
	"gitlab.oit.duke.edu/rz78/ups/pb/amz"
)

var worldAddr = flag.String("sim", ":23456", "world simulator address")

var responsesType = proto.MessageType("amz.Responses")

func main() {
	var sim *world.Sim
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		msg := ParseProto(sc.Text())
		if msg == nil {
			continue
		}
		var err error
		if sim == nil { // not connected
			connected := new(amz.Connected)
			sim, err = world.Connect(*worldAddr, msg, connected)
			if err != nil {
				log.Println("error:", err)
				continue
			}
			defer sim.Close()
			ch := sim.ReadChan(responsesType)
			go func() {
				for msg := range ch {
					log.Println("world:", msg)
				}
			}()
		} else {
			err = sim.WriteProto(msg)
			if err != nil {
				log.Println("error:", err)
				continue
			}
		}
	}
}
