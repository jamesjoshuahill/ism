package kube

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type BrokerRegisterTimeoutErr struct {
	brokerName string
}

func (e BrokerRegisterTimeoutErr) Error() string {
	return fmt.Sprintf("timed out waiting for broker '%s' to be registered", e.brokerName)
}

type Broker struct {
	KubeClient          client.Client
	RegistrationTimeout time.Duration
}

func (b *Broker) FindAll() ([]*osbapi.Broker, error) {
	list := &v1alpha1.BrokerList{}
	if err := b.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Broker{}, err
	}

	brokers := []*osbapi.Broker{}
	for _, broker := range list.Items {
		brokers = append(brokers, &osbapi.Broker{
			Name:      broker.Spec.Name,
			URL:       broker.Spec.URL,
			Username:  broker.Spec.Username,
			Password:  broker.Spec.Password,
			CreatedAt: broker.ObjectMeta.CreationTimestamp.String(),
		})
	}

	return brokers, nil
}

func (b *Broker) Register(broker *osbapi.Broker) error {
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

	if err := b.KubeClient.Create(context.TODO(), brokerResource); err != nil {
		return err
	}

	return b.waitForBrokerRegistration(brokerResource)
}

func (b *Broker) waitForBrokerRegistration(broker *v1alpha1.Broker) error {
	err := wait.Poll(time.Second/2, b.RegistrationTimeout, func() (bool, error) {
		fetchedBroker := &v1alpha1.Broker{}

		err := b.KubeClient.Get(context.TODO(), types.NamespacedName{Name: broker.Name, Namespace: broker.Namespace}, fetchedBroker)
		if err == nil && fetchedBroker.Status.State == v1alpha1.BrokerStateRegistered {
			return true, nil
		}

		return false, nil
	})

	if err != nil {
		if err == wait.ErrWaitTimeout {
			return BrokerRegisterTimeoutErr{brokerName: broker.Name}
		}

		return err
	}

	return nil
}
