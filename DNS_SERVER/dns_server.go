package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Server Configuration
var (
	upstreamServers  = []string{"8.8.8.8:53", "1.1.1.1:53"}
	currentUpstream  = 0
	blocklist        = map[string]bool{"ads.example.com": true, "malware.net": true}
	allowedIPs       = map[string]bool{"127.0.0.1": true, "::1": true}
	queryLimit       = 10
	queryCount       = make(map[string]int)
	blackholeAddress = "0.0.0.0"
	logFile          = "dns_queries.log"

	// Metrics
	totalQueries      int
	blockedQueries    int
	perDomainCounters = make(map[string]int)

	mu sync.Mutex
)

// DNSHeader represents a DNS message header
type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// DNSQuestion represents a DNS question section
type DNSQuestion struct {
	Name  string
	Type  uint16
	Class uint16
}

// Round-robin load balancer for upstream DNS servers
func getNextUpstream() string {
	mu.Lock()
	defer mu.Unlock()
	currentUpstream = (currentUpstream + 1) % len(upstreamServers)
	return upstreamServers[currentUpstream]
}

// Check if IP is allowed (Access Control)
func isAllowedIP(ip string) bool {
	_, allowed := allowedIPs[ip]
	return allowed
}

// Check rate limit
func isRateLimited(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().Minute()
	if queryCount[ip] != now {
		queryCount[ip] = 0
	}
	queryCount[ip]++
	return queryCount[ip] > queryLimit
}

// Check if domain is blocked
func isBlocked(domain string) bool {
	_, exists := blocklist[domain]
	return exists
}

// Encode domain name to DNS format
func encodeDomainName(domain string) []byte {
	parts := strings.Split(domain, ".")
	var encoded []byte
	for _, part := range parts {
		encoded = append(encoded, byte(len(part)))
		encoded = append(encoded, []byte(part)...)
	}
	encoded = append(encoded, 0)
	return encoded
}

// Log DNS queries
func logQuery(clientIP, domain string) {
	mu.Lock()
	defer mu.Unlock()
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()
	entry := fmt.Sprintf("%s - %s queried %s\n", time.Now().Format(time.RFC3339), clientIP, domain)
	file.WriteString(entry)
}

// Increment metrics
func incrementMetrics(domain string, blocked bool) {
	mu.Lock()
	defer mu.Unlock()
	totalQueries++
	if blocked {
		blockedQueries++
	}
	perDomainCounters[domain]++
}

// Display metrics
func printMetrics() {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("---- DNS Server Metrics ----")
	fmt.Printf("Total Queries: %d\n", totalQueries)
	fmt.Printf("Blocked Queries: %d\n", blockedQueries)
	fmt.Println("Per-Domain Query Counts:")
	for domain, count := range perDomainCounters {
		fmt.Printf("  %s: %d\n", domain, count)
	}
	fmt.Println("----------------------------")
}

// Build blackhole response (for blocked domains)
func buildBlackholeResponse(header DNSHeader, question DNSQuestion) []byte {
	response := make([]byte, 12)
	binary.BigEndian.PutUint16(response[0:2], header.ID)
	binary.BigEndian.PutUint16(response[2:4], 0x8180) // Flags: Standard response, No error
	binary.BigEndian.PutUint16(response[4:6], 1)      // Questions: 1
	binary.BigEndian.PutUint16(response[6:8], 1)      // Answers: 1
	response = append(response, encodeDomainName(question.Name)...)
	response = append(response, 0x00, 0x01, 0x00, 0x01) // Type A, Class IN
	response = append(response, encodeDomainName(question.Name)...)
	response = append(response, 0x00, 0x01, 0x00, 0x01)
	response = append(response, 0x00, 0x00, 0x00, 0x3C)                 // TTL 60
	response = append(response, 0x00, 0x04)                             // Data length
	response = append(response, net.ParseIP(blackholeAddress).To4()...) // IP 0.0.0.0
	return response
}

// Forward query to upstream server
func forwardToUpstream(upstream string, request []byte) ([]byte, error) {
	conn, err := net.Dial("udp", upstream)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(request)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 512)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}
	return response[:n], nil
}

// Parse domain name from DNS query
func parseDomainName(data []byte) (string, int) {
	var parts []string
	i := 0
	for {
		length := int(data[i])
		if length == 0 {
			break
		}
		i++
		parts = append(parts, string(data[i:i+length]))
		i += length
	}
	return strings.Join(parts, "."), i + 1
}

// Handle incoming DNS requests
func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, request []byte) {
	clientIP := addr.IP.String()
	header := DNSHeader{
		ID: binary.BigEndian.Uint16(request[:2]),
	}

	domain, _ := parseDomainName(request[12:])
	log.Printf("Received query for %s from %s", domain, clientIP)

	// Access control
	if !isAllowedIP(clientIP) {
		log.Printf("Access denied for %s", clientIP)
		return
	}

	// Rate limiting
	if isRateLimited(clientIP) {
		log.Printf("Rate limit exceeded for %s", clientIP)
		return
	}

	// Blocklist check
	if isBlocked(domain) {
		log.Printf("Blocked domain: %s", domain)
		incrementMetrics(domain, true)
		response := buildBlackholeResponse(header, DNSQuestion{Name: domain})
		conn.WriteToUDP(response, addr)
		return
	}

	// Forward to upstream server
	upstream := getNextUpstream()
	response, err := forwardToUpstream(upstream, request)
	if err != nil {
		log.Printf("Failed to forward query: %v", err)
		return
	}

	// Log query and increment metrics
	logQuery(clientIP, domain)
	incrementMetrics(domain, false)

	// Send response
	conn.WriteToUDP(response, addr)
}

func main() {
	addr := net.UDPAddr{Port: 53, IP: net.ParseIP("0.0.0.0")}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}
	defer conn.Close()
	log.Println("DNS server is running on 0.0.0.0:53")

	// Print metrics every 30 seconds
	go func() {
		for {
			time.Sleep(30 * time.Second)
			printMetrics()
		}
	}()

	for {
		buffer := make([]byte, 512)
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}
		go handleRequest(conn, clientAddr, buffer[:n])
	}
}
