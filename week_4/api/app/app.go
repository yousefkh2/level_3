package app

import (
	"github.com/yousefkh2/level_3/week_4/api/services"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type App struct {
	K8sClient client.Client // controller-runtime client (not Clientset)
	DBService *services.DatabaseService
	// Later: config, etc.
}
