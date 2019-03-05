package kube

import (
	"context"
	"errors"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type Broker struct {
	KubeClient client.Client
}

//TODO change r to b
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

	if err := r.KubeClient.Create(context.TODO(), brokerResource); err != nil {
		return err
	}

	return r.waitForBrokerRegistration(brokerResource)
}

func (r *Broker) waitForBrokerRegistration(broker *v1alpha1.Broker) error {
	err := wait.Poll(1*time.Second, 10*time.Second, func() (bool, error) {
		b := &v1alpha1.Broker{}

		err := r.KubeClient.Get(context.TODO(), types.NamespacedName{Name: broker.Name, Namespace: broker.Namespace}, b)
		if err == nil {
			if b.Status.State == v1alpha1.BrokerStateRegistered {
				return true, nil
			}
		}

		return false, nil
	})

	if err != nil {
		if err == wait.ErrWaitTimeout {
			return errors.New("timed out waiting for broker to be registered")
		}

		return err
	}

	return nil
}
