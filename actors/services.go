package actors

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . ServiceRepository

type ServiceRepository interface {
	FindByBroker(brokerID string) ([]*osbapi.Service, error)
}

type ServicesActor struct {
	Repository ServiceRepository
}

func (a *ServicesActor) GetServices(brokerID string) ([]*osbapi.Service, error) {
	return a.Repository.FindByBroker(brokerID)
}
