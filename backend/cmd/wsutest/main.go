// Command wsutest tests the communication between the
// world simulator and ups server.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb"
	"gitlab.oit.duke.edu/rz78/ups/pb/ups"
)

var (
	simAddr       = flag.String("sim", ":12345", "address of the world simulator")
	worldId       = flag.Uint("world", 0, "ID of world to connect (0 = create a new world)")
	numTrucksInit = flag.Int("truck", 1, "number of initial trucks")
)

type rw struct {
	io.Writer
	*bufio.Reader
}

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *simAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	rw := &rw{conn, bufio.NewReader(conn)}
	if *worldId == 0 {
		createWorld(rw, *numTrucksInit)
	} else {
		connectWorld(rw, *worldId)
	}
	disconnect(rw)
}

func createWorld(conn *rw, numTrucksInit int) {
	c := &ups.Connect{
		NumTrucksInit: proto.Int32(int32(numTrucksInit)),
	}
	_, err := pb.WriteProto(conn, c)
	if err != nil {
		log.Fatal(err)
	}
	r := new(ups.Connected)
	_, err = pb.ReadProto(conn, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
	for i := 0; i < numTrucksInit; i++ {
		r := new(ups.Responses)
		_, err = pb.ReadProto(conn, r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(r)
	}
}

func connectWorld(conn *rw, worldId uint) {
	c := &ups.Connect{
		ReconnectId: proto.Int64(int64(worldId)),
	}
	_, err := pb.WriteProto(conn, c)
	if err != nil {
		log.Fatal(err)
	}
	r := new(ups.Connected)
	_, err = pb.ReadProto(conn, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
}

func disconnect(conn *rw) {
	c := &ups.Commands{
		Disconnect: proto.Bool(true),
	}
	_, err := pb.WriteProto(conn, c)
	if err != nil {
		log.Fatal(err)
	}
	r := new(ups.Responses)
	_, err = pb.ReadProto(conn, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
}
