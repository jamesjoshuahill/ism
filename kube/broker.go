package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type Broker struct {
	KubeClient client.Client
}

func (r *Broker) FindAll() ([]*osbapi.Broker, error) {
	list := &v1alpha1.BrokerList{}
	if err := r.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Broker{}, err
	}

	brokers := []*osbapi.Broker{}
	for _, b := range list.Items {
		brokers = append(brokers, &osbapi.Broker{
			ID:        string(b.UID),
			Name:      b.Spec.Name,
			URL:       b.Spec.URL,
			Username:  b.Spec.Username,
			Password:  b.Spec.Password,
			CreatedAt: b.ObjectMeta.CreationTimestamp.String(),
		})
	}

	return brokers, nil
}

func (r *Broker) Register(broker *osbapi.Broker) error {
	brokerResource := &v1alpha1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name:      broker.Name,
			Namespace: "default",
		},
		Spec: v1alpha1.BrokerSpec{
			Name:     broker.Name,
			URL:      broker.URL,
			Username: broker.Username,
			Password: broker.Password,
		},
	}

	return r.KubeClient.Create(context.TODO(), brokerResource)
}
