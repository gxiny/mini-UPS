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
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
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
	scanned := s.Scan()
	token := s.TokenText()
	negative := false
	if token == "-" {
		negative = true
		scanned = s.Scan()
		token = s.TokenText()
	}
	if scanned == scanner.Int {
		n, err := strconv.ParseInt(token, 10, bitSize)
		if err != nil {
			s.err("ParseInt fail: " + token)
		}
		if negative {
			return -n
		}
		return n
	} else {
		s.err(fmt.Sprintf("near %s: expect %s, got %s",
			token,
			scanner.TokenString(scanner.Int),
			scanner.TokenString(scanned)))
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
	// world commands
	"connect":    parseConnect,
	"disconnect": parseDisconnect,
	"simspeed":   parseSimSpeed,
	"purchase":   parsePurchase,
	"pack":       parsePack,
	"put":        parsePutOnTruck,
	// UPS commands
	"pkg":      parsePackage,
	"truckreq": parseTruckReq,
	"loaded":   parseLoaded,
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

// syntax: "disconnect"
func parseDisconnect(*Scanner) proto.Message {
	return &amz.Commands{
		Disconnect: proto.Bool(true),
	}
}

// syntax: "pack" warehouse_id ship_id <product_list>
func parsePack(sc *Scanner) proto.Message {
	whId := sc.ScanInt(32)
	shId := sc.ScanInt(64)
	msg := &amz.Pack{
		WarehouseId: proto.Int32(int32(whId)),
		ShipId:      &shId,
		Things:      productList(sc),
	}
	if sc.ErrorCount > 0 {
		return nil
	}
	return &amz.Commands{
		ToPack: []*amz.Pack{msg},
	}
}

// syntax: "put" warehouse_id truck_id ship_id
func parsePutOnTruck(sc *Scanner) proto.Message {
	var (
		whId = sc.ScanInt(32)
		trId = sc.ScanInt(32)
		shId = sc.ScanInt(64)
	)
	if sc.ErrorCount > 0 {
		return nil
	}
	msg := &amz.PutOnTruck{
		WarehouseId: proto.Int32(int32(whId)),
		TruckId:     proto.Int32(int32(trId)),
		ShipId:      &shId,
	}
	return &amz.Commands{
		Load: []*amz.PutOnTruck{msg},
	}
}

// product_list ::= {prod_id description count}
func productList(sc *Scanner) []*amz.Product {
	l := []*amz.Product{}
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
			break
		}
		l = append(l, &amz.Product{
			Id:          &prId,
			Description: &desc,
			Count:       proto.Int32(int32(count)),
		})
	}
	return l
}

// syntax: "simspeed" speed
func parseSimSpeed(sc *Scanner) proto.Message {
	speed := sc.ScanInt(32)
	if sc.ErrorCount > 0 {
		return nil
	}
	return &amz.Commands{
		SimSpeed: proto.Uint32(uint32(speed)),
	}
}

// syntax: "purchase" wh_id <product_list>
func parsePurchase(sc *Scanner) proto.Message {
	whId := sc.ScanInt(32)
	msg := &amz.PurchaseMore{
		WarehouseId: proto.Int32(int32(whId)),
		Things:      productList(sc),
	}
	if sc.ErrorCount > 0 {
		return nil
	}
	return &amz.Commands{
		Buy: []*amz.PurchaseMore{msg},
	}
}

// syntax: "pkg" warehouse_id ups_user_id x y {item_id description amount}
// -1 means omitted for ups_user_id
func parsePackage(sc *Scanner) proto.Message {
	var (
		whId  = sc.ScanInt(32)
		upsId = sc.ScanInt(64)
		x     = sc.ScanInt(32)
		y     = sc.ScanInt(32)
	)
	if sc.ErrorCount > 0 {
		return nil
	}
	pkg := &bridge.Package{
		WarehouseId: proto.Int32(int32(whId)),
		X:           proto.Int32(int32(x)),
		Y:           proto.Int32(int32(y)),
	}
	if upsId != -1 {
		pkg.UpsUserId = proto.Int64(upsId)
	}
	for {
		if sc.Peek() == scanner.EOF {
			break
		}
		var (
			itId   = sc.ScanInt(64)
			desc   = sc.ScanString()
			amount = sc.ScanInt(32)
		)
		if sc.ErrorCount > 0 {
			return nil
		}
		pkg.Items = append(pkg.Items, &bridge.Item{
			ItemId:      &itId,
			Description: &desc,
			Amount:      proto.Int32(int32(amount)),
		})
	}
	return &bridge.ACommands{
		PackageIdReq: pkg,
	}
}

// syntax: "truckreq" warehouse_id x y
func parseTruckReq(sc *Scanner) proto.Message {
	whId := sc.ScanInt(32)
	x := sc.ScanInt(32)
	y := sc.ScanInt(32)
	if sc.ErrorCount > 0 {
		return nil
	}
	return &bridge.ACommands{
		TruckReq: &bridge.RequestTruck{
			WarehouseId: proto.Int32(int32(whId)),
			X: proto.Int32(int32(x)),
			Y: proto.Int32(int32(y)),
		},
	}
}

// syntax: "loaded" truck_id warehouse_id {package_id}
func parseLoaded(sc *Scanner) proto.Message {
	trId := sc.ScanInt(32)
	whId := sc.ScanInt(32)
	msg := &bridge.PackagesLoaded{
		TruckId:     proto.Int32(int32(trId)),
		WarehouseId: proto.Int32(int32(whId)),
	}
	for {
		if sc.Peek() == scanner.EOF {
			break
		}
		pkgId := sc.ScanInt(64)
		if sc.ErrorCount > 0 {
			return nil
		}
		msg.PackageIds = append(msg.PackageIds, pkgId)
	}
	return &bridge.ACommands{
		Loaded: msg,
	}
}
