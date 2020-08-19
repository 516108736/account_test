package mpt

import (
	"fmt"
	"github.com/ethereum/go-ethereum/trie"
	"os"

	"github.com/syndtr/goleveldb/leveldb"

	currCommon "github.com/516108736/account_test/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

type MPT struct {
	tr   *trie.Trie
	db   *ethdb.LDBDatabase
	trDB *trie.Database
}

func New() *MPT {
	dbPath := "./mpt_data"
	os.RemoveAll(dbPath)
	db, err := ethdb.NewLDBDatabase(dbPath, currCommon.DBCache, currCommon.DBHandle)
	currCommon.Checkerr(err)

	trDB := trie.NewDatabaseWithCache(db, 0)
	tr, err := trie.New(common.Hash{}, trDB)
	currCommon.Checkerr(err)

	return &MPT{
		tr:   tr,
		db:   db,
		trDB: trDB,
	}

}
func (m *MPT) Update(addr []byte, value []byte) {
	err := m.tr.TryUpdate(addr, value)
	currCommon.Checkerr(err)
}

func (m *MPT) Get(addr []byte) []byte {
	value, _ := m.tr.TryGet(addr)
	return value
}

func (m *MPT) Delete(addr []byte) {
	m.tr.TryDelete(addr)

}

func (m *MPT) Commit() {
	root, err := m.tr.Commit(nil)
	currCommon.Checkerr(err)
	err = m.trDB.Commit(root, false)
	currCommon.Checkerr(err)
	m.tr, err = trie.New(root, m.trDB)
	currCommon.Checkerr(err)
	fmt.Println("MPT new root", root.String())
}

func (m *MPT) Type() string {
	return "MPT"
}

func (m *MPT) DB() *leveldb.DB {
	return m.db.LDB()
}

func (m *MPT) RangeFromRoot() (int, int) {
	return m.tr.RangeFromRoot()
}
