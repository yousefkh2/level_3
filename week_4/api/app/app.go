package app

import (
	"context"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DatabaseService interface {
	CreateDatabase(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error)
	ListDatabases(ctx context.Context) ([]cnpgv1.Cluster, error)
	GetDatabase(ctx context.Context, name string) (*cnpgv1.Cluster, error)
	DeleteDatabase(ctx context.Context, name string) error
	UpdateDatabase(ctx context.Context, name string, instances *int, storage *string) (*cnpgv1.Cluster, error)
}

type App struct {
	K8sClient client.Client
	DBService DatabaseService
	Logger    *zap.Logger
	Metrics   *Metrics
}
