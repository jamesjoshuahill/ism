// Code generated by counterfeiter. DO NOT EDIT.
package actorsfakes

import (
	"sync"

	"github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/osbapi"
)

type FakeBrokerRepository struct {
	DeleteStub        func(string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 string
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	FindAllStub        func() ([]*osbapi.Broker, error)
	findAllMutex       sync.RWMutex
	findAllArgsForCall []struct {
	}
	findAllReturns struct {
		result1 []*osbapi.Broker
		result2 error
	}
	findAllReturnsOnCall map[int]struct {
		result1 []*osbapi.Broker
		result2 error
	}
	FindByNameStub        func(string) (*osbapi.Broker, error)
	findByNameMutex       sync.RWMutex
	findByNameArgsForCall []struct {
		arg1 string
	}
	findByNameReturns struct {
		result1 *osbapi.Broker
		result2 error
	}
	findByNameReturnsOnCall map[int]struct {
		result1 *osbapi.Broker
		result2 error
	}
	RegisterStub        func(*osbapi.Broker) error
	registerMutex       sync.RWMutex
	registerArgsForCall []struct {
		arg1 *osbapi.Broker
	}
	registerReturns struct {
		result1 error
	}
	registerReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBrokerRepository) Delete(arg1 string) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Delete", []interface{}{arg1})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteReturns
	return fakeReturns.result1
}

func (fake *FakeBrokerRepository) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeBrokerRepository) DeleteCalls(stub func(string) error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeBrokerRepository) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBrokerRepository) DeleteReturns(result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBrokerRepository) DeleteReturnsOnCall(i int, result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBrokerRepository) FindAll() ([]*osbapi.Broker, error) {
	fake.findAllMutex.Lock()
	ret, specificReturn := fake.findAllReturnsOnCall[len(fake.findAllArgsForCall)]
	fake.findAllArgsForCall = append(fake.findAllArgsForCall, struct {
	}{})
	fake.recordInvocation("FindAll", []interface{}{})
	fake.findAllMutex.Unlock()
	if fake.FindAllStub != nil {
		return fake.FindAllStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findAllReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBrokerRepository) FindAllCallCount() int {
	fake.findAllMutex.RLock()
	defer fake.findAllMutex.RUnlock()
	return len(fake.findAllArgsForCall)
}

func (fake *FakeBrokerRepository) FindAllCalls(stub func() ([]*osbapi.Broker, error)) {
	fake.findAllMutex.Lock()
	defer fake.findAllMutex.Unlock()
	fake.FindAllStub = stub
}

func (fake *FakeBrokerRepository) FindAllReturns(result1 []*osbapi.Broker, result2 error) {
	fake.findAllMutex.Lock()
	defer fake.findAllMutex.Unlock()
	fake.FindAllStub = nil
	fake.findAllReturns = struct {
		result1 []*osbapi.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBrokerRepository) FindAllReturnsOnCall(i int, result1 []*osbapi.Broker, result2 error) {
	fake.findAllMutex.Lock()
	defer fake.findAllMutex.Unlock()
	fake.FindAllStub = nil
	if fake.findAllReturnsOnCall == nil {
		fake.findAllReturnsOnCall = make(map[int]struct {
			result1 []*osbapi.Broker
			result2 error
		})
	}
	fake.findAllReturnsOnCall[i] = struct {
		result1 []*osbapi.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBrokerRepository) FindByName(arg1 string) (*osbapi.Broker, error) {
	fake.findByNameMutex.Lock()
	ret, specificReturn := fake.findByNameReturnsOnCall[len(fake.findByNameArgsForCall)]
	fake.findByNameArgsForCall = append(fake.findByNameArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FindByName", []interface{}{arg1})
	fake.findByNameMutex.Unlock()
	if fake.FindByNameStub != nil {
		return fake.FindByNameStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findByNameReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBrokerRepository) FindByNameCallCount() int {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	return len(fake.findByNameArgsForCall)
}

func (fake *FakeBrokerRepository) FindByNameCalls(stub func(string) (*osbapi.Broker, error)) {
	fake.findByNameMutex.Lock()
	defer fake.findByNameMutex.Unlock()
	fake.FindByNameStub = stub
}

func (fake *FakeBrokerRepository) FindByNameArgsForCall(i int) string {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	argsForCall := fake.findByNameArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBrokerRepository) FindByNameReturns(result1 *osbapi.Broker, result2 error) {
	fake.findByNameMutex.Lock()
	defer fake.findByNameMutex.Unlock()
	fake.FindByNameStub = nil
	fake.findByNameReturns = struct {
		result1 *osbapi.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBrokerRepository) FindByNameReturnsOnCall(i int, result1 *osbapi.Broker, result2 error) {
	fake.findByNameMutex.Lock()
	defer fake.findByNameMutex.Unlock()
	fake.FindByNameStub = nil
	if fake.findByNameReturnsOnCall == nil {
		fake.findByNameReturnsOnCall = make(map[int]struct {
			result1 *osbapi.Broker
			result2 error
		})
	}
	fake.findByNameReturnsOnCall[i] = struct {
		result1 *osbapi.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBrokerRepository) Register(arg1 *osbapi.Broker) error {
	fake.registerMutex.Lock()
	ret, specificReturn := fake.registerReturnsOnCall[len(fake.registerArgsForCall)]
	fake.registerArgsForCall = append(fake.registerArgsForCall, struct {
		arg1 *osbapi.Broker
	}{arg1})
	fake.recordInvocation("Register", []interface{}{arg1})
	fake.registerMutex.Unlock()
	if fake.RegisterStub != nil {
		return fake.RegisterStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.registerReturns
	return fakeReturns.result1
}

func (fake *FakeBrokerRepository) RegisterCallCount() int {
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	return len(fake.registerArgsForCall)
}

func (fake *FakeBrokerRepository) RegisterCalls(stub func(*osbapi.Broker) error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = stub
}

func (fake *FakeBrokerRepository) RegisterArgsForCall(i int) *osbapi.Broker {
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	argsForCall := fake.registerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBrokerRepository) RegisterReturns(result1 error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = nil
	fake.registerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBrokerRepository) RegisterReturnsOnCall(i int, result1 error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = nil
	if fake.registerReturnsOnCall == nil {
		fake.registerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.registerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBrokerRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.findAllMutex.RLock()
	defer fake.findAllMutex.RUnlock()
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBrokerRepository) recordInvocation(key string, args []interface{}) {
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

var _ actors.BrokerRepository = new(FakeBrokerRepository)
