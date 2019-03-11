package kube

import (
	"context"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type Plan struct {
	KubeClient client.Client
}

func (p *Plan) FindByService(serviceID string) ([]*osbapi.Plan, error) {
	list := &v1alpha1.BrokerServicePlanList{}
	if err := p.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Plan{}, err
	}

	plans := []*osbapi.Plan{}
	for _, p := range list.Items {
		// TODO: This code will be refactored so filtering happens in the API - for now
		// we are assuming there will never be multiple owner references. See #164327846
		if len(p.ObjectMeta.OwnerReferences) == 0 {
			break
		}

		if string(p.ObjectMeta.OwnerReferences[0].UID) == serviceID {
			plans = append(plans, &osbapi.Plan{
				ID:        p.ObjectMeta.Name,
				Name:      p.Spec.Name,
				ServiceID: serviceID,
			})
		}
	}

	return plans, nil
}
