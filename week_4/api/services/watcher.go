package services

import (
	"context"
	"time"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StatusWatcher struct {
	k8sClient client.Client
	logger    *zap.Logger
	namespace string

	lastKnownStatus map[string]string
}

func NewStatusWatcher(k8sClient client.Client, logger *zap.Logger, namespace string) *StatusWatcher {
	return &StatusWatcher{
		k8sClient:       k8sClient,
		logger:          logger,
		namespace:       namespace,
		lastKnownStatus: make(map[string]string),
	}
}

// Start begins the watch loop. Call this in a goroutine at startup.
// It polls every 15 seconds — a pragmatic alternative to a true Watch stream
// which requires additional reconnect/resync logic.
//
// Why polling instead of Watch API?
// The controller-runtime Watch requires setting up a full manager/cache.
// For a lightweight status logger, polling every 15s is simpler and sufficient
// since status changes (provisioning, scaling) take minutes, not milliseconds.
func (w *StatusWatcher) Start(ctx context.Context) {
	w.logger.Info("status watcher started", zap.String("namespace", w.namespace))

	// Do an initial sync immediately, then tick every 15 seconds
	w.reconcile(ctx)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("status watcher stopped")
			return
		case <-ticker.C:
			w.reconcile(ctx)
		}
	}
}

func (w *StatusWatcher) reconcile(ctx context.Context) {
	var clusterList cnpgv1.ClusterList
	if err := w.k8sClient.List(ctx, &clusterList, client.InNamespace(w.namespace)); err != nil {
		w.logger.Error("status watcher: failed to list clusters", zap.Error(err))
		return
	}

	// Build a set of currently existing database names
	// so we can detect deletions
	existing := make(map[string]bool)

	for _, cluster := range clusterList.Items {
		name := cluster.Name
		phase := statusPhase(&cluster)
		existing[name] = true

		previous, seen := w.lastKnownStatus[name]

		if !seen {
			// First time we see this database — emit an initial status event so
			// status logs are useful immediately after startup/redeploy.
			w.logger.Info("service event",
				zap.String("log_type", "service"),
				zap.String("resource", name),
				zap.String("current_status", phase),
				zap.String("event", "initial status observed: "+phase),
			)
			w.lastKnownStatus[name] = phase
			continue
		}

		if phase != previous {
			w.logger.Info("service event",
				zap.String("log_type", "service"),
				zap.String("resource", name),
				zap.String("previous_status", previous),
				zap.String("current_status", phase),
				zap.String("event", statusEventMessage(previous, phase)),
			)
			w.lastKnownStatus[name] = phase
		}
	}

	for name := range w.lastKnownStatus {
		if !existing[name] {
			w.logger.Info("service event",
				zap.String("log_type", "service"),
				zap.String("resource", name),
				zap.String("event", "database deleted from cluster"),
			)
			delete(w.lastKnownStatus, name)
		}
	}
}

func statusPhase(cluster *cnpgv1.Cluster) string {
	phase := string(cluster.Status.Phase)
	if phase == "" {
		return "unknown"
	}
	return phase
}

func statusEventMessage(previous, current string) string {
	_ = previous // available for more specific messages in the future
	switch current {
	case "Cluster in healthy state":
		return "database is ready"
	case "Setting up primary":
		return "database is provisioning"
	case "Cluster is not ready":
		return "database entered degraded state"
	default:
		return "database status changed to: " + current
	}
}

var _ runtime.Object = &cnpgv1.ClusterList{}
