package actors

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . InstanceRepository

type InstanceRepository interface {
	Create(*osbapi.Instance) error
}

type InstancesActor struct {
	Repository InstanceRepository
}

func (a *InstancesActor) Create(instance *osbapi.Instance) error {
	return a.Repository.Create(instance)
}
