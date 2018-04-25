package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb/web"
)

func handleRequest(database *sql.DB, req *web.Request) (resp *web.Response) {
	var err error
	switch {
	case req.NewUser != nil:
		resp, err = handleNewUser(database, *req.NewUser)
	case req.GetPackages != nil:
		resp, err = handleGetPackages(database, *req.GetPackages)
	case req.GetPackageStatus != nil:
		resp, err = handleGetPackageStatus(database, req.GetPackageStatus)
	case req.ChangeDestination != nil:
		resp, err = handleChangeDestination(database, req.ChangeDestination)
	}
	if err != nil {
		s := err.Error()
		log.Println(s)
		resp.Error = &s
	}
	return
}

func handleNewUser(database *sql.DB, username string) (resp *web.Response, err error) {
	resp = new(web.Response)
	err = db.WithTx(database, func(tx *sql.Tx) (err error) {
		var user db.User
		err = user.Create(tx)
		if err != nil {
			return
		}
		resp.UserId = proto.Int64(int64(user))
		return
	})
	return
}

const pkgCommonSQL = `SELECT id, items, destination, create_time, load_time, deliver_time, truck_status FROM package_view `

func handleGetPackages(database *sql.DB, userId int64) (resp *web.Response, err error) {
	resp = new(web.Response)
	rows, err := database.Query(pkgCommonSQL+`WHERE user_id = $1 ORDER BY create_time DESC`, userId)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = appendPackage(resp, rows)
		if err != nil {
			return
		}
	}
	return
}

func handleGetPackageStatus(database *sql.DB, pkgIds []int64) (resp *web.Response, err error) {
	resp = new(web.Response)
	var invalidIds []int64
	stmt, err := database.Prepare(pkgCommonSQL + `WHERE id = $1 ORDER BY create_time DESC`)
	if err != nil {
		return
	}
	defer stmt.Close()
	for _, pkgId := range pkgIds {
		err = appendPackage(resp, stmt.QueryRow(pkgId))
		if err == sql.ErrNoRows {
			invalidIds = append(invalidIds, pkgId)
			err = nil
		}
		if err != nil {
			log.Printf("%#v %#v", err, sql.ErrNoRows)
			return
		}
	}
	if err == nil && len(invalidIds) > 0 {
		err = fmt.Errorf("invalid package id: %v", invalidIds)
	}
	return
}

type scanner interface {
	Scan(...interface{}) error
}

func appendPackage(resp *web.Response, sc scanner) (err error) {
	var (
		pkgId  int64
		items  db.PackageItems
		dest   db.Coord
		ctime  int64
		ltime  sql.NullInt64
		dtime  sql.NullInt64
		status *db.TruckStatus
	)
	err = sc.Scan(&pkgId, &items, &dest, &ctime, &ltime, &dtime, &status)
	if err != nil {
		return
	}
	pkg := &web.Package{
		PackageId: &pkgId,
		Detail: &web.PkgDetail{
			Items: convertItems(&items),
			X:     &dest.X,
			Y:     &dest.Y,
		},
		Status: convertStatus(ctime, ltime, dtime, status),
	}
	resp.Packages = append(resp.Packages, pkg)
	return
}

func convertItems(items *db.PackageItems) (r []*web.Item) {
	for _, item := range items.Items {
		r = append(r, &web.Item{
			Description: item.Description,
			Amount:      item.Count,
		})
	}
	return
}

func convertStatus(ctime int64, ltime, dtime sql.NullInt64, status *db.TruckStatus) (r []*web.Status) {
	r = append(r, &web.Status{
		Status:    proto.String("created"),
		Timestamp: &ctime,
	})
	if ltime.Valid {
		r = append(r, &web.Status{
			Status:    proto.String("loaded to truck"),
			Timestamp: &ltime.Int64,
		})
	}
	if dtime.Valid {
		r = append(r, &web.Status{
			Status:    proto.String("delivered"),
			Timestamp: &dtime.Int64,
		})
	} else if status != nil {
		var s string
		switch *status {
		case db.TO_WAREHOUSE:
			s = "truck en route"
		case db.AT_WAREHOUSE:
			s = "truck waiting for package"
		case db.DELIVERING:
			s = "delivering"
		}
		r = append(r, &web.Status{Status: &s})
	}
	return
}

var (
	errPkgLoaded  = errors.New("package is delivered / being delivered")
	errPermDenied = errors.New("permission denied")
)

func handleChangeDestination(database *sql.DB, dest *web.PkgDest) (resp *web.Response, err error) {
	resp = new(web.Response)
	err = db.WithTx(database, func(tx *sql.Tx) (err error) {
		var (
			ltime  sql.NullInt64
			userId sql.NullInt64
		)
		err = tx.QueryRow(`SELECT load_time, user_id FROM package WHERE id = $1 FOR UPDATE`,
			dest.PackageId).Scan(&ltime, &userId)
		if err != nil {
			return
		}
		if !userId.Valid || dest.GetUserId() != userId.Int64 {
			err = errPermDenied
		} else if ltime.Valid {
			err = errPkgLoaded
		}
		if err != nil {
			return
		}
		_, err = tx.Exec(`UPDATE package SET destination = $1 WHERE id = $2`,
			db.CoordXY(dest), dest.PackageId)
		return
	})
	return
}
