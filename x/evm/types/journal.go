package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var ripemd = ethcmn.HexToAddress("0000000000000000000000000000000000000003")

// journalEntry is a modification entry in the state change journal that can be
// reverted on demand.
type journalEntry interface {
	// revert undoes the changes introduced by this journal entry.
	revert(*CommitStateDB)

	// dirtied returns the Ethereum address modified by this journal entry.
	dirtied() *ethcmn.Address
}

// journal contains the list of state modifications applied since the last state
// commit. These are tracked to be able to be reverted in case of an execution
// exception or revertal request.
type journal struct {
	entries []journalEntry         // Current changes tracked by the journal
	dirties map[ethcmn.Address]int // Dirty accounts and the number of changes
}

// newJournal create a new initialized journal.
func newJournal() *journal {
	return &journal{
		dirties: make(map[ethcmn.Address]int),
	}
}

// append inserts a new modification entry to the end of the change journal.
func (j *journal) append(entry journalEntry) {
	j.entries = append(j.entries, entry)
	if addr := entry.dirtied(); addr != nil {
		j.dirties[*addr]++
	}
}

// revert undoes a batch of journalled modifications along with any reverted
// dirty handling too.
func (j *journal) revert(statedb *CommitStateDB, snapshot int) {
	for i := len(j.entries) - 1; i >= snapshot; i-- {
		// Undo the changes made by the operation
		j.entries[i].revert(statedb)

		// Drop any dirty tracking induced by the change
		if addr := j.entries[i].dirtied(); addr != nil {
			if j.dirties[*addr]--; j.dirties[*addr] == 0 {
				delete(j.dirties, *addr)
			}
		}
	}
	j.entries = j.entries[:snapshot]
}

// dirty explicitly sets an address to dirty, even if the change entries would
// otherwise suggest it as clean. This method is an ugly hack to handle the RIPEMD
// precompile consensus exception.
func (j *journal) dirty(addr ethcmn.Address) {
	j.dirties[addr]++
}

// length returns the current number of entries in the journal.
func (j *journal) length() int {
	return len(j.entries)
}

type (
	// Changes to the account trie.
	createObjectChange struct {
		account *ethcmn.Address
	}

	resetObjectChange struct {
		prev *stateObject
	}

	suicideChange struct {
		account     *ethcmn.Address
		prev        bool // whether account had already suicided
		prevBalance sdk.Dec
	}

	// Changes to individual accounts.
	balanceChange struct {
		account *ethcmn.Address
		prev    sdk.Dec
	}

	nonceChange struct {
		account *ethcmn.Address
		prev    uint64
	}

	storageChange struct {
		account        *ethcmn.Address
		key, prevValue ethcmn.Hash
	}

	codeChange struct {
		account            *ethcmn.Address
		prevCode, prevHash []byte
	}

	// Changes to other state values.
	refundChange struct {
		prev uint64
	}

	addLogChange struct {
		txhash ethcmn.Hash
	}

	addPreimageChange struct {
		hash ethcmn.Hash
	}

	touchChange struct {
		account *ethcmn.Address
		// prev      bool
		// prevDirty bool
	}
	accessListAddAccountChange struct {
		address *ethcmn.Address
	}
	accessListAddSlotChange struct {
		address *ethcmn.Address
		slot    *ethcmn.Hash
	}
)

func (ch createObjectChange) revert(s *CommitStateDB) {
	delete(s.stateObjects, *ch.account)
	delete(s.stateObjectsDirty, *ch.account)
}

func (ch createObjectChange) dirtied() *ethcmn.Address {
	return ch.account
}

func (ch resetObjectChange) revert(s *CommitStateDB) {
	s.setStateObject(ch.prev)
}

func (ch resetObjectChange) dirtied() *ethcmn.Address {
	return nil
}

func (ch suicideChange) revert(s *CommitStateDB) {
	so := s.getStateObject(*ch.account)
	if so != nil {
		so.suicided = ch.prev
		so.setBalance(sdk.DefaultBondDenom, ch.prevBalance)
	}
}

func (ch suicideChange) dirtied() *ethcmn.Address {
	return ch.account
}

func (ch touchChange) revert(s *CommitStateDB) {
}

func (ch touchChange) dirtied() *ethcmn.Address {
	return ch.account
}

func (ch balanceChange) revert(s *CommitStateDB) {
	s.getStateObject(*ch.account).setBalance(sdk.DefaultBondDenom, ch.prev)
}

func (ch balanceChange) dirtied() *ethcmn.Address {
	return ch.account
}

func (ch nonceChange) revert(s *CommitStateDB) {
	s.getStateObject(*ch.account).setNonce(ch.prev)
}

func (ch nonceChange) dirtied() *ethcmn.Address {
	return ch.account
}

func (ch codeChange) revert(s *CommitStateDB) {
	s.getStateObject(*ch.account).setCode(ethcmn.BytesToHash(ch.prevHash), ch.prevCode)
}

func (ch codeChange) dirtied() *ethcmn.Address {
	return ch.account
}

func (ch storageChange) revert(s *CommitStateDB) {
	s.getStateObject(*ch.account).setState(ch.key, ch.prevValue)
}

func (ch storageChange) dirtied() *ethcmn.Address {
	return ch.account
}

func (ch refundChange) revert(s *CommitStateDB) {
	s.refund = ch.prev
}

func (ch refundChange) dirtied() *ethcmn.Address {
	return nil
}

func (ch addLogChange) revert(s *CommitStateDB) {
	logs, err := s.GetLogs(ch.txhash)
	if err != nil {
		// panic on unmarshal error
		panic(err)
	}

	// delete logs if entry is empty or has only one item
	if len(logs) <= 1 {
		s.DeleteLogs(ch.txhash)
	} else if err := s.SetLogs(ch.txhash, logs[:len(logs)-1]); err != nil {
		// panic on marshal error
		panic(err)
	}

	s.logSize--
}

func (ch addLogChange) dirtied() *ethcmn.Address {
	return nil
}

func (ch addPreimageChange) revert(s *CommitStateDB) {
	delete(s.preimages, ch.hash)
}

func (ch addPreimageChange) dirtied() *ethcmn.Address {
	return nil
}

func (ch accessListAddAccountChange) revert(s *CommitStateDB) {
	/*
		One important invariant here, is that whenever a (addr, slot) is added, if the
		addr is not already present, the add causes two journal entries:
		- one for the address,
		- one for the (address,slot)
		Therefore, when unrolling the change, we can always blindly delete the
		(addr) at this point, since no storage adds can remain when come upon
		a single (addr) change.
	*/
	s.accessList.DeleteAddress(*ch.address)
}

func (ch accessListAddAccountChange) dirtied() *ethcmn.Address {
	return nil
}

func (ch accessListAddSlotChange) revert(s *CommitStateDB) {
	s.accessList.DeleteSlot(*ch.address, *ch.slot)
}

func (ch accessListAddSlotChange) dirtied() *ethcmn.Address {
	return nil
}
