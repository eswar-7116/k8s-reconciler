package apiserver

import (
	"sync"
	"testing"
)

func TestInitialState(t *testing.T) {
	api := NewAPIServer(3)
	api.Close()

	pods, replicas := api.cluster.State()
	if pods != replicas {
		t.Errorf("expected %d pods, got %d", replicas, pods)
	}
}

func TestState(t *testing.T) {
	api := NewAPIServer(4)
	api.SetReplicas(6)
	api.Close()

	pods, replicas := api.State()
	if replicas != 6 {
		t.Errorf("expected desired replicas 6, got %d", replicas)
	}
	if pods != 6 {
		t.Errorf("expected 6 running pods, got %d", pods)
	}
}

func TestScaleUp(t *testing.T) {
	api := NewAPIServer(2)
	api.SetReplicas(5)
	api.Close()

	pods, _ := api.cluster.State()
	if pods != 5 {
		t.Errorf("expected 5 pods, got %d", pods)
	}
}

func TestScaleDown(t *testing.T) {
	api := NewAPIServer(5)
	api.SetReplicas(2)
	api.Close()

	pods, _ := api.cluster.State()
	if pods != 2 {
		t.Errorf("expected 2 pods, got %d", pods)
	}
}

func TestScaleToZero(t *testing.T) {
	api := NewAPIServer(4)
	api.SetReplicas(0)
	api.Close()

	pods, _ := api.cluster.State()
	if pods != 0 {
		t.Errorf("expected 0 pods, got %d", pods)
	}
}

func TestRapidReplicaChanges(t *testing.T) {
	api := NewAPIServer(1)
	api.SetReplicas(10)
	api.SetReplicas(3)
	api.SetReplicas(7)
	api.SetReplicas(0)
	api.SetReplicas(5)
	api.Close()

	pods, _ := api.cluster.State()
	if pods != 5 {
		t.Errorf("expected 5 pods, got %d", pods)
	}
}

func TestConcurrentSetReplicas(t *testing.T) {
	api := NewAPIServer(0)
	var wg sync.WaitGroup

	for range 20 {
		wg.Go(func() {
			api.SetReplicas(5)
		})
	}
	wg.Wait()
	api.Close()

	pods, replicas := api.cluster.State()
	if pods != replicas {
		t.Errorf("expected %d pods, got %d", replicas, pods)
	}
}
