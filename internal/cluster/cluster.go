package cluster

import (
	"fmt"
	"sync"
)

type Cluster struct {
	mu         sync.RWMutex
	deployment Deployment
	pods       []Pod
	nextId     int
}

func NewCluster(replicas int) *Cluster {
	return &Cluster{
		deployment: Deployment{
			Replicas: replicas,
		},
		nextId: 1,
	}
}

func (c *Cluster) CreatePod() {
	c.mu.Lock()
	c.pods = append(c.pods, Pod{
		Id: c.nextId,
	})
	fmt.Println("Created pod with id", c.nextId)
	c.nextId++
	c.mu.Unlock()
}

func (c *Cluster) TerminatePod() {
	c.mu.Lock()
	if len(c.pods) == 0 {
		return
	}

	id := c.pods[len(c.pods)-1].Id
	c.pods = c.pods[:len(c.pods)-1]
	fmt.Println("Terminated pod with id", id)
	c.mu.Unlock()
}

func (c *Cluster) SetReplicas(replicas int) {
	c.mu.Lock()
	c.deployment.Replicas = replicas
	c.mu.Unlock()
}

func (c *Cluster) State() (pods, replicas int) {
	c.mu.RLock()
	pods = len(c.pods)
	replicas = c.deployment.Replicas
	c.mu.RUnlock()
	return
}
