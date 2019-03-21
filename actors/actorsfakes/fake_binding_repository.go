// Code generated by counterfeiter. DO NOT EDIT.
package actorsfakes

import (
	sync "sync"

	actors "github.com/pivotal-cf/ism/actors"
	osbapi "github.com/pivotal-cf/ism/osbapi"
)

type FakeBindingRepository struct {
	CreateStub        func(*osbapi.Binding) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 *osbapi.Binding
	}
	createReturns struct {
		result1 error
	}
	createReturnsOnCall map[int]struct {
		result1 error
	}
	FindAllStub        func() ([]*osbapi.Binding, error)
	findAllMutex       sync.RWMutex
	findAllArgsForCall []struct {
	}
	findAllReturns struct {
		result1 []*osbapi.Binding
		result2 error
	}
	findAllReturnsOnCall map[int]struct {
		result1 []*osbapi.Binding
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBindingRepository) Create(arg1 *osbapi.Binding) error {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 *osbapi.Binding
	}{arg1})
	fake.recordInvocation("Create", []interface{}{arg1})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.createReturns
	return fakeReturns.result1
}

func (fake *FakeBindingRepository) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeBindingRepository) CreateCalls(stub func(*osbapi.Binding) error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = stub
}

func (fake *FakeBindingRepository) CreateArgsForCall(i int) *osbapi.Binding {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	argsForCall := fake.createArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBindingRepository) CreateReturns(result1 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBindingRepository) CreateReturnsOnCall(i int, result1 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBindingRepository) FindAll() ([]*osbapi.Binding, error) {
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

func (fake *FakeBindingRepository) FindAllCallCount() int {
	fake.findAllMutex.RLock()
	defer fake.findAllMutex.RUnlock()
	return len(fake.findAllArgsForCall)
}

func (fake *FakeBindingRepository) FindAllCalls(stub func() ([]*osbapi.Binding, error)) {
	fake.findAllMutex.Lock()
	defer fake.findAllMutex.Unlock()
	fake.FindAllStub = stub
}

func (fake *FakeBindingRepository) FindAllReturns(result1 []*osbapi.Binding, result2 error) {
	fake.findAllMutex.Lock()
	defer fake.findAllMutex.Unlock()
	fake.FindAllStub = nil
	fake.findAllReturns = struct {
		result1 []*osbapi.Binding
		result2 error
	}{result1, result2}
}

func (fake *FakeBindingRepository) FindAllReturnsOnCall(i int, result1 []*osbapi.Binding, result2 error) {
	fake.findAllMutex.Lock()
	defer fake.findAllMutex.Unlock()
	fake.FindAllStub = nil
	if fake.findAllReturnsOnCall == nil {
		fake.findAllReturnsOnCall = make(map[int]struct {
			result1 []*osbapi.Binding
			result2 error
		})
	}
	fake.findAllReturnsOnCall[i] = struct {
		result1 []*osbapi.Binding
		result2 error
	}{result1, result2}
}

func (fake *FakeBindingRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.findAllMutex.RLock()
	defer fake.findAllMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBindingRepository) recordInvocation(key string, args []interface{}) {
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

var _ actors.BindingRepository = new(FakeBindingRepository)
