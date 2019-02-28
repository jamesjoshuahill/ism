package commands

type RootCommand struct {
	BrokerCommand  BrokerCommand  `command:"broker" long-description:"The broker command group lets you register and list service brokers from the marketplace"`
	ServiceCommand ServiceCommand `command:"service" long-description:"The service command group lets you list the available services in the marketplace"`
}
