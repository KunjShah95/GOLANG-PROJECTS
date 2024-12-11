package middlewares

import (
	"net/http"
	"time"
)

// HealthCheckMiddleware verifies the health of backend servers
func HealthCheckMiddleware(servers []string, interval time.Duration) {
	// Implement health check logic, possibly updating the LoadBalancer's server list
	// This is a placeholder for actual implementation
	go func() {
		for {
			for _, server := range servers {
				go func(s string) {
					resp, err := http.Get(s + "/health")
					if err != nil || resp.StatusCode != http.StatusOK {
						// Mark server as unhealthy in LoadBalancer
					} else {
						// Mark server as healthy in LoadBalancer
					}
				}(server)
			}
			time.Sleep(interval)
		}
	}()
}
