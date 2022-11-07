// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	types "github.com/tendermint/tendermint/abci/types"
)

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

// ApplySnapshotChunk provides a mock function with given fields: _a0
func (_m *Application) ApplySnapshotChunk(_a0 types.RequestApplySnapshotChunk) types.ResponseApplySnapshotChunk {
	ret := _m.Called(_a0)

	var r0 types.ResponseApplySnapshotChunk
	if rf, ok := ret.Get(0).(func(types.RequestApplySnapshotChunk) types.ResponseApplySnapshotChunk); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseApplySnapshotChunk)
	}

	return r0
}

// BeginBlock provides a mock function with given fields: _a0
func (_m *Application) BeginBlock(_a0 types.RequestBeginBlock) types.ResponseBeginBlock {
	ret := _m.Called(_a0)

	var r0 types.ResponseBeginBlock
	if rf, ok := ret.Get(0).(func(types.RequestBeginBlock) types.ResponseBeginBlock); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseBeginBlock)
	}

	return r0
}

// CheckTx provides a mock function with given fields: _a0
func (_m *Application) CheckTx(_a0 types.RequestCheckTx) types.ResponseCheckTx {
	ret := _m.Called(_a0)

	var r0 types.ResponseCheckTx
	if rf, ok := ret.Get(0).(func(types.RequestCheckTx) types.ResponseCheckTx); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseCheckTx)
	}

	return r0
}

// Commit provides a mock function with given fields:
func (_m *Application) Commit() types.ResponseCommit {
	ret := _m.Called()

	var r0 types.ResponseCommit
	if rf, ok := ret.Get(0).(func() types.ResponseCommit); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.ResponseCommit)
	}

	return r0
}

// DeliverTx provides a mock function with given fields: _a0
func (_m *Application) DeliverTx(_a0 types.RequestDeliverTx) types.ResponseDeliverTx {
	ret := _m.Called(_a0)

	var r0 types.ResponseDeliverTx
	if rf, ok := ret.Get(0).(func(types.RequestDeliverTx) types.ResponseDeliverTx); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseDeliverTx)
	}

	return r0
}

// EndBlock provides a mock function with given fields: _a0
func (_m *Application) EndBlock(_a0 types.RequestEndBlock) types.ResponseEndBlock {
	ret := _m.Called(_a0)

	var r0 types.ResponseEndBlock
	if rf, ok := ret.Get(0).(func(types.RequestEndBlock) types.ResponseEndBlock); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseEndBlock)
	}

	return r0
}

// ExtendVote provides a mock function with given fields: _a0
func (_m *Application) ExtendVote(_a0 types.RequestExtendVote) types.ResponseExtendVote {
	ret := _m.Called(_a0)

	var r0 types.ResponseExtendVote
	if rf, ok := ret.Get(0).(func(types.RequestExtendVote) types.ResponseExtendVote); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseExtendVote)
	}

	return r0
}

// Info provides a mock function with given fields: _a0
func (_m *Application) Info(_a0 types.RequestInfo) types.ResponseInfo {
	ret := _m.Called(_a0)

	var r0 types.ResponseInfo
	if rf, ok := ret.Get(0).(func(types.RequestInfo) types.ResponseInfo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseInfo)
	}

	return r0
}

// InitChain provides a mock function with given fields: _a0
func (_m *Application) InitChain(_a0 types.RequestInitChain) types.ResponseInitChain {
	ret := _m.Called(_a0)

	var r0 types.ResponseInitChain
	if rf, ok := ret.Get(0).(func(types.RequestInitChain) types.ResponseInitChain); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseInitChain)
	}

	return r0
}

// ListSnapshots provides a mock function with given fields: _a0
func (_m *Application) ListSnapshots(_a0 types.RequestListSnapshots) types.ResponseListSnapshots {
	ret := _m.Called(_a0)

	var r0 types.ResponseListSnapshots
	if rf, ok := ret.Get(0).(func(types.RequestListSnapshots) types.ResponseListSnapshots); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseListSnapshots)
	}

	return r0
}

// LoadSnapshotChunk provides a mock function with given fields: _a0
func (_m *Application) LoadSnapshotChunk(_a0 types.RequestLoadSnapshotChunk) types.ResponseLoadSnapshotChunk {
	ret := _m.Called(_a0)

	var r0 types.ResponseLoadSnapshotChunk
	if rf, ok := ret.Get(0).(func(types.RequestLoadSnapshotChunk) types.ResponseLoadSnapshotChunk); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseLoadSnapshotChunk)
	}

	return r0
}

// OfferSnapshot provides a mock function with given fields: _a0
func (_m *Application) OfferSnapshot(_a0 types.RequestOfferSnapshot) types.ResponseOfferSnapshot {
	ret := _m.Called(_a0)

	var r0 types.ResponseOfferSnapshot
	if rf, ok := ret.Get(0).(func(types.RequestOfferSnapshot) types.ResponseOfferSnapshot); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseOfferSnapshot)
	}

	return r0
}

// PrepareProposal provides a mock function with given fields: _a0
func (_m *Application) PrepareProposal(_a0 types.RequestPrepareProposal) types.ResponsePrepareProposal {
	ret := _m.Called(_a0)

	var r0 types.ResponsePrepareProposal
	if rf, ok := ret.Get(0).(func(types.RequestPrepareProposal) types.ResponsePrepareProposal); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponsePrepareProposal)
	}

	return r0
}

// ProcessProposal provides a mock function with given fields: _a0
func (_m *Application) ProcessProposal(_a0 types.RequestProcessProposal) types.ResponseProcessProposal {
	ret := _m.Called(_a0)

	var r0 types.ResponseProcessProposal
	if rf, ok := ret.Get(0).(func(types.RequestProcessProposal) types.ResponseProcessProposal); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseProcessProposal)
	}

	return r0
}

// Query provides a mock function with given fields: _a0
func (_m *Application) Query(_a0 types.RequestQuery) types.ResponseQuery {
	ret := _m.Called(_a0)

	var r0 types.ResponseQuery
	if rf, ok := ret.Get(0).(func(types.RequestQuery) types.ResponseQuery); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseQuery)
	}

	return r0
}

// VerifyVoteExtension provides a mock function with given fields: _a0
func (_m *Application) VerifyVoteExtension(_a0 types.RequestVerifyVoteExtension) types.ResponseVerifyVoteExtension {
	ret := _m.Called(_a0)

	var r0 types.ResponseVerifyVoteExtension
	if rf, ok := ret.Get(0).(func(types.RequestVerifyVoteExtension) types.ResponseVerifyVoteExtension); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseVerifyVoteExtension)
	}

	return r0
}

type mockConstructorTestingTNewApplication interface {
	mock.TestingT
	Cleanup(func())
}

// NewApplication creates a new instance of Application. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewApplication(t mockConstructorTestingTNewApplication) *Application {
	mock := &Application{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
