package db

import (
	"database/sql"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestPackage(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	var pkg Package
	items := &PackageItems{
		Items: []*PackageItem{
			{
				Description: proto.String("abc"),
				Count:       proto.Int32(123),
			},
		},
	}
	err = pkg.Create(tx, items, Coord{3, 4}, sql.NullInt64{}, 1)
	if err != nil {
		t.Error(err)
	}
	items1 := new(PackageItems)
	err = tx.QueryRow(`SELECT items FROM package WHERE id = $1`,
		pkg).Scan(items1)
	if !proto.Equal(items, items1) {
		t.Errorf("%v != %v", items, items1)
	}

	err = pkg.SetDelivered(tx)
	if err != nil {
		t.Error(err)
	}
}
