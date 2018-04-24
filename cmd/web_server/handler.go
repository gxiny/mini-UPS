package main

import (
	"database/sql"
	"errors"
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
		resp, err = handleGetPackageStatus(database, *req.GetPackageStatus)
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

const pkgCommonSQL = `
SELECT package.id, package.items, package.destination,
	package.create_time, package.deliver_time,
	truck.status
FROM package LEFT JOIN truck
ON package.truck_id = truck.id WHERE `

func handleGetPackages(database *sql.DB, userId int64) (resp *web.Response, err error) {
	return pkgQuery(database, pkgCommonSQL+`package.user_id = $1`, userId)
}

var errInvalidPkgId = errors.New("invalid package id")

func handleGetPackageStatus(database *sql.DB, pkgId int64) (resp *web.Response, err error) {
	resp, err = pkgQuery(database, pkgCommonSQL+`package.id = $1`, pkgId)
	if err == nil && len(resp.GetPackages()) == 0 {
		err = errInvalidPkgId
	}
	return
}

func pkgQuery(database *sql.DB, query string, args ...interface{}) (resp *web.Response, err error) {
	resp = new(web.Response)
	err = db.WithTx(database, func(tx *sql.Tx) (err error) {
		tx.Exec(`SET TRANSACTION READ ONLY`)
		rows, err := tx.Query(query, args...)
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			var (
				pkgId  int64
				items  db.PackageItems
				dest   db.Coord // TODO: include dest in the response
				ctime  int64
				dtime  sql.NullInt64
				status *db.TruckStatus
			)
			err = rows.Scan(&pkgId, &items, &dest, &ctime, &dtime, &status)
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
				Status: convertStatus(ctime, dtime, status),
			}
			resp.Packages = append(resp.Packages, pkg)
		}
		return
	})
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

func convertStatus(ctime int64, dtime sql.NullInt64, status *db.TruckStatus) (r []*web.Status) {
	r = append(r, &web.Status{
		Status:    proto.String("created"),
		Timestamp: &ctime,
	})
	if !dtime.Valid {
		if status != nil {
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
	} else {
		r = append(r, &web.Status{
			Status:    proto.String("delivered"),
			Timestamp: &dtime.Int64,
		})
	}
	return
}

var (
	errPkgDelivered  = errors.New("package is delivered")
	errPkgDelivering = errors.New("package is out for delivery")
	errPermDenied    = errors.New("permission denied")
)

func handleChangeDestination(database *sql.DB, dest *web.PkgDest) (resp *web.Response, err error) {
	resp = new(web.Response)
	err = db.WithTx(database, func(tx *sql.Tx) (err error) {
		var (
			dtime  sql.NullInt64
			userId sql.NullInt64
			truck  *db.Truck
		)
		err = tx.QueryRow(`SELECT deliver_time, user_id, truck_id FROM package WHERE id = $1 FOR UPDATE`,
			dest.PackageId).Scan(&dtime, &userId, &truck)
		if err != nil {
			return
		}
		if !userId.Valid || dest.GetUserId() != userId.Int64 {
			err = errPermDenied
		} else if dtime.Valid { // delivered
			err = errPkgDelivered
		} else if truck != nil {
			var status db.TruckStatus
			err = tx.QueryRow(`SELECT status FROM truck WHERE id = $1 FOR SHARE`,
				*truck).Scan(&status)
			if err == nil && status == db.DELIVERING { // delivering
				err = errPkgDelivering
			}
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
