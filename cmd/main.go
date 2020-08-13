package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/516108736/account_test"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/syndtr/goleveldb/leveldb"
)

func IntTo20Bytes(data int) []byte {
	ret := make([]byte, 8)
	binary.BigEndian.PutUint64(ret, uint64(data))
	h := sha256.New()
	h.Write(ret)
	hash := h.Sum(nil)
	return hash[:20]
}

func AddAccounts(store account_test.Store, from int, to int) {
	ts := time.Now()
	for index := from; index < to; index++ {
		store.SetBalance(IntTo20Bytes(index), coinsForTrieInit)

		if index != 0 && index%100000 == 0 {
			fmt.Println("SetAccounts handle index", index, time.Now().Sub(ts).Seconds())
		}
		if index != 0 && index%1000000 == 0 {
			store.Commit()
			fmt.Println("SetAccounts commit index", index, time.Now().Sub(ts).Seconds())
		}

	}
	store.Commit()
	fmt.Println(store.Type(), "ADD END from", from, "to", to, "ts", time.Now().Sub(ts).Seconds())
}

func GetAccounts(store account_test.Store, from int, to int, length int, random bool, shouldCoins *big.Int) []exported.Account {
	ts := time.Now()
	addrInt := 0
	for index := 0; index < length; index++ {
		addrInt = index + from
		if random {
			addrInt = rand.Intn(to-from) + from
		}
		store.GetBalance(IntTo20Bytes(addrInt))
	}

	if store.GetBalance(IntTo20Bytes(addrInt)).Cmp(shouldCoins) != 0 {
		panic("GetAccounts failed")
	}
	fmt.Println(store.Type(), "GET End from", from, "to", to, "length", length, time.Now().Sub(ts).Seconds())
	return nil
}

func UpdateAccounts(store account_test.Store, from int, to int, length int, random bool) {
	ts := time.Now()
	addrInt := 0
	for index := 0; index < length; index++ {
		addrInt = index + from
		if random {
			addrInt = rand.Intn(to-from) + from
		}
		store.SetBalance(IntTo20Bytes(addrInt), coinsForUpdate)
	}
	store.Commit()
	if store.GetBalance(IntTo20Bytes(addrInt)).Cmp(coinsForUpdate) != 0 {
		panic("UpdateAccounts failed")
	}
	fmt.Println(store.Type(), "UPDATE END from", from, "to", to, "length", length, time.Now().Sub(ts).Seconds())
}

func RemoveAccounts(store account_test.Store, from int, to int, length int, random bool) {
	ts := time.Now()
	addrInt := 0

	for index := 0; index < length; index++ {
		addrInt = index + from
		if random {
			addrInt = rand.Intn(to-from) + from
		}
		store.DeleteAddr(IntTo20Bytes(addrInt))

	}
	store.Commit()

	if store.GetBalance(IntTo20Bytes(addrInt)).Cmp(common.Big0) != 0 {
		panic("RemoveAccounts failed")
	}
	fmt.Println(store.Type(), "IAVL Delete End", time.Now().Sub(ts).Seconds())
}

func PrintDB(db *leveldb.DB) {
	cnt := 0
	cntBytes := 0
	itr := db.NewIterator(nil, nil)
	for itr.Next() {
		key := itr.Key()
		value := itr.Value()
		cnt++
		cntBytes += len(key)
		cntBytes += len(value)
	}
	fmt.Println("DB Print", "len(key)", cnt, "sum{len(key)+len(value)}", cntBytes)
}

var (
	updateNumber     = 10000
	coinsForTrieInit = new(big.Int).SetUint64(1000)
	coinsForUpdate   = new(big.Int).SetUint64(6666)

	initNumber = flag.Int("initNumber", 1000*100*1, "init number")
	typ        = flag.String("typ", "fastdb", "send or query")
)

func main() {
	flag.Parse()
	initNumber := *initNumber
	fmt.Println("initNumber", initNumber, "updateNumber", updateNumber, "type", *typ)
	store := account_test.NewStore(*typ)

	AddAccounts(store, 0, initNumber)
	PrintDB(store.DB())

	AddAccounts(store, initNumber, initNumber+updateNumber)

	GetAccounts(store, 0, updateNumber, updateNumber, false, coinsForTrieInit)
	GetAccounts(store, initNumber, initNumber+updateNumber, updateNumber, false, coinsForTrieInit)
	GetAccounts(store, 0, initNumber+updateNumber, updateNumber, true, coinsForTrieInit)

	UpdateAccounts(store, 0, 0+updateNumber, updateNumber, false)
	UpdateAccounts(store, initNumber, initNumber+updateNumber, updateNumber, false)
	UpdateAccounts(store, updateNumber, initNumber, updateNumber, true)

	RemoveAccounts(store, 0, 0+updateNumber, updateNumber, false)
	RemoveAccounts(store, initNumber, initNumber+updateNumber, updateNumber, false)
	RemoveAccounts(store, updateNumber, initNumber, updateNumber, true)
}
