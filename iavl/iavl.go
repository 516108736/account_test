package iavl

import (
	"math/big"
	"os"

	"github.com/syndtr/goleveldb/leveldb"

	currCommon "github.com/516108736/account_test/common"
	"github.com/cosmos/cosmos-sdk/store/iavl"
	sTypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
	dbm "github.com/tendermint/tm-db"
)

func GetDB() dbm.DB {
	var db dbm.DB

	dirPath := "./iavl_data"
	err := os.RemoveAll(dirPath)
	if err != nil {
		panic(err)
	}
	db, err = dbm.NewGoLevelDB("trie_test", dirPath)
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

func (i *IAVL) SetBalance(addr []byte, coin *big.Int) {
	value := state.Account{
		Nonce:    0,
		Balance:  coin,
		Root:     common.Hash{},
		CodeHash: nil,
	}
	valueBytes, err := rlp.EncodeToBytes(value)
	currCommon.Checkerr(err)
	i.iavl.Set(addr, valueBytes)
}

func (i *IAVL) GetBalance(addr []byte) *big.Int {
	bz := i.iavl.Get(addr)
	if len(bz) == 0 {
		return common.Big0
	}
	a := new(state.Account)
	err := rlp.DecodeBytes(bz, a)
	currCommon.Checkerr(err)
	return a.Balance
}

func (i *IAVL) DeleteAddr(addr []byte) {
	i.iavl.Delete(addr)
}

func (i *IAVL) Commit() {
	i.iavl.Commit()
}

func (i *IAVL) Type() string {
	return "IAVL"
}

func (i *IAVL) DB() *leveldb.DB {
	return i.db.DB()
}
