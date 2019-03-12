package usecases

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . BrokerFetcher

type BrokerFetcher interface {
	GetBrokers() ([]*osbapi.Broker, error)
}

//go:generate counterfeiter . ServiceFetcher

type ServiceFetcher interface {
	GetServices(brokerName string) ([]*osbapi.Service, error)
}

//go:generate counterfeiter . PlanFetcher

type PlanFetcher interface {
	GetPlans(serviceID string) ([]*osbapi.Plan, error)
}

type ServiceListUsecase struct {
	BrokerFetcher  BrokerFetcher
	ServiceFetcher ServiceFetcher
	PlanFetcher    PlanFetcher
}

func (u *ServiceListUsecase) GetServices() ([]*Service, error) {
	brokers, err := u.BrokerFetcher.GetBrokers()
	if err != nil {
		return []*Service{}, err
	}

	var servicesToDisplay []*Service
	for _, b := range brokers {
		services, err := u.ServiceFetcher.GetServices(b.Name)
		if err != nil {
			return []*Service{}, err
		}

		for _, s := range services {
			plans, err := u.PlanFetcher.GetPlans(s.ID)
			if err != nil {
				return []*Service{}, err
			}

			serviceToDisplay := &Service{
				Name:        s.Name,
				Description: s.Description,
				PlanNames:   plansToNames(plans),
				BrokerName:  b.Name,
			}
			servicesToDisplay = append(servicesToDisplay, serviceToDisplay)
		}
	}

	return servicesToDisplay, nil
}

type Service struct {
	Name        string
	Description string
	PlanNames   []string
	BrokerName  string
}

func plansToNames(plans []*osbapi.Plan) []string {
	names := []string{}
	for _, p := range plans {
		names = append(names, p.Name)
	}

	return names
}
