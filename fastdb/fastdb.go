package fastdb

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb"

	currCommon "github.com/516108736/account_test/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

type FastDB struct {
	db *ethdb.LDBDatabase
}

func New() *FastDB {
	dbPath := "./fastdb_data"
	os.RemoveAll(dbPath)
	db, err := ethdb.NewLDBDatabase(dbPath, 128, 128)
	currCommon.Checkerr(err)
	return &FastDB{
		db: db,
	}
}

func (f *FastDB) Update(addr []byte, valueBytes []byte) {
	f.db.Put(addr, valueBytes)
}

func (f *FastDB) Type() string {
	return "FastDB"
}

func (f *FastDB) Commit() {
}

func (f *FastDB) Get(addr []byte) []byte {
	bz, _ := f.db.Get(addr)
	return bz
}

func (f *FastDB) Delete(addr []byte) {
	f.db.Delete(addr)
}

func (f *FastDB) DB() *leveldb.DB {
	return f.db.LDB()
}
