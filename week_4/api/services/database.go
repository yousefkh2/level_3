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
