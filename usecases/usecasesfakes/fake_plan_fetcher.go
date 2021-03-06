// Code generated by counterfeiter. DO NOT EDIT.
package usecasesfakes

import (
	sync "sync"

	osbapi "github.com/pivotal-cf/ism/osbapi"
	usecases "github.com/pivotal-cf/ism/usecases"
)

type FakePlanFetcher struct {
	GetPlanStub        func(string) (*osbapi.Plan, error)
	getPlanMutex       sync.RWMutex
	getPlanArgsForCall []struct {
		arg1 string
	}
	getPlanReturns struct {
		result1 *osbapi.Plan
		result2 error
	}
	getPlanReturnsOnCall map[int]struct {
		result1 *osbapi.Plan
		result2 error
	}
	GetPlansStub        func(string) ([]*osbapi.Plan, error)
	getPlansMutex       sync.RWMutex
	getPlansArgsForCall []struct {
		arg1 string
	}
	getPlansReturns struct {
		result1 []*osbapi.Plan
		result2 error
	}
	getPlansReturnsOnCall map[int]struct {
		result1 []*osbapi.Plan
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePlanFetcher) GetPlan(arg1 string) (*osbapi.Plan, error) {
	fake.getPlanMutex.Lock()
	ret, specificReturn := fake.getPlanReturnsOnCall[len(fake.getPlanArgsForCall)]
	fake.getPlanArgsForCall = append(fake.getPlanArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetPlan", []interface{}{arg1})
	fake.getPlanMutex.Unlock()
	if fake.GetPlanStub != nil {
		return fake.GetPlanStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPlanReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePlanFetcher) GetPlanCallCount() int {
	fake.getPlanMutex.RLock()
	defer fake.getPlanMutex.RUnlock()
	return len(fake.getPlanArgsForCall)
}

func (fake *FakePlanFetcher) GetPlanCalls(stub func(string) (*osbapi.Plan, error)) {
	fake.getPlanMutex.Lock()
	defer fake.getPlanMutex.Unlock()
	fake.GetPlanStub = stub
}

func (fake *FakePlanFetcher) GetPlanArgsForCall(i int) string {
	fake.getPlanMutex.RLock()
	defer fake.getPlanMutex.RUnlock()
	argsForCall := fake.getPlanArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakePlanFetcher) GetPlanReturns(result1 *osbapi.Plan, result2 error) {
	fake.getPlanMutex.Lock()
	defer fake.getPlanMutex.Unlock()
	fake.GetPlanStub = nil
	fake.getPlanReturns = struct {
		result1 *osbapi.Plan
		result2 error
	}{result1, result2}
}

func (fake *FakePlanFetcher) GetPlanReturnsOnCall(i int, result1 *osbapi.Plan, result2 error) {
	fake.getPlanMutex.Lock()
	defer fake.getPlanMutex.Unlock()
	fake.GetPlanStub = nil
	if fake.getPlanReturnsOnCall == nil {
		fake.getPlanReturnsOnCall = make(map[int]struct {
			result1 *osbapi.Plan
			result2 error
		})
	}
	fake.getPlanReturnsOnCall[i] = struct {
		result1 *osbapi.Plan
		result2 error
	}{result1, result2}
}

func (fake *FakePlanFetcher) GetPlans(arg1 string) ([]*osbapi.Plan, error) {
	fake.getPlansMutex.Lock()
	ret, specificReturn := fake.getPlansReturnsOnCall[len(fake.getPlansArgsForCall)]
	fake.getPlansArgsForCall = append(fake.getPlansArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetPlans", []interface{}{arg1})
	fake.getPlansMutex.Unlock()
	if fake.GetPlansStub != nil {
		return fake.GetPlansStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPlansReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePlanFetcher) GetPlansCallCount() int {
	fake.getPlansMutex.RLock()
	defer fake.getPlansMutex.RUnlock()
	return len(fake.getPlansArgsForCall)
}

func (fake *FakePlanFetcher) GetPlansCalls(stub func(string) ([]*osbapi.Plan, error)) {
	fake.getPlansMutex.Lock()
	defer fake.getPlansMutex.Unlock()
	fake.GetPlansStub = stub
}

func (fake *FakePlanFetcher) GetPlansArgsForCall(i int) string {
	fake.getPlansMutex.RLock()
	defer fake.getPlansMutex.RUnlock()
	argsForCall := fake.getPlansArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakePlanFetcher) GetPlansReturns(result1 []*osbapi.Plan, result2 error) {
	fake.getPlansMutex.Lock()
	defer fake.getPlansMutex.Unlock()
	fake.GetPlansStub = nil
	fake.getPlansReturns = struct {
		result1 []*osbapi.Plan
		result2 error
	}{result1, result2}
}

func (fake *FakePlanFetcher) GetPlansReturnsOnCall(i int, result1 []*osbapi.Plan, result2 error) {
	fake.getPlansMutex.Lock()
	defer fake.getPlansMutex.Unlock()
	fake.GetPlansStub = nil
	if fake.getPlansReturnsOnCall == nil {
		fake.getPlansReturnsOnCall = make(map[int]struct {
			result1 []*osbapi.Plan
			result2 error
		})
	}
	fake.getPlansReturnsOnCall[i] = struct {
		result1 []*osbapi.Plan
		result2 error
	}{result1, result2}
}

func (fake *FakePlanFetcher) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getPlanMutex.RLock()
	defer fake.getPlanMutex.RUnlock()
	fake.getPlansMutex.RLock()
	defer fake.getPlansMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePlanFetcher) recordInvocation(key string, args []interface{}) {
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

var _ usecases.PlanFetcher = new(FakePlanFetcher)
