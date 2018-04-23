package main

import (
	"database/sql"
	"testing"

	"github.com/golang/protobuf/proto"
	_ "github.com/lib/pq"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb/web"
)

var (
	testdb *sql.DB
)

const testDbOptions = "dbname=test user=postgres password=passw0rd"

func TestMain(m *testing.M) {
	var err error
	testdb, err = sql.Open("postgres", testDbOptions)
	if err != nil {
		panic(err)
	}
	defer testdb.Close()
	err = db.WithTx(testdb, func(tx *sql.Tx) error {
		db.DestroySchema(tx)
		return db.InitSchema(tx)
	})
	m.Run()
}

func TestNewUser(t *testing.T) {
	req := &web.Request{
		NewUser: proto.String("xyzzy"),
	}
	resp := handleRequest(testdb, req)
	if resp.Error != nil {
		t.Error(*resp.Error)
	}

	var userId, userId1 int64
	userId = resp.GetUserId()
	err := testdb.QueryRow(`SELECT id FROM "user" WHERE id = $1`, userId).Scan(&userId1)
	if err != nil {
		t.Error(err)
	} else if userId != userId1 {
		t.Errorf("%v != %v", userId, userId1)
	}
}

func TestPackages(t *testing.T) {
	var (
		user       db.User
		truck      db.Truck
		pkg1, pkg2 db.Package
	)
	items1 := &db.PackageItems{
		Items: []*db.PackageItem{
			{
				Description: proto.String("all of programming"),
				Count:       proto.Int32(10),
			},
		},
	}
	items2 := &db.PackageItems{
		Items: []*db.PackageItem{
			{
				Description: proto.String("orange chicken"),
				Count:       proto.Int32(1),
			},
			{
				Description: proto.String("broccoli beef"),
				Count:       proto.Int32(1),
			},
		},
	}
	err := db.WithTx(testdb, func(tx *sql.Tx) (err error) {
		err = user.Create(tx)
		if err != nil {
			return
		}
		userId := sql.NullInt64{Valid: true, Int64: int64(user)}

		truck.UpdatePos(tx, db.Coord{X: 0, Y: 0})

		err = pkg1.Create(tx, items1, db.Coord{X: 1, Y: 2}, userId, 1)
		if err != nil {
			return
		}
		_, err = tx.Exec(`UPDATE package SET truck_id = $1 WHERE id = $2`,
			truck, pkg1)
		if err != nil {
			return
		}
		err = pkg2.Create(tx, items2, db.Coord{X: 3, Y: 4}, userId, 1)
		return
	})
	if err != nil {
		t.Fatal(err)
	}

	// requesting a single package
	req := &web.Request{
		GetPackageStatus: proto.Int64(int64(pkg2)),
	}
	resp := handleRequest(testdb, req)
	if resp.Error != nil {
		t.Error(*resp.Error)
	}
	if n := len(resp.GetPackages()); n != 1 {
		t.Fatalf("expect 1 packages, got %d", n)
	}
	if n := len(resp.GetPackages()[0].GetDetail().GetItems()); n != 2 {
		t.Fatalf("expect 2 items, got %d", n)
	}

	// changing destination
	req = &web.Request{
		ChangeDestination: &web.PkgDest{
			PackageId: proto.Int64(int64(pkg1)),
			X:         proto.Int32(99),
			Y:         proto.Int32(100),
		},
	}
	resp = handleRequest(testdb, req)
	if resp.Error != nil {
		t.Error(*resp.Error)
	}
	var dest db.Coord
	err = testdb.QueryRow(`SELECT destination FROM package WHERE id=$1`, pkg1).Scan(&dest)
	if err != nil {
		t.Error(err)
	} else if (dest != db.Coord{X: 99, Y: 100}) {
		t.Errorf("%v != (99,100)", dest)
	}
	err = db.WithTx(testdb, func(tx *sql.Tx) error {
		return truck.UpdateStatus(tx, db.DELIVERING)
	})
	if err != nil {
		t.Fatal(err)
	}
	resp = handleRequest(testdb, req) // cannot change because package is being delivered
	if e := resp.Error; e == nil || *e != errCannotChangeDest.Error() {
		t.Errorf("resp.Error = %v", e)
	}
	err = db.WithTx(testdb, func(tx *sql.Tx) error {
		truck.UpdateStatus(tx, db.IDLE)
		return pkg1.SetDelivered(tx)
	})
	if err != nil {
		t.Fatal(err)
	}
	resp = handleRequest(testdb, req) // cannot change because package is delivered
	if e := resp.Error; e == nil || *e != errCannotChangeDest.Error() {
		t.Errorf("resp.Error = %v", e)
	}

	// requesting all package for a user
	req = &web.Request{
		GetPackages: proto.Int64(int64(user)),
	}
	resp = handleRequest(testdb, req)
	if resp.Error != nil {
		t.Error(*resp.Error)
	}
	if n := len(resp.GetPackages()); n != 2 {
		t.Errorf("expect 2 packages, got %d", n)
	}
}
