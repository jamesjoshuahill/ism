package repositories

import (
	"context"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type KubeServiceRepo struct {
	client client.Client
	scheme *runtime.Scheme
}

func NewKubeServiceRepo(client client.Client) *KubeServiceRepo {
	return &KubeServiceRepo{
		client: client,
		scheme: scheme.Scheme,
	}
}

func (repo *KubeServiceRepo) Create(broker *v1alpha1.Broker, catalogService osbapi.Service) (*v1alpha1.BrokerService, error) {
	service := &v1alpha1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      catalogService.ID,
			Namespace: broker.Namespace,
		},
		Spec: v1alpha1.BrokerServiceSpec{
			BrokerID:    broker.Name,
			Name:        catalogService.Name,
			Description: catalogService.Description,
		},
	}

	if err := controllerutil.SetControllerReference(broker, service, repo.scheme); err != nil {
		return nil, err
	}

	if err := repo.client.Create(context.TODO(), service); err != nil {
		return nil, err
	}

	return service, nil
}
