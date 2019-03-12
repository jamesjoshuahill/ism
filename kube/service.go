package kube

import (
	"context"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type Service struct {
	KubeClient client.Client
}

func (s *Service) FindByBroker(brokerName string) ([]*osbapi.Service, error) {
	list := &v1alpha1.BrokerServiceList{}
	if err := s.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Service{}, err
	}

	services := []*osbapi.Service{}
	for _, s := range list.Items {
		// TODO: This code will be refactored so filtering happens in the API - for now
		// we are assuming there will never be multiple owner references. See #164327846
		if len(s.ObjectMeta.OwnerReferences) == 0 {
			break
		}

		if string(s.ObjectMeta.OwnerReferences[0].Name) == brokerName {
			services = append(services, &osbapi.Service{
				ID:          s.ObjectMeta.Name,
				Name:        s.Spec.Name,
				Description: s.Spec.Description,
				BrokerName:  brokerName,
			})
		}
	}

	return services, nil
}
