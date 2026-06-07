package cluster

import (
	"sync"
	"testing"
)

func TestCreatePod(t *testing.T) {
	c := NewCluster(0)

	c.CreatePod()
	c.CreatePod()

	pods, _ := c.State()
	if pods != 2 {
		t.Errorf("expected 2 pods, got %d", pods)
	}
}

func TestTerminatePod(t *testing.T) {
	c := NewCluster(0)
	c.CreatePod()
	c.CreatePod()
	c.CreatePod()

	c.TerminatePod()

	pods, _ := c.State()
	if pods != 2 {
		t.Errorf("expected 2 pods, got %d", pods)
	}
}

func TestTerminatePodOnEmpty(t *testing.T) {
	c := NewCluster(0)

	c.TerminatePod()

	pods, _ := c.State()
	if pods != 0 {
		t.Errorf("expected 0 pods, got %d", pods)
	}
}

func TestSetReplicas(t *testing.T) {
	c := NewCluster(3)

	_, replicas := c.State()
	if replicas != 3 {
		t.Errorf("expected 3 replicas, got %d", replicas)
	}

	c.SetReplicas(7)

	_, replicas = c.State()
	if replicas != 7 {
		t.Errorf("expected 7 replicas, got %d", replicas)
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := NewCluster(0)
	var wg sync.WaitGroup

	// Concurrent creates
	for range 50 {
		wg.Go(func() {
			c.CreatePod()
		})
	}
	wg.Wait()

	pods, _ := c.State()
	if pods != 50 {
		t.Errorf("expected 50 pods after concurrent creates, got %d", pods)
	}

	// Concurrent terminates
	for range 50 {
		wg.Go(func() {
			c.TerminatePod()
		})
	}
	wg.Wait()

	pods, _ = c.State()
	if pods != 0 {
		t.Errorf("expected 0 pods after concurrent terminates, got %d", pods)
	}

	// Concurrent mixed operations
	for range 30 {
		wg.Add(3)
		go func() {
			defer wg.Done()
			c.CreatePod()
		}()
		go func() {
			defer wg.Done()
			c.SetReplicas(10)
		}()
		go func() {
			defer wg.Done()
			c.State()
		}()
	}
	wg.Wait()
}
