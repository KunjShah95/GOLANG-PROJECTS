package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	IP      string `json:"ip" xml:"ip"`
	Port    int    `json:"port" xml:"port"`
	Status  string `json:"status" xml:"status"`
	Service string `json:"service" xml:"service"`
	Banner  string `json:"banner,omitempty" xml:"banner,omitempty"`
}

// ParsePorts parses the port input (e.g., "20-80,8080").
func ParsePorts(portInput string) ([]int, error) {
	ports := []int{}
	ranges := strings.Split(portInput, ",")
	for _, r := range ranges {
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			start, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid port range: %s", r)
			}
			end, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid port range: %s", r)
			}
			for i := start; i <= end; i++ {
				ports = append(ports, i)
			}
		} else {
			port, err := strconv.Atoi(r)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", r)
			}
			ports = append(ports, port)
		}
	}
	return ports, nil
}

// PingHost checks if a host is alive using ICMP ping.
func PingHost(host string) bool {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", host)
	err := cmd.Run()
	return err == nil
}

// BannerGrab grabs a service banner from the given address.
func BannerGrab(address string, timeout time.Duration) string {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return ""
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	return strings.TrimSpace(string(buffer[:n]))
}

// ServiceDetection maps common ports to services.
func ServiceDetection(port int) string {
	services := map[int]string{
		20: "FTP", 21: "FTP", 22: "SSH", 23: "Telnet", 25: "SMTP",
		53: "DNS", 80: "HTTP", 110: "POP3", 443: "HTTPS", 3306: "MySQL",
	}
	if service, exists := services[port]; exists {
		return service
	}
	return "Unknown"
}

// ScanPort scans a single port, performs banner grabbing, and detects service versions.
func ScanPort(ip string, port int, results chan ScanResult, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	result := ScanResult{IP: ip, Port: port, Service: ServiceDetection(port)}

	if err != nil {
		result.Status = "Closed"
		results <- result
		return
	}
	defer conn.Close()

	result.Status = "Open"
	result.Banner = BannerGrab(address, timeout)
	results <- result
}

// SaveResults saves scan results to a specified file format (JSON, XML, HTML).
func SaveResults(filename string, format string, results []ScanResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create results file: %v", err)
	}
	defer file.Close()

	switch format {
	case "json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		return encoder.Encode(results)
	case "xml":
		encoder := xml.NewEncoder(file)
		encoder.Indent("", "  ")
		return encoder.Encode(results)
	case "html":
		writer := bufio.NewWriter(file)
		writer.WriteString("<html><body><h1>Scan Results</h1><table border='1'><tr><th>IP</th><th>Port</th><th>Status</th><th>Service</th><th>Banner</th></tr>")
		for _, result := range results {
			writer.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>", result.IP, result.Port, result.Status, result.Service, result.Banner))
		}
		writer.WriteString("</table></body></html>")
		writer.Flush()
	}
	return nil
}

// expandCIDR expands a CIDR notation into a list of IP addresses.
func expandCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// Remove network address and broadcast address
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}
	return ips, nil
}

// inc increments an IP address.
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// main function
func main() {
	fmt.Print("Enter target IP or CIDR (e.g., 192.168.1.1 or 192.168.1.0/24): ")
	var ipInput string
	fmt.Scanln(&ipInput)

	fmt.Print("Enter ports to scan (e.g., 20-80,443,8080): ")
	var portInput string
	fmt.Scanln(&portInput)

	fmt.Print("Enter output file name (without extension): ")
	var outputFile string
	fmt.Scanln(&outputFile)

	fmt.Print("Choose output format (json/xml/html): ")
	var outputFormat string
	fmt.Scanln(&outputFormat)

	ports, err := ParsePorts(portInput)
	if err != nil {
		fmt.Printf("Error parsing ports: %v\n", err)
		return
	}

	var ips []string
	if strings.Contains(ipInput, "/") {
		ips, err = expandCIDR(ipInput)
		if err != nil {
			fmt.Printf("Error expanding CIDR: %v\n", err)
			return
		}
	} else {
		ips = []string{ipInput}
	}

	var scanResults []ScanResult
	results := make(chan ScanResult)
	var wg sync.WaitGroup
	rateLimiter := make(chan struct{}, 10) // Rate limit concurrency to 10

	for _, ip := range ips {
		if !PingHost(ip) {
			fmt.Printf("Host %s is not reachable. Skipping...\n", ip)
			continue
		}

		for _, port := range ports {
			wg.Add(1)
			rateLimiter <- struct{}{}
			go func(ip string, port int) {
				defer func() { <-rateLimiter }()
				ScanPort(ip, port, results, &wg, 2*time.Second)
			}(ip, port)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		scanResults = append(scanResults, result)
		fmt.Printf("%s:%d - %s (%s) %s\n", result.IP, result.Port, result.Status, result.Service, result.Banner)
	}

	err = SaveResults(outputFile+"."+outputFormat, outputFormat, scanResults)
	if err != nil {
		fmt.Printf("Error saving results: %v\n", err)
	} else {
		fmt.Printf("Results saved to %s.%s\n", outputFile, outputFormat)
	}
}
