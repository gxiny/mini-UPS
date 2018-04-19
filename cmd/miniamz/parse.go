package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/pb/amz"
)

type Scanner struct {
	scanner.Scanner
}

func (s *Scanner) Init(r io.Reader) {
	s.Scanner.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanStrings
	s.Error = func(_ *scanner.Scanner, msg string) {
		fmt.Fprintln(os.Stderr, msg)
	}
}

func (s *Scanner) err(msg string) {
	s.Error(&s.Scanner, msg)
	s.ErrorCount++
}

func (s *Scanner) Expect(expected rune) (token string) {
	scanned := s.Scan()
	token = s.TokenText()
	if scanned != expected {
		s.err(fmt.Sprintf("near %s: expect %s, got %s",
			token,
			scanner.TokenString(expected),
			scanner.TokenString(scanned)))
	}
	return
}

func (s *Scanner) ScanInt(bitSize int) int64 {
	token := s.Expect(scanner.Int)
	if token != "" {
		n, err := strconv.ParseInt(token, 10, bitSize)
		if err != nil {
			s.err("ParseInt fail: " + token)
		}
		return n
	}
	return 0
}

func (s *Scanner) ScanIdent() string {
	return s.Expect(scanner.Ident)
}

func (s *Scanner) ScanString() string {
	token := s.Expect(scanner.String)
	if len(token) >= 2 {
		token = token[1 : len(token)-1] // strip ""
	}
	return token
}

var parserMap = map[string]func(*Scanner) proto.Message{
	"connect":    parseConnect,
	"disconnect": parseDisconnect,
	"simspeed":   parseSimSpeed,
	"purchase":   parsePurchase,
}

func ParseProto(s string) proto.Message {
	var sc Scanner
	sc.Init(strings.NewReader(s))

	name := sc.ScanIdent()
	if sc.ErrorCount > 0 {
		return nil
	}
	parser, ok := parserMap[name]
	if !ok {
		sc.err("unknown message: " + name)
		return nil
	}
	return parser(&sc)
}

// syntax: "connect" world_id {x y}
func parseConnect(sc *Scanner) proto.Message {
	worldId := sc.ScanInt(64)
	if sc.ErrorCount > 0 {
		return nil
	}
	msg := &amz.Connect{
		WorldId: &worldId,
	}
	for {
		if sc.Peek() == scanner.EOF {
			break
		}
		x := sc.ScanInt(32)
		y := sc.ScanInt(32)
		if sc.ErrorCount > 0 {
			return nil
		}
		msg.InitWarehouses = append(msg.InitWarehouses, &amz.InitWarehouse{
			X: proto.Int32(int32(x)),
			Y: proto.Int32(int32(y)),
		})
	}
	return msg
}

func parseDisconnect(*Scanner) proto.Message {
	return &amz.Commands{
		Disconnect: proto.Bool(true),
	}
}

func parseSimSpeed(sc *Scanner) proto.Message {
	speed := sc.ScanInt(32)
	if sc.ErrorCount > 0 {
		return nil
	}
	return &amz.Commands{
		SimSpeed: proto.Uint32(uint32(speed)),
	}
}

// syntax: "purchase" wh_id {prod_id description count}
func parsePurchase(sc *Scanner) proto.Message {
	whId := sc.ScanInt(32)
	if sc.ErrorCount > 0 {
		return nil
	}
	msg := &amz.PurchaseMore{
		WarehouseId: proto.Int32(int32(whId)),
	}
	for {
		if sc.Peek() == scanner.EOF {
			break
		}
		var (
			prId  = sc.ScanInt(64)
			desc  = sc.ScanString()
			count = sc.ScanInt(32)
		)
		if sc.ErrorCount > 0 {
			return nil
		}
		msg.Things = append(msg.Things, &amz.Product{
			Id:          &prId,
			Description: &desc,
			Count:       proto.Int32(int32(count)),
		})
	}
	return &amz.Commands{
		Buy: []*amz.PurchaseMore{msg},
	}
}
