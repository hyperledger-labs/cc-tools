package mock

import (
	"errors"

	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
)

// StateQueryIterator implements shim.StateQueryIteratorInterface for rich queries
type StateQueryIterator struct {
	stub   *MockStub
	keys   []string
	pos    int
	closed bool
}

// HasNext returns true if the iterator contains more items
func (iter *StateQueryIterator) HasNext() bool {
	if iter.closed {
		return false
	}
	return iter.pos < len(iter.keys)
}

// Next returns the next key-value pair
func (iter *StateQueryIterator) Next() (*queryresult.KV, error) {
	if iter.closed {
		return nil, errors.New("StateQueryIterator.Next() called after Close()")
	}

	if !iter.HasNext() {
		return nil, errors.New("StateQueryIterator.Next() called when it does not have next")
	}

	// Get the key
	key := iter.keys[iter.pos]
	iter.pos++

	// Read the latest value from state
	value, err := iter.stub.GetState(key)
	if err != nil {
		return nil, err
	}

	// Skip keys that no longer exist in state
	if value == nil {
		if iter.HasNext() {
			return iter.Next()
		}
		return nil, errors.New("no more results")
	}

	return &queryresult.KV{
		Key:   key,
		Value: value,
	}, nil
}

// Close closes the iterator
func (iter *StateQueryIterator) Close() error {
	if iter.closed {
		return errors.New("StateQueryIterator.Close() called after Close()")
	}
	iter.closed = true
	iter.stub = nil
	iter.keys = nil
	return nil
}

// NewStateQueryIterator creates a new state query iterator
func NewStateQueryIterator(stub *MockStub, keys []string) *StateQueryIterator {
	return &StateQueryIterator{
		stub:   stub,
		keys:   keys,
		pos:    0,
		closed: false,
	}
}
