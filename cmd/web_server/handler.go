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
	case req.GetPackageList != nil:
		resp, err = handleGetPackageList(database, req.GetPackageList)
	case req.GetPackageDetail != nil:
		resp, err = handleGetPackageDetail(database, *req.GetPackageDetail)
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

func handleGetPackageList(database *sql.DB, req *web.PkgListReq) (resp *web.Response, err error) {
	resp = &web.Response{
		PackageList: new(web.PkgList),
	}
	const (
		pre = `SELECT id, create_time, load_time, deliver_time, truck_status FROM package_view `
		post = ` ORDER BY create_time DESC LIMIT $1 OFFSET $2`
	)
	var limit sql.NullInt64
	if req.Limit != nil {
		limit.Valid = true
		limit.Int64 = *req.Limit
	}
	offset := req.GetOffset()

	var rows *sql.Rows

	if len(req.PackageIds) > 0 {
		resp.PackageList.Total = proto.Int64(int64(len(req.PackageIds)))
		var stmt *sql.Stmt
		stmt, err = database.Prepare(pre + `WHERE id = $3` + post)
		if err != nil {
			return
		}
		defer stmt.Close()
		for _, pkgId := range req.PackageIds {
			err = appendPkgInfo(resp, stmt.QueryRow(nil, nil, pkgId))
			if err == sql.ErrNoRows {
				err = fmt.Errorf("invalid package id: %v", pkgId)
			}
			if err != nil {
				return
			}
		}
	} else {
		var queryTotal *sql.Row
		var (
			query string
			args = []interface{}{limit, offset, nil}[:2]
		)
		if req.UserId != nil {
			queryTotal = database.QueryRow(`SELECT COUNT(*) FROM package WHERE user_id = $1`, *req.UserId)
			query = pre + `WHERE user_id = $3` + post
			args = append(args, *req.UserId)
		} else {
			queryTotal = database.QueryRow(`SELECT COUNT(*) FROM package`)
			query = pre + post
		}
		var total int64
		err = queryTotal.Scan(&total)
		if err != nil {
			return
		}
		resp.PackageList.Total = &total
		rows, err = database.Query(query, args...)
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			err = appendPkgInfo(resp, rows)
			if err != nil {
				return
			}
		}
	}
	return
}

func handleGetPackageDetail(database *sql.DB, pkgId int64) (resp *web.Response, err error) {
	resp = new(web.Response)
	const query = `SELECT user_id, items, destination, create_time, load_time, deliver_time, truck_status `+
		`FROM package_view WHERE id = $1`
	var (
		userId sql.NullInt64
		items  db.PackageItems
		dest   db.Coord
		ctime  int64
		ltime  sql.NullInt64
		dtime  sql.NullInt64
		status *db.TruckStatus
	)
	err = database.QueryRow(query, pkgId).Scan(&userId, &items, &dest, &ctime, &ltime, &dtime, &status)
	resp.PackageDetail = &web.PkgDetail{
		Items: convertItems(&items),
		X: &dest.X,
		Y: &dest.Y,
		UserId: &userId.Int64,
		Status: convertStatus(ctime, ltime, dtime, status),
	}

	return
}

func convertItems(items *db.PackageItems) (r []*web.PkgDetail_Item) {
	for _, item := range items.Items {
		r = append(r, &web.PkgDetail_Item{
			Description: item.Description,
			Amount:      item.Count,
		})
	}
	return
}

func convertStatus(ctime int64, ltime, dtime sql.NullInt64, status *db.TruckStatus) (r []*web.PkgDetail_Status) {
	r = append(r, &web.PkgDetail_Status{
		Status:    proto.String("created"),
		Timestamp: &ctime,
	})
	if ltime.Valid {
		r = append(r, &web.PkgDetail_Status{
			Status:    proto.String("loaded to truck"),
			Timestamp: &ltime.Int64,
		})
	}
	if dtime.Valid {
		r = append(r, &web.PkgDetail_Status{
			Status:    proto.String("delivered"),
			Timestamp: &dtime.Int64,
		})
	} else if status != nil {
		s := convertTruckStatus(*status)
		r = append(r, &web.PkgDetail_Status{Status: &s})
	}
	return
}

func convertTruckStatus(status db.TruckStatus) string {
	switch status {
	case db.TO_WAREHOUSE:
		return "truck en route"
	case db.AT_WAREHOUSE:
		return "truck waiting for package"
	case db.DELIVERING:
		return "delivering"
	}
	panic("BUG")
}

type scanner interface {
	Scan(...interface{}) error
}

func appendPkgInfo(resp *web.Response, sc scanner) (err error) {
	var (
		pkgId  int64
		ctime  int64
		ltime  sql.NullInt64
		dtime  sql.NullInt64
		status *db.TruckStatus
	)
	err = sc.Scan(&pkgId, &ctime, &ltime, &dtime, &status)
	if err != nil {
		return
	}
	var s string
	switch {
	case dtime.Valid:
		s = "delivered"
	case status != nil:
		s = convertTruckStatus(*status)
	default:
		s = "created"
	}
	pkg := &web.PkgList_Info{
		PackageId: &pkgId,
		Status: &s,
		CreateTime: &ctime,
	}
	resp.PackageList.Packages = append(resp.PackageList.Packages, pkg)
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
