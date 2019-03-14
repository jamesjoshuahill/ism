// Code generated by counterfeiter. DO NOT EDIT.
package commandsfakes

import (
	sync "sync"

	commands "github.com/pivotal-cf/ism/commands"
	usecases "github.com/pivotal-cf/ism/usecases"
)

type FakeInstanceListUsecase struct {
	GetInstancesStub        func() ([]*usecases.Instance, error)
	getInstancesMutex       sync.RWMutex
	getInstancesArgsForCall []struct {
	}
	getInstancesReturns struct {
		result1 []*usecases.Instance
		result2 error
	}
	getInstancesReturnsOnCall map[int]struct {
		result1 []*usecases.Instance
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeInstanceListUsecase) GetInstances() ([]*usecases.Instance, error) {
	fake.getInstancesMutex.Lock()
	ret, specificReturn := fake.getInstancesReturnsOnCall[len(fake.getInstancesArgsForCall)]
	fake.getInstancesArgsForCall = append(fake.getInstancesArgsForCall, struct {
	}{})
	fake.recordInvocation("GetInstances", []interface{}{})
	fake.getInstancesMutex.Unlock()
	if fake.GetInstancesStub != nil {
		return fake.GetInstancesStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getInstancesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeInstanceListUsecase) GetInstancesCallCount() int {
	fake.getInstancesMutex.RLock()
	defer fake.getInstancesMutex.RUnlock()
	return len(fake.getInstancesArgsForCall)
}

func (fake *FakeInstanceListUsecase) GetInstancesCalls(stub func() ([]*usecases.Instance, error)) {
	fake.getInstancesMutex.Lock()
	defer fake.getInstancesMutex.Unlock()
	fake.GetInstancesStub = stub
}

func (fake *FakeInstanceListUsecase) GetInstancesReturns(result1 []*usecases.Instance, result2 error) {
	fake.getInstancesMutex.Lock()
	defer fake.getInstancesMutex.Unlock()
	fake.GetInstancesStub = nil
	fake.getInstancesReturns = struct {
		result1 []*usecases.Instance
		result2 error
	}{result1, result2}
}

func (fake *FakeInstanceListUsecase) GetInstancesReturnsOnCall(i int, result1 []*usecases.Instance, result2 error) {
	fake.getInstancesMutex.Lock()
	defer fake.getInstancesMutex.Unlock()
	fake.GetInstancesStub = nil
	if fake.getInstancesReturnsOnCall == nil {
		fake.getInstancesReturnsOnCall = make(map[int]struct {
			result1 []*usecases.Instance
			result2 error
		})
	}
	fake.getInstancesReturnsOnCall[i] = struct {
		result1 []*usecases.Instance
		result2 error
	}{result1, result2}
}

func (fake *FakeInstanceListUsecase) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getInstancesMutex.RLock()
	defer fake.getInstancesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeInstanceListUsecase) recordInvocation(key string, args []interface{}) {
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

var _ commands.InstanceListUsecase = new(FakeInstanceListUsecase)
