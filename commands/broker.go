package commands

import (
	"github.com/pivotal-cf/ism/osbapi"
)

//go:generate counterfeiter . BrokerRegistrar

type BrokerRegistrar interface {
	Register(*osbapi.Broker) error
}

//go:generate counterfeiter . BrokerFetcher

type BrokerFetcher interface {
	GetBrokers() ([]*osbapi.Broker, error)
}

type BrokerCommand struct {
	BrokerRegisterCommand BrokerRegisterCommand `command:"register" long-description:"Register a service broker into the marketplace"`
	BrokerListCommand     BrokerListCommand     `command:"list" long-description:"Lists the service brokers in the marketplace"`
}

type BrokerListCommand struct {
	UI            UI
	BrokerFetcher BrokerFetcher
}

type BrokerRegisterCommand struct {
	Name     string `long:"name" description:"Name of the service broker" required:"true"`
	URL      string `long:"url" description:"URL of the service broker" required:"true"`
	Username string `long:"username" description:"Username of the service broker" required:"true"`
	Password string `long:"password" description:"Password of the service broker" required:"true"`

	UI              UI
	BrokerRegistrar BrokerRegistrar
}

func (cmd *BrokerRegisterCommand) Execute([]string) error {
	newBroker := &osbapi.Broker{
		Name:     cmd.Name,
		URL:      cmd.URL,
		Username: cmd.Username,
		Password: cmd.Password,
	}

	if err := cmd.BrokerRegistrar.Register(newBroker); err != nil {
		return err
	}

	cmd.UI.DisplayText("Broker '{{.BrokerName}}' registered.", map[string]interface{}{"BrokerName": cmd.Name})

	return nil
}

func (cmd *BrokerListCommand) Execute([]string) error {
	brokers, err := cmd.BrokerFetcher.GetBrokers()
	if err != nil {
		return err
	}

	if len(brokers) == 0 {
		cmd.UI.DisplayText("No brokers found.")
		return nil
	}

	brokersTable := buildBrokerTableData(brokers)
	cmd.UI.DisplayTable(brokersTable)
	return nil
}

func buildBrokerTableData(brokers []*osbapi.Broker) [][]string {
	headers := []string{"NAME", "URL", "CREATED AT"}
	data := [][]string{headers}

	for _, b := range brokers {
		row := []string{b.Name, b.URL, b.CreatedAt}
		data = append(data, row)
	}
	return data
}
