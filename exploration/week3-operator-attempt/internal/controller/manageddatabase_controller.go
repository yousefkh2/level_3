/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	paasv1alpha1 "github.com/yousefkh2/database-operator/api/v1alpha1"
)

// ManagedDatabaseReconciler reconciles a ManagedDatabase object
type ManagedDatabaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=paas.paas.level3.cloud,resources=manageddatabases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=paas.paas.level3.cloud,resources=manageddatabases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=paas.paas.level3.cloud,resources=manageddatabases/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ManagedDatabase object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile
func (r *ManagedDatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// TODO(user): your logic here
	// 1) Load ManagedDatabase
	var mdb paasv1alpha1.ManagedDatabase
	err := r.Get(ctx, req.NamespacedName, &mdb)

    // 2) Compute desired Cluster spec from tier
	clusterSpec, err := r.computeClusterSpec(&mdb)
    if err != nil {
        log.Error(err, "invalid tier configuration")
        mdb.Status.Phase = "Failed"
        mdb.Status.Message = err.Error()
        r.Status().Update(ctx, &mdb)
        return ctrl.Result{}, err
    }
    // 3) Ensure Cluster exists (create/patch)
	
    // 4) Update status (phase/service/secret)
    // 5) Requeue or finish

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ManagedDatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&paasv1alpha1.ManagedDatabase{}).
		Named("manageddatabase").
		Complete(r)
}


func (r *ManagedDatabaseReconciler) computeClusterSpec(db *paasv1alpha1.ManagedDatabase) (ClusterSpec, error) {
    // Map tier to resources
    tiers := map[string]struct{
        Instances int
        CPU       string
        Memory    string
        Storage   string
    }{
        "small":  {1, "500m", "1Gi", "10Gi"},
        "medium": {3, "1", "2Gi", "50Gi"},
        "large":  {3, "2", "4Gi", "100Gi"},
    }
    
    config, ok := tiers[db.Spec.Tier]
    if !ok {
        return ClusterSpec{}, fmt.Errorf("unknown tier: %s", db.Spec.Tier)
    }
    
    // Build and return cluster spec
    // ...
}