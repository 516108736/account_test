package account_test

import (
	"fmt"
	"github.com/516108736/account_test/mpts"

	"github.com/516108736/account_test/fastdb"
	"github.com/516108736/account_test/iavl"
	"github.com/516108736/account_test/mpt"
	"github.com/syndtr/goleveldb/leveldb"
)

type Store interface {
	Update(addr []byte, value []byte)
	Type() string
	Commit()
	Get(addr []byte) []byte
	Delete(addr []byte)
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
	case "mpts":
		return mpts.New()

	default:
		panic(fmt.Errorf("%v not support yet", key))
	}

}
