package lock

import (
	"strconv"
	"sync/atomic"

	"6.5840/kvsrv1/rpc"
	kvtest "6.5840/kvtest1"
)

const valUnlocked = ""

var nextLockID int32 = 0 // Atomic counter for lock IDs

type Lock struct {
	// IKVClerk is a go interface for k/v clerks: the interface hides
	// the specific Clerk type of ck but promises that ck supports
	// Put and Get.  The tester passes the clerk in when calling
	// MakeLock().
	ck kvtest.IKVClerk
	// You may add code here
	name string // The name of the lock, used as the key in the k/v store.
	id   int32
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// Use l as the key to store the "lock state" (you would have to decide
// precisely what the lock state is).
func MakeLock(ck kvtest.IKVClerk, l string) *Lock {
	lk := &Lock{ck: ck}
	// You may add code here
	lk.name = l
	lk.id = atomic.AddInt32(&nextLockID, 1) // Assign a unique ID to the lock
	return lk
}

func (lk *Lock) Acquire() {
	// Your code here
	acquired := false
	ck := lk.ck

	for !acquired {
		currId, currVer, err := ck.Get(lk.name)
		if currId != valUnlocked {
			// If the lock is already held, we need to wait and retry.
			continue
		}
		if err == rpc.ErrNoKey {
			currVer = 0
		}

		reply := ck.Put(lk.name, strconv.Itoa(int(lk.id)), currVer) // Attempt to acquire the lock by setting a value.
		switch reply {
		case rpc.OK:
			acquired = true // Lock acquired successfully
		case rpc.ErrMaybe:
			currId, _, _ := ck.Get(lk.name)
			if currId == strconv.Itoa(int(lk.id)) {
				acquired = true // Lock was acquired by this client
			}
		}
	}

}

func (lk *Lock) Release() {
	// Your code here
	ck := lk.ck
	currId, currVer, _ := ck.Get(lk.name)
	if currId == valUnlocked || currId != strconv.Itoa(int(lk.id)) {
		return
	}

	reply := ck.Put(lk.name, valUnlocked, currVer) // Attempt to acquire the lock by setting a value.
	switch reply {
	case rpc.OK:

	case rpc.ErrMaybe:

	case rpc.ErrVersion:
	}
	// if reply == rpc.ErrMaybe {
	// 	triedVer := currVer
	// 	state := valLocked
	// 	for state == valLocked && currVer == triedVer {
	// 		state, currVer, _ = ck.Get(lk.name)
	// 		if state == valLocked && currVer == triedVer {
	// 			reply = ck.Put(lk.name, valUnlocked, currVer) // Ensure we release the lock correctly.
	// 			if reply == rpc.OK {
	// 				return // Successfully released the lock.
	// 			}
	// 		}
	// 	}
	// }

}
