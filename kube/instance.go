package kube

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type InstanceCreateTimeoutErr struct {
	instanceName string
}

func (e InstanceCreateTimeoutErr) Error() string {
	return fmt.Sprintf("timed out waiting for instance '%s' to be created", e.instanceName)
}

type Instance struct {
	KubeClient client.Client
	// CreationTimeout time.Duration
}

func (b *Instance) Create(instance *osbapi.Instance) error {
	instanceResource := &v1alpha1.ServiceInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: "default",
		},
		Spec: v1alpha1.ServiceInstanceSpec{
			Name:       instance.Name,
			PlanID:     instance.PlanID,
			ServiceID:  instance.ServiceID,
			BrokerName: instance.BrokerName,
		},
	}

	if err := b.KubeClient.Create(context.TODO(), instanceResource); err != nil {
		return err
	}
	return nil

	// return b.waitForInstanceRegistration(instanceResource)
}

// func (b *Instance) waitForInstanceRegistration(instance *v1alpha1.Instance) error {
// 	err := wait.Poll(time.Second/2, b.RegistrationTimeout, func() (bool, error) {
// 		fetchedInstance := &v1alpha1.Instance{}
//
// 		err := b.KubeClient.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, fetchedInstance)
// 		if err == nil && fetchedInstance.Status.State == v1alpha1.InstanceStateCreateed {
// 			return true, nil
// 		}
//
// 		return false, nil
// 	})
//
// 	if err != nil {
// 		if err == wait.ErrWaitTimeout {
// 			return InstanceCreateTimeoutErr{instanceName: instance.Name}
// 		}
//
// 		return err
// 	}
//
// 	return nil
// }
