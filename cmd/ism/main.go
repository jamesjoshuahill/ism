package main

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	flags "github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/kube"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"github.com/pivotal-cf/ism/ui"
	"github.com/pivotal-cf/ism/usecases"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	UI := &ui.UI{
		Out: os.Stdout,
		Err: os.Stderr,
	}

	kubeClient, err := buildKubeClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	brokerRepository := &kube.Broker{KubeClient: kubeClient}
	serviceRepository := &kube.Service{KubeClient: kubeClient}
	planRepository := &kube.Plan{KubeClient: kubeClient}

	brokersActor := &actors.BrokersActor{
		Repository: brokerRepository,
	}
	servicesActor := &actors.ServicesActor{
		Repository: serviceRepository,
	}
	plansActor := &actors.PlansActor{
		Repository: planRepository,
	}

	serviceListUsecase := &usecases.ServiceListUsecase{
		BrokerFetcher:  brokersActor,
		ServiceFetcher: servicesActor,
		PlanFetcher:    plansActor,
	}

	rootCommand := commands.RootCommand{
		BrokerCommand: commands.BrokerCommand{
			BrokerRegisterCommand: commands.BrokerRegisterCommand{
				UI:              UI,
				BrokerRegistrar: brokersActor,
			},
			BrokerListCommand: commands.BrokerListCommand{
				UI:            UI,
				BrokerFetcher: brokersActor,
			},
		},
		ServiceCommand: commands.ServiceCommand{
			ServiceListCommand: commands.ServiceListCommand{
				UI:                 UI,
				ServiceListUsecase: serviceListUsecase,
			},
		},
	}
	parser := flags.NewParser(&rootCommand, flags.HelpFlag|flags.PassDoubleDash)

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	_, err = parser.Parse()

	if err != nil {
		fmt.Println(err)
	}
}

func buildKubeClient() (client.Client, error) {
	home := os.Getenv("HOME")
	kubeconfigFilepath := fmt.Sprintf("%s/.kube/config", home)
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigFilepath)
	if err != nil {
		return nil, err
	}

	if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	return client.New(clientConfig, client.Options{Scheme: scheme.Scheme})
}
