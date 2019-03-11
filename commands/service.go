package commands

import (
	"strings"

	"github.com/pivotal-cf/ism/usecases"
)

//go:generate counterfeiter . ServiceListUsecase

type ServiceListUsecase interface {
	GetServices() ([]*usecases.Service, error)
}

type ServiceCommand struct {
	ServiceListCommand ServiceListCommand `command:"list" long-description:"List the services that are available in the marketplace."`
}

type ServiceListCommand struct {
	UI                 UI
	ServiceListUsecase ServiceListUsecase
}

func (cmd *ServiceListCommand) Execute([]string) error {
	services, err := cmd.ServiceListUsecase.GetServices()
	if err != nil {
		return err
	}

	if len(services) == 0 {
		cmd.UI.DisplayText("No services found.")
		return nil
	}

	servicesTable := buildServiceTableData(services)
	cmd.UI.DisplayTable(servicesTable)

	return nil
}

func buildServiceTableData(services []*usecases.Service) [][]string {
	headers := []string{"SERVICE", "PLANS", "BROKER", "DESCRIPTION"}
	data := [][]string{headers}

	for _, s := range services {
		row := []string{s.Name, strings.Join(s.PlanNames, ", "), s.BrokerName, s.Description}
		data = append(data, row)
	}

	return data
}
