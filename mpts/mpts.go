package mpts

import (
	"fmt"
	currCommon "github.com/516108736/account_test/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

type MPTS struct {
	st   *trie.SecureTrie
	db   *ethdb.LDBDatabase
	trDB *trie.Database
}

func New() *MPTS {
	dbPath := "./mpts_data"
	os.RemoveAll(dbPath)
	db, err := ethdb.NewLDBDatabase(dbPath, currCommon.DBCache, currCommon.DBHandle)
	currCommon.Checkerr(err)

	trDB := trie.NewDatabaseWithCache(db, 0)
	st, err := trie.NewSecure(common.Hash{}, trDB, 0)
	currCommon.Checkerr(err)

	return &MPTS{
		st:   st,
		db:   db,
		trDB: trDB,
	}

}
func (m *MPTS) Update(addr []byte, value []byte) {
	err := m.st.TryUpdate(addr, value)
	currCommon.Checkerr(err)
}

func (m *MPTS) Get(addr []byte) []byte {
	value, _ := m.st.TryGet(addr)
	return value
}

func (m *MPTS) Delete(addr []byte) {
	m.st.TryDelete(addr)

}

func (m *MPTS) Commit() {
	root, err := m.st.Commit(nil)
	currCommon.Checkerr(err)
	err = m.trDB.Commit(root, false)
	currCommon.Checkerr(err)
	m.st, err = trie.NewSecure(root, m.trDB, 0)
	currCommon.Checkerr(err)
	fmt.Println("MPTS new root", root.String())
}

func (m *MPTS) Type() string {
	return "MPTS"
}

func (m *MPTS) DB() *leveldb.DB {
	return m.db.LDB()
}

func (m *MPTS) RangeFromRoot() (int, int) {
	return m.st.RangeFromRoot()
}
