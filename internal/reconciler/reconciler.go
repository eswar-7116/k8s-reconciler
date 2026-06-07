package reconciler

import "github.com/eswar-7116/k8s-reconciler/internal/cluster"

type Reconciler struct {
	cluster   *cluster.Cluster
	workqueue <-chan struct{}
}

func NewReconciler(cluster *cluster.Cluster, workqueue <-chan struct{}) *Reconciler {
	return &Reconciler{
		cluster:   cluster,
		workqueue: workqueue,
	}
}

func (r *Reconciler) reconcile() {
	for {
		pods, replicas := r.cluster.State()

		if pods == replicas {
			return
		}

		if replicas > pods {
			r.cluster.CreatePod()
		} else {
			r.cluster.TerminatePod()
		}
	}
}

func (r *Reconciler) Start() {
	go func() {
		for range r.workqueue {
			r.reconcile()
		}
	}()
}
