package repositories

import (
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// var ctx = context.TODO()

type KubeServiceInstanceRepo struct {
	client client.Client
}

func NewKubeServiceInstanceRepo(client client.Client) *KubeServiceInstanceRepo {
	return &KubeServiceInstanceRepo{
		client: client,
	}
}

func (repo *KubeServiceInstanceRepo) Get(resource types.NamespacedName) (*v1alpha1.ServiceInstance, error) {
	serviceInstance := &v1alpha1.ServiceInstance{}

	err := repo.client.Get(ctx, resource, serviceInstance)
	if err != nil {
		return nil, err
	}

	return serviceInstance, nil
}

func (repo *KubeServiceInstanceRepo) UpdateState(serviceInstance *v1alpha1.ServiceInstance, newState v1alpha1.ServiceInstanceState) error {
	serviceInstance.Status.State = newState

	return repo.client.Status().Update(ctx, serviceInstance)
}
