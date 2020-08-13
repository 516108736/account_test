package fastdb

import (
	"math/big"
	"os"

	"github.com/syndtr/goleveldb/leveldb"

	currCommon "github.com/516108736/account_test/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
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

func (f *FastDB) SetBalance(addr []byte, coin *big.Int) {
	value := state.Account{
		Nonce:    0,
		Balance:  coin,
		Root:     common.Hash{},
		CodeHash: nil,
	}
	valueBytes, err := rlp.EncodeToBytes(value)
	currCommon.Checkerr(err)
	f.db.Put(addr, valueBytes)
}

func (f *FastDB) Type() string {
	return "FastDB"
}

func (f *FastDB) Commit() {

}

func (f *FastDB) GetBalance(addr []byte) *big.Int {
	bz, err := f.db.Get(addr)
	if err != nil {
		return new(big.Int)
	}
	a := new(state.Account)
	err = rlp.DecodeBytes(bz, a)
	currCommon.Checkerr(err)
	return a.Balance
}

func (f *FastDB) DeleteAddr(addr []byte) {
	currCommon.Checkerr(f.db.Delete(addr))
}

func (f *FastDB) DB() *leveldb.DB {
	return f.db.LDB()
}
