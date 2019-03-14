// Code generated by counterfeiter. DO NOT EDIT.
package actorsfakes

import (
	sync "sync"

	actors "github.com/pivotal-cf/ism/actors"
	osbapi "github.com/pivotal-cf/ism/osbapi"
)

type FakeServiceRepository struct {
	FindStub        func(string) (*osbapi.Service, error)
	findMutex       sync.RWMutex
	findArgsForCall []struct {
		arg1 string
	}
	findReturns struct {
		result1 *osbapi.Service
		result2 error
	}
	findReturnsOnCall map[int]struct {
		result1 *osbapi.Service
		result2 error
	}
	FindByBrokerStub        func(string) ([]*osbapi.Service, error)
	findByBrokerMutex       sync.RWMutex
	findByBrokerArgsForCall []struct {
		arg1 string
	}
	findByBrokerReturns struct {
		result1 []*osbapi.Service
		result2 error
	}
	findByBrokerReturnsOnCall map[int]struct {
		result1 []*osbapi.Service
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeServiceRepository) Find(arg1 string) (*osbapi.Service, error) {
	fake.findMutex.Lock()
	ret, specificReturn := fake.findReturnsOnCall[len(fake.findArgsForCall)]
	fake.findArgsForCall = append(fake.findArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Find", []interface{}{arg1})
	fake.findMutex.Unlock()
	if fake.FindStub != nil {
		return fake.FindStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServiceRepository) FindCallCount() int {
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	return len(fake.findArgsForCall)
}

func (fake *FakeServiceRepository) FindCalls(stub func(string) (*osbapi.Service, error)) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = stub
}

func (fake *FakeServiceRepository) FindArgsForCall(i int) string {
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	argsForCall := fake.findArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeServiceRepository) FindReturns(result1 *osbapi.Service, result2 error) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = nil
	fake.findReturns = struct {
		result1 *osbapi.Service
		result2 error
	}{result1, result2}
}

func (fake *FakeServiceRepository) FindReturnsOnCall(i int, result1 *osbapi.Service, result2 error) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = nil
	if fake.findReturnsOnCall == nil {
		fake.findReturnsOnCall = make(map[int]struct {
			result1 *osbapi.Service
			result2 error
		})
	}
	fake.findReturnsOnCall[i] = struct {
		result1 *osbapi.Service
		result2 error
	}{result1, result2}
}

func (fake *FakeServiceRepository) FindByBroker(arg1 string) ([]*osbapi.Service, error) {
	fake.findByBrokerMutex.Lock()
	ret, specificReturn := fake.findByBrokerReturnsOnCall[len(fake.findByBrokerArgsForCall)]
	fake.findByBrokerArgsForCall = append(fake.findByBrokerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FindByBroker", []interface{}{arg1})
	fake.findByBrokerMutex.Unlock()
	if fake.FindByBrokerStub != nil {
		return fake.FindByBrokerStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findByBrokerReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServiceRepository) FindByBrokerCallCount() int {
	fake.findByBrokerMutex.RLock()
	defer fake.findByBrokerMutex.RUnlock()
	return len(fake.findByBrokerArgsForCall)
}

func (fake *FakeServiceRepository) FindByBrokerCalls(stub func(string) ([]*osbapi.Service, error)) {
	fake.findByBrokerMutex.Lock()
	defer fake.findByBrokerMutex.Unlock()
	fake.FindByBrokerStub = stub
}

func (fake *FakeServiceRepository) FindByBrokerArgsForCall(i int) string {
	fake.findByBrokerMutex.RLock()
	defer fake.findByBrokerMutex.RUnlock()
	argsForCall := fake.findByBrokerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeServiceRepository) FindByBrokerReturns(result1 []*osbapi.Service, result2 error) {
	fake.findByBrokerMutex.Lock()
	defer fake.findByBrokerMutex.Unlock()
	fake.FindByBrokerStub = nil
	fake.findByBrokerReturns = struct {
		result1 []*osbapi.Service
		result2 error
	}{result1, result2}
}

func (fake *FakeServiceRepository) FindByBrokerReturnsOnCall(i int, result1 []*osbapi.Service, result2 error) {
	fake.findByBrokerMutex.Lock()
	defer fake.findByBrokerMutex.Unlock()
	fake.FindByBrokerStub = nil
	if fake.findByBrokerReturnsOnCall == nil {
		fake.findByBrokerReturnsOnCall = make(map[int]struct {
			result1 []*osbapi.Service
			result2 error
		})
	}
	fake.findByBrokerReturnsOnCall[i] = struct {
		result1 []*osbapi.Service
		result2 error
	}{result1, result2}
}

func (fake *FakeServiceRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	fake.findByBrokerMutex.RLock()
	defer fake.findByBrokerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeServiceRepository) recordInvocation(key string, args []interface{}) {
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

var _ actors.ServiceRepository = new(FakeServiceRepository)
