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
	pods, replicas := r.cluster.State()
	diff := replicas - pods

	if diff > 0 {
		for range diff {
			r.cluster.CreatePod()
		}
	} else {
		for range -diff {
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
