package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "sync"
    "sync/atomic"
    "time"
)

type Server struct {
    URL            string
    Healthy        bool
    LastChecked    time.Time
    CircuitBreaker bool
}

type LoadBalancer struct {
    servers        []*Server
    index          uint64
    mu             sync.RWMutex
    shutdown       chan struct{}
}

func NewLoadBalancer(servers []string) *LoadBalancer {
    lbServers := make([]*Server, len(servers))
    for i, server := range servers {
        lbServers[i] = &Server{URL: server, Healthy: true, CircuitBreaker: false}
    }
    return &LoadBalancer{
        servers:  lbServers,
        shutdown: make(chan struct{}),
    }
}

func (lb *LoadBalancer) getNextServer() *Server {
    lb.mu.RLock()
    defer lb.mu.RUnlock()
    nextIndex := atomic.AddUint64(&lb.index, 1)
    return lb.servers[(int(nextIndex)-1)%len(lb.servers)]
}

func (lb *LoadBalancer) healthCheck(server *Server) {
    resp, err := http.Get(server.URL + "/health")
    if err != nil || resp.StatusCode != http.StatusOK {
        server.Healthy = false
        server.CircuitBreaker = true
    } else {
        server.Healthy = true
        server.CircuitBreaker = false
    }
    server.LastChecked = time.Now()
    if resp != nil {
        resp.Body.Close()
    }
}

func (lb *LoadBalancer) handleRequest(w http.ResponseWriter, r *http.Request) {
    var server *Server
    for i := 0; i < len(lb.servers); i++ {
        server = lb.getNextServer()
        if server.Healthy && !server.CircuitBreaker {
            break
        }
    }

    if server.CircuitBreaker {
        http.Error(w, "Server is temporarily unavailable", http.StatusServiceUnavailable)
        return
    }

    if !server.Healthy {
        http.Error(w, "No healthy backend servers available", http.StatusServiceUnavailable)
        return
    }

    proxyURL := fmt.Sprintf("%s%s", server.URL, r.URL.Path)
    client := &http.Client{
        Timeout: 5 * time.Second, //Set a timeout for the request
    }
    resp, err := client.Get(proxyURL)
    if err != nil {
        http.Error(w, "Error contacting backend server", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

func (lb *LoadBalancer) startHealthChecks(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            lb.mu.Lock()
            for _, server := range lb.servers {
                go lb.healthCheck(server) // Run health checks concurrently
            }
            lb.mu.Unlock()
        case <-lb.shutdown:
            return
        }
    }
}

func (lb *LoadBalancer) addServer(url string) {
    lb.mu.Lock()
    defer lb.mu.Unlock()
    lb.servers = append(lb.servers, &Server{URL: url, Healthy: true, CircuitBreaker: false})
}

func (lb *LoadBalancer) removeServer(url string) {
    lb.mu.Lock()
    defer lb.mu.Unlock()
    for i, server := range lb.servers {
        if server.URL == url {
            lb.servers = append(lb.servers[:i], lb.servers[i+1:]...)
            break
        }
    }
}

func (lb *LoadBalancer) gracefulShutdown() {
    log.Println("Shutting down gracefully...")
    close(lb.shutdown)
}

func main() {
    servers := []string{"http://localhost:8081", "http://localhost:8082", "http://localhost:8083"} // List of backend servers
 lb := NewLoadBalancer(servers)

    go lb.startHealthChecks(10 * time.Second) // Start health checks every 10 seconds

    http.HandleFunc("/", lb.handleRequest)

    // Example endpoints to add/remove servers
    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
        url := r.URL.Query().Get("url")
        if url != "" {
            lb.addServer(url)
            fmt.Fprintf(w, "Added server: %s\n", url)
        } else {
            http.Error(w, "URL parameter is required", http.StatusBadRequest)
        }
    })

    http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
        url := r.URL.Query().Get("url")
        if url != "" {
            lb.removeServer(url)
            fmt.Fprintf(w, "Removed server: %s\n", url)
        } else {
            http.Error(w, "URL parameter is required", http.StatusBadRequest)
        }
    })

    log.Println("Load Balancer is running on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Error starting server:", err)
    }

    defer lb.gracefulShutdown() // Ensure graceful shutdown on exit
}