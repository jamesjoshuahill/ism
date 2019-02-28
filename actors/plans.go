package actors

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . PlanRepository

type PlanRepository interface {
	FindByService(serviceID string) ([]*osbapi.Plan, error)
}

type PlansActor struct {
	Repository PlanRepository
}

func (a *PlansActor) GetPlans(serviceID string) ([]*osbapi.Plan, error) {
	return a.Repository.FindByService(serviceID)
}
