package actors

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . BrokerRepository

type BrokerRepository interface {
	FindAll() ([]*osbapi.Broker, error)
	Register(*osbapi.Broker) error
}

type BrokersActor struct {
	Repository BrokerRepository
}

func (a *BrokersActor) GetBrokers() ([]*osbapi.Broker, error) {
	return a.Repository.FindAll()
}

//TODO: Make names consistent
func (a *BrokersActor) Register(broker *osbapi.Broker) error {
	return a.Repository.Register(broker)
}
