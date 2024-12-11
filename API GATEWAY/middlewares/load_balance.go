package middlewares

import (
	"net/http"
	"sync/atomic"
	"time"
)

// LoadBalancer manages backend server selection
type LoadBalancer struct {
	servers []string
	counter uint64
	healthy []int32
}

// NewLoadBalancer initializes a new LoadBalancer
func NewLoadBalancer(servers []string) *LoadBalancer {
	healthy := make([]int32, len(servers))
	for i := range healthy {
		healthy[i] = 1 // Assume all servers are healthy initially
	}
	return &LoadBalancer{
		servers: servers,
		healthy: healthy,
	}
}

// RoundRobinLoadBalancer selects the next healthy server using round-robin
func (lb *LoadBalancer) RoundRobinLoadBalancer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalServers := len(lb.servers)
		for i := 0; i < totalServers; i++ {
			idx := atomic.AddUint64(&lb.counter, 1)
			serverIndex := int(idx % uint64(totalServers))
			if atomic.LoadInt32(&lb.healthy[serverIndex]) == 1 {
				target := lb.servers[serverIndex]
				// Modify the request to point to the selected backend
				r.URL.Host = target
				r.URL.Scheme = "http"
				next.ServeHTTP(w, r)
				return
			}
		}
		// If no healthy servers found
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	})
}

// MarkServerHealthy sets the health status of a server
func (lb *LoadBalancer) MarkServerHealthy(index int, healthy bool) {
	if index >= 0 && index < len(lb.healthy) {
		if healthy {
			atomic.StoreInt32(&lb.healthy[index], 1)
		} else {
			atomic.StoreInt32(&lb.healthy[index], 0)
		}
	}
}

// HealthCheck periodically checks the health of backend servers
func (lb *LoadBalancer) HealthCheck(interval time.Duration) {
	for {
		for i, server := range lb.servers {
			go func(index int, url string) {
				client := http.Client{
					Timeout: 5 * time.Second,
				}
				resp, err := client.Get(url + "/health")
				if err != nil || resp.StatusCode != http.StatusOK {
					lb.MarkServerHealthy(index, false)
				} else {
					lb.MarkServerHealthy(index, true)
				}
			}(i, server)
		}
		time.Sleep(interval)
	}
}
