//business logic, separate from handlers to keep things clean

package services

import (
	"context"                                                // timer for the request. if the API takes too long, ctx tells it to give up
	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1" //this is cnpg specific stuff
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DatabaseService struct {
	K8sClient client.Client // controller-runtime client (instead of kubectl)
}

func NewDatabaseService(k8sClient client.Client) *DatabaseService {
	return &DatabaseService{K8sClient: k8sClient}
}

func (s *DatabaseService) CreateDatabase(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error) {
	cluster := &cnpgv1.Cluster{ // now we initialize the object... this is the empty form we are filling
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: cnpgv1.ClusterSpec{
			Instances: instances,
			StorageConfiguration: cnpgv1.StorageConfiguration{
				Size: storage,
			},
		},
	}

	//  this is the magic line - creates it in kubernetes
	err := s.K8sClient.Create(ctx, cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (s *DatabaseService) ListDatabases(ctx context.Context) ([]cnpgv1.Cluster, error) {
	var clusterList cnpgv1.ClusterList

	err := s.K8sClient.List(ctx, &clusterList, client.InNamespace("default"))
	if err != nil {
		return nil, err
	}

	return clusterList.Items, nil
}

func (s *DatabaseService) GetDatabase(ctx context.Context, name string) (*cnpgv1.Cluster, error) {
	var cluster cnpgv1.Cluster

	err := s.K8sClient.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: "default",
	}, &cluster)

	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (s *DatabaseService) DeleteDatabase(ctx context.Context, name string) error {
	cluster := &cnpgv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
	}

	return s.K8sClient.Delete(ctx, cluster)
}
