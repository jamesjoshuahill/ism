// Code generated by counterfeiter. DO NOT EDIT.
package reconcilersfakes

import (
	sync "sync"

	v1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	reconcilers "github.com/pivotal-cf/ism/pkg/internal/reconcilers"
	v1 "k8s.io/api/core/v1"
)

type FakeKubeSecretRepo struct {
	CreateStub        func(*v1alpha1.ServiceBinding, map[string]interface{}) (*v1.Secret, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 *v1alpha1.ServiceBinding
		arg2 map[string]interface{}
	}
	createReturns struct {
		result1 *v1.Secret
		result2 error
	}
	createReturnsOnCall map[int]struct {
		result1 *v1.Secret
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeKubeSecretRepo) Create(arg1 *v1alpha1.ServiceBinding, arg2 map[string]interface{}) (*v1.Secret, error) {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 *v1alpha1.ServiceBinding
		arg2 map[string]interface{}
	}{arg1, arg2})
	fake.recordInvocation("Create", []interface{}{arg1, arg2})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeKubeSecretRepo) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeKubeSecretRepo) CreateCalls(stub func(*v1alpha1.ServiceBinding, map[string]interface{}) (*v1.Secret, error)) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = stub
}

func (fake *FakeKubeSecretRepo) CreateArgsForCall(i int) (*v1alpha1.ServiceBinding, map[string]interface{}) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	argsForCall := fake.createArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeKubeSecretRepo) CreateReturns(result1 *v1.Secret, result2 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 *v1.Secret
		result2 error
	}{result1, result2}
}

func (fake *FakeKubeSecretRepo) CreateReturnsOnCall(i int, result1 *v1.Secret, result2 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 *v1.Secret
			result2 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 *v1.Secret
		result2 error
	}{result1, result2}
}

func (fake *FakeKubeSecretRepo) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeKubeSecretRepo) recordInvocation(key string, args []interface{}) {
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

var _ reconcilers.KubeSecretRepo = new(FakeKubeSecretRepo)