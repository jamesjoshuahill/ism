// Code generated by counterfeiter. DO NOT EDIT.
package usecasesfakes

import (
	sync "sync"

	osbapi "github.com/pivotal-cf/ism/osbapi"
	usecases "github.com/pivotal-cf/ism/usecases"
)

type FakeBrokersActor struct {
	GetBrokersStub        func() ([]*osbapi.Broker, error)
	getBrokersMutex       sync.RWMutex
	getBrokersArgsForCall []struct {
	}
	getBrokersReturns struct {
		result1 []*osbapi.Broker
		result2 error
	}
	getBrokersReturnsOnCall map[int]struct {
		result1 []*osbapi.Broker
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBrokersActor) GetBrokers() ([]*osbapi.Broker, error) {
	fake.getBrokersMutex.Lock()
	ret, specificReturn := fake.getBrokersReturnsOnCall[len(fake.getBrokersArgsForCall)]
	fake.getBrokersArgsForCall = append(fake.getBrokersArgsForCall, struct {
	}{})
	fake.recordInvocation("GetBrokers", []interface{}{})
	fake.getBrokersMutex.Unlock()
	if fake.GetBrokersStub != nil {
		return fake.GetBrokersStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBrokersReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBrokersActor) GetBrokersCallCount() int {
	fake.getBrokersMutex.RLock()
	defer fake.getBrokersMutex.RUnlock()
	return len(fake.getBrokersArgsForCall)
}

func (fake *FakeBrokersActor) GetBrokersCalls(stub func() ([]*osbapi.Broker, error)) {
	fake.getBrokersMutex.Lock()
	defer fake.getBrokersMutex.Unlock()
	fake.GetBrokersStub = stub
}

func (fake *FakeBrokersActor) GetBrokersReturns(result1 []*osbapi.Broker, result2 error) {
	fake.getBrokersMutex.Lock()
	defer fake.getBrokersMutex.Unlock()
	fake.GetBrokersStub = nil
	fake.getBrokersReturns = struct {
		result1 []*osbapi.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBrokersActor) GetBrokersReturnsOnCall(i int, result1 []*osbapi.Broker, result2 error) {
	fake.getBrokersMutex.Lock()
	defer fake.getBrokersMutex.Unlock()
	fake.GetBrokersStub = nil
	if fake.getBrokersReturnsOnCall == nil {
		fake.getBrokersReturnsOnCall = make(map[int]struct {
			result1 []*osbapi.Broker
			result2 error
		})
	}
	fake.getBrokersReturnsOnCall[i] = struct {
		result1 []*osbapi.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBrokersActor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getBrokersMutex.RLock()
	defer fake.getBrokersMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBrokersActor) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ usecases.BrokersActor = new(FakeBrokersActor)
