//business logic, separate from handlers to keep things clean

package services

import (
	"context"                                                // timer for the request. if the API takes too long, ctx tells it to give up
	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1" //this is cnpg specific stuff
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DatabaseService struct {
	K8sClient client.Client // controller-runtime client (instead of clientSet)
}

const DatabaseNamespace = "default"

func NewDatabaseService(k8sClient client.Client) *DatabaseService {
	return &DatabaseService{K8sClient: k8sClient}
}

func (s *DatabaseService) CreateDatabase(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error) {
	cluster := &cnpgv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: DatabaseNamespace,
		},
		Spec: cnpgv1.ClusterSpec{
			Instances: instances,
			StorageConfiguration: cnpgv1.StorageConfiguration{
				Size: storage,
			},
		},
	}

	err := s.K8sClient.Create(ctx, cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (s *DatabaseService) ListDatabases(ctx context.Context) ([]cnpgv1.Cluster, error) {
	var clusterList cnpgv1.ClusterList

	err := s.K8sClient.List(ctx, &clusterList, client.InNamespace(DatabaseNamespace))
	if err != nil {
		return nil, err
	}

	return clusterList.Items, nil
}

func (s *DatabaseService) GetDatabase(ctx context.Context, name string) (*cnpgv1.Cluster, error) {
	var cluster cnpgv1.Cluster

	err := s.K8sClient.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: DatabaseNamespace,
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
			Namespace: DatabaseNamespace,
		},
	}

	return s.K8sClient.Delete(ctx, cluster)
}

func (s *DatabaseService) UpdateDatabase(ctx context.Context, name string, instances *int, storage *string) (*cnpgv1.Cluster, error) {
	var cluster cnpgv1.Cluster
	err := s.K8sClient.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: DatabaseNamespace,
	}, &cluster)
	if err != nil {
		return nil, err
	}

	// create a patch base (snapshot of current state)
	patch := client.MergeFrom(cluster.DeepCopy())

	// mutate only the fields that were provided
	if instances != nil {
		cluster.Spec.Instances = *instances
	}
	if storage != nil {
		cluster.Spec.StorageConfiguration.Size = *storage
	}

	// apply the patch (sends only the diff to k8s)
	err = s.K8sClient.Patch(ctx, &cluster, patch) // patch (old); cluster (new)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}
