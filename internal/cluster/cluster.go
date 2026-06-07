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
	defer c.mu.Unlock()
	c.pods = append(c.pods, Pod{
		Id: c.nextId,
	})
	fmt.Println("Created pod with id", c.nextId)
	c.nextId++
}

func (c *Cluster) TerminatePod() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.pods) == 0 {
		return
	}

	id := c.pods[len(c.pods)-1].Id
	c.pods = c.pods[:len(c.pods)-1]
	fmt.Println("Terminated pod with id", id)
}

func (c *Cluster) SetReplicas(replicas int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deployment.Replicas = replicas
}

func (c *Cluster) State() (pods, replicas int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	pods = len(c.pods)
	replicas = c.deployment.Replicas
	return
}
