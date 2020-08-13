package mpt

import (
	"math/big"
	"os"

	"github.com/syndtr/goleveldb/leveldb"

	currCommon "github.com/516108736/account_test/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
)

type MPT struct {
	state *state.StateDB
	db    *ethdb.LDBDatabase
}

func New() *MPT {
	dbPath := "./mpt_data"
	os.RemoveAll(dbPath)
	db, err := ethdb.NewLDBDatabase(dbPath, 128, 128)
	currCommon.Checkerr(err)

	stateDB, err := state.New(common.Hash{}, state.NewDatabaseWithCache(db, 128))
	currCommon.Checkerr(err)
	return &MPT{
		state: stateDB,
		db:    db,
	}
}
func (m *MPT) SetBalance(addr []byte, coin *big.Int) {
	m.state.SetBalance(common.BytesToAddress(addr), coin)
}

func (m *MPT) GetBalance(addr []byte) *big.Int {
	return m.state.GetBalance(common.BytesToAddress(addr))
}

func (m *MPT) DeleteAddr(addr []byte) {
	m.state.Suicide(common.BytesToAddress(addr))
}

func (m *MPT) Commit() {
	root, err := m.state.Commit(true)
	currCommon.Checkerr(err)
	m.state.Database().TrieDB().Commit(root, false)
	m.state, err = state.New(root, state.NewDatabaseWithCache(m.db, 128))
	currCommon.Checkerr(err)
}

func (m *MPT) Type() string {
	return "MPT"
}

func (m *MPT) DB() *leveldb.DB {
	return m.db.LDB()
}
