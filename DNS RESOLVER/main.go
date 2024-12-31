package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	cache           = make(map[string][]net.IP)
	mu              sync.RWMutex
	cacheFile       = "dns_cache.json"
	customDNSServer string
)

// Main function
func main() {
	loadCache()

	for {
		showMenu()

		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Invalid input, please try again.")
			continue
		}

		switch choice {
		case 1:
			resolveRecords("A")
		case 2:
			resolveRecords("AAAA")
		case 3:
			resolveMXRecords()
		case 4:
			setCustomDNSServer()
		case 5:
			clearCache()
		case 6:
			showCache()
		case 7:
			fmt.Println("Exiting...")
			saveCache() // Save cache before exiting
			return
		default:
			fmt.Println("Invalid option, please select again.")
		}
	}
}

// Show the main menu
func showMenu() {
	fmt.Println("\nDNS Resolver")
	fmt.Println("1. Resolve A Record")
	fmt.Println("2. Resolve AAAA Record")
	fmt.Println("3. Resolve MX Record")
	fmt.Println("4. Set Custom DNS Server")
	fmt.Println("5. Clear Cache")
	fmt.Println("6. Show Cache")
	fmt.Println("7. Exit")
	fmt.Print("Select an option: ")
}

// Resolve DNS records for specified types
func resolveRecords(recordType string) {
	fmt.Print("Enter domain name (comma separated for multiple): ")
	var domains string
	fmt.Scan(&domains)

	domainList := splitAndTrim(domains)
	var wg sync.WaitGroup

	for _, domain := range domainList {
		if !isValidDomain(domain) {
			fmt.Printf("Invalid domain name: %s\n", domain)
			continue
		}

		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			if recordType == "A" {
				resolveA(domain)
			} else {
				resolveAAAA(domain)
			}
		}(domain)
	}

	wg.Wait()
}

// Resolve A records for a domain
func resolveA(domain string) {
	mu.RLock()
	if cachedIPs, found := cache[domain]; found {
		fmt.Printf("Cached A Record for %s: %v\n", domain, cachedIPs)
		mu.RUnlock()
		return
	}
	mu.RUnlock()

	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("Failed to resolve A record for %s: %v\n", domain, err)
		return
	}

	var aRecords []net.IP
	for _, ip := range ips {
		if ip.To4() != nil {
			aRecords = append(aRecords, ip)
		}
	}

	mu.Lock()
	cache[domain] = aRecords
	mu.Unlock()

	fmt.Printf("Resolved A Record for %s: %v\n", domain, aRecords)
}

// Resolve AAAA records for a domain
func resolveAAAA(domain string) {
	mu.RLock()
	if cachedIPs, found := cache[domain]; found {
		fmt.Printf("Cached AAAA Record for %s: %v\n", domain, cachedIPs)
		mu.RUnlock()
		return
	}
	mu.RUnlock()

	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("Failed to resolve AAAA record for %s: %v\n", domain, err)
		return
	}

	var aaaaRecords []net.IP
	for _, ip := range ips {
		if ip.To16() != nil && ip.To4() == nil {
			aaaaRecords = append(aaaaRecords, ip)
		}
	}

	mu.Lock()
	cache[domain] = aaaaRecords
	mu.Unlock()

	fmt.Printf("Resolved AAAA Record for %s: %v\n", domain, aaaaRecords)
}

// Resolve MX records for a domain
func resolveMXRecords() {
	fmt.Print("Enter domain name: ")
	var domain string
	fmt.Scan(&domain)

	if !isValidDomain(domain) {
		fmt.Printf("Invalid domain name: %s\n", domain)
		return
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		fmt.Printf("Failed to resolve MX records for %s: %v\n", domain, err)
		return
	}

	fmt.Printf("Resolved MX Records for %s:\n", domain)
	for _, mx := range mxRecords {
		fmt.Printf("  - %s (priority: %d)\n", mx.Host, mx.Pref)
	}
}

// Set a custom DNS server
func setCustomDNSServer() {
	fmt.Print("Enter custom DNS server (e.g., 8.8.8.8): ")
	fmt.Scan(&customDNSServer)
	fmt.Printf("Custom DNS server set to: %s\n", customDNSServer)
}

// Clear the DNS cache
func clearCache() {
	mu.Lock()
	cache = make(map[string][]net.IP)
	mu.Unlock()
	fmt.Println("Cache cleared.")
}

// Show the current cache
func showCache() {
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println("Current Cache:")
	for domain, ips := range cache {
		fmt.Printf("  %s: %v\n", domain, ips)
	}
}

// Load cache from file
func loadCache() {
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error loading cache: %v\n", err)
		}
		return
	}

	if err := json.Unmarshal(data, &cache); err != nil {
		fmt.Printf("Error parsing cache: %v\n", err)
	}
}

// Save cache to file
func saveCache() {
	mu.RLock()
	defer mu.RUnlock()
	data, err := json.Marshal(cache)
	if err != nil {
		fmt.Printf("Error saving cache: %v\n", err)
		return
	}

	if err := ioutil.WriteFile(cacheFile, data, 0644); err != nil {
		fmt.Printf("Error writing cache to file: %v\n", err)
	}
}

// Split input by comma and trim whitespace
func splitAndTrim(input string) []string {
	parts := strings.Split(input, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// Validate domain name using regex
func isValidDomain(domain string) bool {
	const domainPattern = `^[A-Za-z0-9-]{1,63}(\.[A-Za-z]{2,6})+$`
	re := regexp.MustCompile(domainPattern)
	return re.MatchString(domain)
}
