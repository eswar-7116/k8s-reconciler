package reconciler

import (
	"sync"

	"github.com/eswar-7116/k8s-reconciler/internal/cluster"
)

type Reconciler struct {
	cluster   *cluster.Cluster
	workqueue <-chan struct{}
	wg        *sync.WaitGroup
}

func NewReconciler(cluster *cluster.Cluster, workqueue <-chan struct{}, wg *sync.WaitGroup) *Reconciler {
	return &Reconciler{
		cluster:   cluster,
		workqueue: workqueue,
		wg:        wg,
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
	r.wg.Go(func() {
		for range r.workqueue {
			r.reconcile()
		}
	})
}
