package commands

//TODO: Where should this be defined?

//go:generate counterfeiter . UI

type UI interface {
	DisplayText(text string, data ...map[string]interface{})
	DisplayTable(table [][]string)
}

type RootCommand struct {
	BrokerCommand   BrokerCommand   `command:"broker" long-description:"The broker command group lets you register and list service brokers from the marketplace"`
	ServiceCommand  ServiceCommand  `command:"service" long-description:"The service command group lets you list the available services in the marketplace"`
	InstanceCommand InstanceCommand `command:"instance" long-description:"The instance command group lets you create service instances"`
}
