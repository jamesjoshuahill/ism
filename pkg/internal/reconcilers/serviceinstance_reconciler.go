package reconcilers

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	v1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// var ctx = context.TODO()

//go:generate counterfeiter . KubeServiceInstanceRepo

type KubeServiceInstanceRepo interface {
	Get(resource types.NamespacedName) (*v1alpha1.ServiceInstance, error)
	UpdateState(serviceinstance *v1alpha1.ServiceInstance, newState v1alpha1.ServiceInstanceState) error
}

// //go:generate counterfeiter . BrokerClient
//
// type BrokerClient interface {
// 	osbapi.Client
// }

type ServiceInstanceReconciler struct {
	createBrokerClient      osbapi.CreateFunc
	kubeServiceInstanceRepo KubeServiceInstanceRepo
	kubeBrokerRepo          KubeBrokerRepo
}

func NewServiceInstanceReconciler(
	createBrokerClient osbapi.CreateFunc,
	kubeServiceInstanceRepo KubeServiceInstanceRepo,
	kubeBrokerRepo KubeBrokerRepo,
) *ServiceInstanceReconciler {
	return &ServiceInstanceReconciler{
		createBrokerClient:      createBrokerClient,
		kubeServiceInstanceRepo: kubeServiceInstanceRepo,
		kubeBrokerRepo:          kubeBrokerRepo,
	}
}

func (r *ServiceInstanceReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	instance, err := r.kubeServiceInstanceRepo.Get(request.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	broker, err := r.kubeBrokerRepo.Get(types.NamespacedName{Name: instance.Spec.BrokerName, Namespace: instance.ObjectMeta.Namespace})
	if err != nil {
		//TODO what if the broker does not exist?
		return reconcile.Result{}, err
	}

	if instance.Status.State == v1alpha1.ServiceInstanceStateProvisioned {
		return reconcile.Result{}, nil
	}

	osbapiConfig := brokerClientConfig(broker)

	osbapiClient, err := r.createBrokerClient(osbapiConfig)
	if err != nil {
		return reconcile.Result{}, err
	}

	_, err = osbapiClient.ProvisionInstance(&osbapi.ProvisionRequest{
		InstanceID:        string(instance.ObjectMeta.UID),
		AcceptsIncomplete: false,
		ServiceID:         instance.Spec.ServiceID,
		PlanID:            instance.Spec.PlanID,
		OrganizationGUID:  instance.ObjectMeta.Namespace,
		SpaceGUID:         instance.ObjectMeta.Namespace,
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.kubeServiceInstanceRepo.UpdateState(instance, v1alpha1.ServiceInstanceStateProvisioned); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
