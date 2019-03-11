package commands

//go:generate counterfeiter . InstanceCreateUsecase

type InstanceCreateUsecase interface {
	Create(name, planName, serviceName, brokerName string) error
}

type InstanceCommand struct {
	InstanceCreateCommand InstanceCreateCommand `command:"create" long-description:"Create a service instance"`
}

type InstanceCreateCommand struct {
	Name    string `long:"name" description:"Name of the service instance" required:"true"`
	Service string `long:"service" description:"Name of the service" required:"true"`
	Plan    string `long:"plan" description:"Name of the plan" required:"true"`
	Broker  string `long:"broker" description:"Name of the broker" required:"true"`

	UI                    UI
	InstanceCreateUsecase InstanceCreateUsecase
}

func (cmd *InstanceCreateCommand) Execute([]string) error {
	if err := cmd.InstanceCreateUsecase.Create(cmd.Name, cmd.Plan, cmd.Service, cmd.Broker); err != nil {
		return err
	}

	cmd.UI.DisplayText("Instance '{{.InstanceName}}' created.", map[string]interface{}{"InstanceName": cmd.Name})

	return nil
}
