package stubwrapper

import (
	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/cc-tools/mock"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/util"
)

type StubWrapper struct {
	Stub        shim.ChaincodeStubInterface
	WriteSet    map[string][]byte
	PvtWriteSet map[string]map[string][]byte
}

func (sw *StubWrapper) PutState(key string, obj []byte) errors.ICCError {
	err := sw.Stub.PutState(key, obj)
	if err != nil {
		return errors.WrapError(err, "stub.PutState call error")
	}

	if sw.WriteSet == nil {
		sw.WriteSet = make(map[string][]byte)
	}
	sw.WriteSet[key] = obj

	return nil
}

func (sw *StubWrapper) GetState(key string) ([]byte, errors.ICCError) {
	obj, inSet := sw.WriteSet[key]
	if inSet {
		return obj, nil
	}

	return sw.GetCommittedState(key)
}

func (sw *StubWrapper) GetCommittedState(key string) ([]byte, errors.ICCError) {
	obj, err := sw.Stub.GetState(key)
	if err != nil {
		return nil, errors.WrapError(err, "stub.GetState call error")
	}

	return obj, nil
}

func (sw *StubWrapper) DelState(key string) errors.ICCError {
	err := sw.Stub.DelState(key)
	if err != nil {
		return errors.WrapError(err, "stub.DelState call error")
	}

	if sw.WriteSet == nil {
		sw.WriteSet = make(map[string][]byte)
	}
	sw.WriteSet[key] = nil

	return nil
}

func (sw *StubWrapper) PutPrivateData(collection, key string, obj []byte) errors.ICCError {
	err := sw.Stub.PutPrivateData(collection, key, obj)
	if err != nil {
		return errors.WrapError(err, "stub.PutPrivateData call error")
	}

	if sw.PvtWriteSet == nil {
		sw.PvtWriteSet = make(map[string]map[string][]byte)
	}
	if sw.PvtWriteSet[collection] == nil {
		sw.PvtWriteSet[collection] = make(map[string][]byte)
	}

	sw.PvtWriteSet[collection][key] = obj

	return nil
}

func (sw *StubWrapper) GetPrivateData(collection, key string) ([]byte, errors.ICCError) {
	obj, inSet := sw.PvtWriteSet[collection][key]
	if inSet {
		return obj, nil
	}

	return sw.GetCommittedPrivateData(collection, key)
}

func (sw *StubWrapper) GetCommittedPrivateData(collection, key string) ([]byte, errors.ICCError) {

	obj, err := sw.Stub.GetPrivateData(collection, key)
	if err != nil {
		return nil, errors.WrapError(err, "stub.GetPrivateData call error")
	}

	return obj, nil
}

func (sw *StubWrapper) GetPrivateDataHash(collection, key string) ([]byte, errors.ICCError) {
	obj, inSet := sw.PvtWriteSet[collection][key]
	if inSet {
		if obj != nil {
			return util.ComputeSHA256(obj), nil
		} else {
			return nil, nil
		}
	}

	obj, err := sw.Stub.GetPrivateDataHash(collection, key)
	if err != nil {
		return nil, errors.WrapError(err, "stub.GetPrivateData call error")
	}

	return obj, nil
}

func (sw *StubWrapper) DelPrivateData(collection, key string) errors.ICCError {
	err := sw.Stub.DelPrivateData(collection, key)
	if err != nil {
		return errors.WrapError(err, "stub.DelPrivateData call error")
	}

	if sw.PvtWriteSet == nil {
		sw.PvtWriteSet = make(map[string]map[string][]byte)
	}
	if sw.PvtWriteSet[collection] == nil {
		sw.PvtWriteSet[collection] = make(map[string][]byte)
	}

	sw.WriteSet[key] = nil

	return nil
}

func (sw *StubWrapper) CreateCompositeKey(objectType string, attributes []string) (string, errors.ICCError) {
	compositeKey, err := sw.Stub.CreateCompositeKey(objectType, attributes)
	if err != nil {
		return compositeKey, errors.WrapError(err, "stub.CreateCompositeKey call error")
	}
	return compositeKey, nil
}

// GetQueryResult does not return non-commited ledger states
func (sw *StubWrapper) GetQueryResult(query string) (shim.StateQueryIteratorInterface, errors.ICCError) {
	it, err := sw.Stub.GetQueryResult(query)
	if err != nil {
		return it, errors.WrapError(err, "stub.GetQueryResult call error")
	}
	return it, nil
}

// GetPrivateDataQueryResult does not return non-commited ledger states
func (sw *StubWrapper) GetPrivateDataQueryResult(collection, query string) (shim.StateQueryIteratorInterface, errors.ICCError) {
	it, err := sw.Stub.GetPrivateDataQueryResult(collection, query)
	if err != nil {
		return it, errors.WrapError(err, "stub.GetPrivateDataQueryResult call error")
	}
	return it, nil
}

// GetQueryResultWithPagination does not return non-commited ledger states
func (sw *StubWrapper) GetQueryResultWithPagination(query string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, errors.ICCError) {

	it, metadata, err := sw.Stub.GetQueryResultWithPagination(query, pageSize, bookmark)
	if err != nil {
		return it, metadata, errors.WrapError(err, "stub.GetQueryResultWithPagination call error")
	}
	return it, metadata, nil
}

// GetStateByPartialCompositeKey does not return non-commited ledger states
func (sw *StubWrapper) GetStateByPartialCompositeKey(objectType string, keys []string) (shim.StateQueryIteratorInterface, errors.ICCError) {
	it, err := sw.Stub.GetStateByPartialCompositeKey(objectType, keys)
	if err != nil {
		return it, errors.WrapError(err, "stub.GetStateByPartialCompositeKey call error")
	}
	return it, nil
}

// GetHistoryForKey does not return non-commited ledger states
func (sw *StubWrapper) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, errors.ICCError) {
	it, err := sw.Stub.GetHistoryForKey(key)
	if err != nil {
		return it, errors.WrapError(err, "stub.GetHistoryForKey call error")
	}
	return it, nil
}

// GetMSPID wraps cid.GetMSPID allowing for automated testing
func (sw *StubWrapper) GetMSPID() (string, errors.ICCError) {
	mockStub, isMock := sw.Stub.(*mock.MockStub)
	if isMock {
		return mockStub.Name, nil
	}
	mspid, err := cid.GetMSPID(sw.Stub)
	if err != nil {
		return mspid, errors.WrapError(err, "cid.GetMSPID call error")
	}
	return mspid, nil
}

// SplitCompositeKey returns composite keys
func (sw *StubWrapper) SplitCompositeKey(compositeKey string) (string, []string, errors.ICCError) {
	key, keys, err := sw.Stub.SplitCompositeKey(compositeKey)
	if err != nil {
		return "", nil, errors.WrapError(err, "stub.SplitCompositeKey call error")
	}
	return key, keys, nil
}

func (sw *StubWrapper) SetEvent(name string, payload []byte) errors.ICCError {
	err := sw.Stub.SetEvent(name, payload)
	if err != nil {
		return errors.WrapError(err, "stub.SetEvent call error")
	}

	return nil
}
