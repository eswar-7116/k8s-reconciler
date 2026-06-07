package apiserver

import (
	"sync"

	"github.com/eswar-7116/k8s-reconciler/internal/cluster"
	"github.com/eswar-7116/k8s-reconciler/internal/reconciler"
)

type APIServer struct {
	cluster    *cluster.Cluster
	workqueue  chan struct{}
	wg         *sync.WaitGroup
	reconciler *reconciler.Reconciler
}

func NewAPIServer(replicas int) *APIServer {
	c := cluster.NewCluster(replicas)
	wq := make(chan struct{}, 10)
	wg := &sync.WaitGroup{}

	r := reconciler.NewReconciler(c, wq, wg)
	r.Start()

	a := &APIServer{
		cluster:    c,
		workqueue:  wq,
		reconciler: r,
		wg:         wg,
	}
	a.notify()

	return a
}

func (a *APIServer) CreatePod() {
	a.cluster.CreatePod()
	a.notify()
}

func (a *APIServer) TerminatePod() {
	a.cluster.TerminatePod()
	a.notify()
}

func (a *APIServer) SetReplicas(replicas int) {
	a.cluster.SetReplicas(replicas)
	a.notify()
}

func (a *APIServer) State() (int, int) {
	return a.cluster.State()
}

func (a *APIServer) notify() {
	select {
	case a.workqueue <- struct{}{}:
	default:
	}
}

func (a *APIServer) Close() {
	close(a.workqueue)
	a.wg.Wait()
}
