package server

import (
	"database/sql"

	"github.com/golang/protobuf/proto"
	"gitlab.oit.duke.edu/rz78/ups/db"
	"gitlab.oit.duke.edu/rz78/ups/pb/bridge"
)

// PackageIdReqs produces tracking numbers for multiple packages.
func (s *Server) PackageIdReqs(packages []*bridge.Package) (resp []*bridge.ResponsePackageId, err error) {
	resp = make([]*bridge.ResponsePackageId, len(packages))
	for i, pkg := range packages {
		resp[i], err = s.PackageIdReq(pkg)
		if err != nil {
			s := err.Error()
			resp[i].Error = &s
		}
	}
	return
}

// PackageIdReq produces a tracking number for one package.
func (s *Server) PackageIdReq(pkg *bridge.Package) (resp *bridge.ResponsePackageId, err error) {
	resp = new(bridge.ResponsePackageId)
	// TODO(rz78): package detail is ignored
	err = db.WithTx(s.db, func(tx *sql.Tx) (err error) {
		var pkgId db.Package
		err = pkgId.Create(tx, "N/A", db.Coord{pkg.GetX(), pkg.GetY()}, pkg.GetWarehouseId())
		if err != nil {
			return
		}
		resp.PackageId = proto.Int64(int64(pkgId))
		return
	})
	return
}

func (s *Server) onPackageDelivered(pkg db.Package) error {
	return db.WithTx(s.db, func(tx *sql.Tx) error {
		return pkg.SetDelivered(tx)
	})
}
