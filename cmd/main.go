package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	currCommon "github.com/516108736/account_test/common"
	"github.com/516108736/account_test/mpt"
	"github.com/516108736/account_test/mpts"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rlp"
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

func EncodeBalance(balance *big.Int) []byte {
	value := state.Account{
		Nonce:    0,
		Balance:  balance,
		Root:     common.Hash{},
		CodeHash: nil,
	}
	valueBytes, err := rlp.EncodeToBytes(value)
	currCommon.Checkerr(err)
	return valueBytes
}

func DecodeBalance(bz []byte) *big.Int {
	if len(bz) == 0 {
		return common.Big0
	}
	a := new(state.Account)
	err := rlp.DecodeBytes(bz, a)
	currCommon.Checkerr(err)
	return a.Balance
}
func AddAccounts(store account_test.Store, from int, to int) {
	ts := time.Now()
	for index := from; index < to; index++ {
		store.Update(IntTo20Bytes(index), EncodeBalance(coinsForTrieInit))

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
		store.Get(IntTo20Bytes(addrInt))
	}

	if DecodeBalance(store.Get(IntTo20Bytes(addrInt))).Cmp(shouldCoins) != 0 {
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
		store.Update(IntTo20Bytes(addrInt), EncodeBalance(coinsForUpdate))
	}
	store.Commit()
	if DecodeBalance(store.Get(IntTo20Bytes(addrInt))).Cmp(coinsForUpdate) != 0 {
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
		store.Delete(IntTo20Bytes(addrInt))

	}
	store.Commit()

	if DecodeBalance(store.Get(IntTo20Bytes(addrInt))).Cmp(common.Big0) != 0 {
		panic("RemoveAccounts failed")
	}
	fmt.Println(store.Type(), "Delete End", time.Now().Sub(ts).Seconds())
}

func PrintDB(db *leveldb.DB) {
	mapp := make(map[int]map[int]int)
	cnt := 0
	cntBytes := 0
	itr := db.NewIterator(nil, nil)
	for itr.Next() {
		key := itr.Key()
		value := itr.Value()
		cnt++
		cntBytes += len(key)
		cntBytes += len(value)

		if data, ok := mapp[len(key)]; ok {
			if _, ok1 := data[len(value)]; ok1 {
				mapp[len(key)][len(value)]++
			} else {
				mapp[len(key)][len(value)] = 1
			}
		} else {
			mapp[len(key)] = make(map[int]int)
			mapp[len(key)][len(value)] = 1
		}
	}
	fmt.Println("DB Print", "len(key)", cnt, "sum{len(key)+len(value)}", cntBytes)
	fmt.Println("mapp", mapp)
	sum := 0
	for k, vs := range mapp {
		for value, c := range vs {
			t := (k + value) * c
			sum += t
			fmt.Println("-----------", k, value, c)
		}
	}
	fmt.Println("qingsuan", sum)
}

func RangeMPT(store account_test.Store) {
	switch n := store.(type) {
	case *mpt.MPT:
		a, b := n.RangeFromRoot()
		fmt.Println("MPT RangeFromRoot", "keyCnt", a, "bytesAll", b)
	case *mpts.MPTS:
		a, b := n.RangeFromRoot()
		fmt.Println("MPTS RangeFromRoot", "keyCnt", a, "bytesAll", b)
	}
}

var (
	updateNumber     = 10000
	coinsForTrieInit = new(big.Int).SetUint64(1000)
	coinsForUpdate   = new(big.Int).SetUint64(6666)

	initNumber = flag.Int("initNumber", 1000*100*1, "init number")
	typ        = flag.String("typ", "fastdb", "send or query")
)

func CloseDB(db *leveldb.DB) {
	if err := db.Close(); err != nil {
		panic(err)
	}
}
func main() {
	flag.Parse()
	initNumber := *initNumber
	fmt.Println("initNumber", initNumber, "updateNumber", updateNumber, "type", *typ)
	store := account_test.NewStore(*typ)

	AddAccounts(store, 0, initNumber)
	PrintDB(store.DB())
	RangeMPT(store)

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

	CloseDB(store.DB())
}
