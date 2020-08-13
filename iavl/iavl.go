package iavl

import (
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"

	currCommon "github.com/516108736/account_test/common"
	"github.com/cosmos/cosmos-sdk/store/iavl"
	sTypes "github.com/cosmos/cosmos-sdk/store/types"
	dbm "github.com/tendermint/tm-db"
)

func GetDB() dbm.DB {
	var db dbm.DB

	dirPath := "./iavl_data"
	err := os.RemoveAll(dirPath)
	if err != nil {
		panic(err)
	}
	db, err = dbm.NewGoLevelDBWithOpts("trie_test", dirPath, &opt.Options{
		OpenFilesCacheCapacity: currCommon.DBHandle,
		BlockCacheCapacity:     currCommon.DBCache / 2 * opt.MiB,
		WriteBuffer:            currCommon.DBCache / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if err != nil {
		panic(err)
	}
	return db
}

type IAVL struct {
	iavl *iavl.Store
	db   *dbm.GoLevelDB
}

func New() *IAVL {
	db := GetDB()
	ivalStore, err := iavl.LoadStore(db, sTypes.CommitID{}, sTypes.PruningOptions{}, false)
	currCommon.Checkerr(err)
	return &IAVL{
		iavl: ivalStore.(*iavl.Store),
		db:   db.(*dbm.GoLevelDB),
	}
}

func (i *IAVL) Update(addr []byte, value []byte) {
	i.iavl.Set(addr, value)
}

func (i *IAVL) Get(addr []byte) []byte {
	return i.iavl.Get(addr)
}

func (i *IAVL) Delete(addr []byte) {
	i.iavl.Delete(addr)
}

func (i *IAVL) Commit() {
	id := i.iavl.Commit()
	t, err := iavl.LoadStore(i.db, id, sTypes.PruningOptions{}, false)
	currCommon.Checkerr(err)
	i.iavl = t.(*iavl.Store)
}

func (i *IAVL) Type() string {
	return "IAVL"
}

func (i *IAVL) DB() *leveldb.DB {
	return i.db.DB()
}
