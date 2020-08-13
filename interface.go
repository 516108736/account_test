package account_test

import (
	"fmt"
	"math/big"

	"github.com/516108736/account_test/fastdb"
	"github.com/516108736/account_test/iavl"
	"github.com/516108736/account_test/mpt"
	"github.com/syndtr/goleveldb/leveldb"
)

type Store interface {
	SetBalance(addr []byte, coin *big.Int)
	Type() string
	Commit()
	GetBalance(addr []byte) *big.Int
	DeleteAddr(addr []byte)
	DB() *leveldb.DB
}

func NewStore(key string) Store {
	switch key {
	case "iavl":
		return iavl.New()
	case "mpt":
		return mpt.New()
	case "fastdb":
		return fastdb.New()
	default:
		panic(fmt.Errorf("%v not support yet", key))
	}

}
