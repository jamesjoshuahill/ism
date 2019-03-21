/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms
of the under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	flags "github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/commands"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"github.com/pivotal-cf/ism/repositories/kube"
	"github.com/pivotal-cf/ism/ui"
	"github.com/pivotal-cf/ism/usecases"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const timeout = time.Second * 30

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

	brokerRepository := &kube.Broker{KubeClient: kubeClient, RegistrationTimeout: timeout}
	serviceRepository := &kube.Service{KubeClient: kubeClient}
	planRepository := &kube.Plan{KubeClient: kubeClient}
	instanceRepository := &kube.Instance{KubeClient: kubeClient}
	bindingRepository := &kube.Binding{KubeClient: kubeClient}

	brokersActor := &actors.BrokersActor{
		Repository: brokerRepository,
	}
	servicesActor := &actors.ServicesActor{
		Repository: serviceRepository,
	}
	plansActor := &actors.PlansActor{
		Repository: planRepository,
	}
	instancesActor := &actors.InstancesActor{
		Repository: instanceRepository,
	}
	bindingsActor := &actors.BindingsActor{
		Repository: bindingRepository,
	}

	serviceListUsecase := &usecases.ServiceListUsecase{
		BrokerFetcher:  brokersActor,
		ServiceFetcher: servicesActor,
		PlanFetcher:    plansActor,
	}

	instanceCreateUsecase := &usecases.InstanceCreateUsecase{
		BrokerFetcher:   brokersActor,
		ServiceFetcher:  servicesActor,
		PlanFetcher:     plansActor,
		InstanceCreator: instancesActor,
	}

	instanceListUsecase := &usecases.InstanceListUsecase{
		InstanceFetcher: instancesActor,
		ServiceFetcher:  servicesActor,
		PlanFetcher:     plansActor,
	}

	bindingCreateUsecase := &usecases.BindingCreateUsecase{
		BindingCreator:  bindingsActor,
		InstanceFetcher: instancesActor,
	}

	bindingListUsecase := &usecases.BindingListUsecase{
		BindingFetcher:  bindingsActor,
		InstanceFetcher: instancesActor,
	}

	rootCommand := commands.RootCommand{
		BrokerCommand: commands.BrokerCommand{
			BrokerRegisterCommand: commands.BrokerRegisterCommand{
				UI:              UI,
				BrokerRegistrar: brokersActor,
			},
			BrokerListCommand: commands.BrokerListCommand{
				UI:             UI,
				BrokersFetcher: brokersActor,
			},
		},
		ServiceCommand: commands.ServiceCommand{
			ServiceListCommand: commands.ServiceListCommand{
				UI:                 UI,
				ServiceListUsecase: serviceListUsecase,
			},
		},
		InstanceCommand: commands.InstanceCommand{
			InstanceCreateCommand: commands.InstanceCreateCommand{
				UI:                    UI,
				InstanceCreateUsecase: instanceCreateUsecase,
			},
			InstanceListCommand: commands.InstanceListCommand{
				UI:                  UI,
				InstanceListUsecase: instanceListUsecase,
			},
		},
		BindingCommand: commands.BindingCommand{
			BindingCreateCommand: commands.BindingCreateCommand{
				UI:                   UI,
				BindingCreateUsecase: bindingCreateUsecase,
			},
			BindingListCommand: commands.BindingListCommand{
				UI:                 UI,
				BindingListUsecase: bindingListUsecase,
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

		if outErr, ok := err.(*flags.Error); ok && outErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
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
