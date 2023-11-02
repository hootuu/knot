package knotx

import (
	"github.com/hootuu/knot/dbx"
	"github.com/hootuu/tome/vn"
	"github.com/hootuu/utils/configure"
	"github.com/hootuu/utils/errors"
	"github.com/hootuu/utils/sys"
	"path/filepath"
	"sync"
)

var gChainDbxDict = make(map[string]*dbx.DBX)
var gChainDbxLock sync.Mutex

func doGetChainDBX(vnID vn.ID) (*dbx.DBX, *errors.Error) {
	dbxM, ok := gChainDbxDict[vnID.S()]
	if ok {
		return dbxM, nil
	}
	knotPathStr := configure.GetString("knot.db.path", filepath.Join(sys.WorkingDirectory, ".knot"))
	gChainDbxLock.Lock()
	defer gChainDbxLock.Unlock()
	dbxM, err := dbx.NewDBX(vnID.S(), knotPathStr)
	if err != nil {
		return nil, err
	}
	gChainDbxDict[vnID.S()] = dbxM
	return dbxM, nil
}

func Close() {
	for k, dbxM := range gChainDbxDict {
		err := dbxM.Close()
		if err != nil {
			sys.Warn("knot.", k, " close failed: ", err.Err)
			continue
		}
	}
}
