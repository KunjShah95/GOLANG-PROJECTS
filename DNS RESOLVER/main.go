package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	cacheFile = "dns_cache.json"
	logFile   = "dns_resolver.log"
)

type DNSCache struct {
	mu    sync.RWMutex
	cache map[string]DNSRecord
}

type DNSRecord struct {
	Domain      string
	RecordTypes map[string][]string
	Timestamp   time.Time
}

type DNSResolver struct {
	cache        *DNSCache
	logger       *log.Logger
	customServer []string
}

func NewDNSResolver() *DNSResolver {
	return &DNSResolver{
		cache: &DNSCache{
			cache: make(map[string]DNSRecord),
		},
		logger: setupLogger(),
	}
}

func setupLogger() *log.Logger {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return log.New(os.Stdout, "DNS_RESOLVER: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return log.New(file, "DNS_RESOLVER: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (r *DNSResolver) ResolveDomain(domain, recordType string) ([]string, error) {
	// Check cache first
	if cachedRecord := r.checkCache(domain, recordType); cachedRecord != nil {
		return cachedRecord, nil
	}

	var results []string
	var err error

	// Prioritize custom DNS servers if provided
	if len(r.customServer) > 0 {
		results, err = r.resolveWithCustomDNS(domain, recordType, r.customServer)
	} else {
		results, err = r.resolveWithDefaultDNS(domain, recordType)
	}

	if err != nil {
		r.logger.Printf("Resolution error for %s (%s): %v", domain, recordType, err)
		return nil, err
	}

	// Cache the results
	r.cacheResults(domain, recordType, results)
	return results, nil
}

func (r *DNSResolver) resolveWithCustomDNS(_, _ string, _ []string) ([]string, error) {
	// Custom DNS resolution logic here
	return nil, fmt.Errorf("custom DNS resolution not implemented")
}

func (r *DNSResolver) resolveWithDefaultDNS(domain, recordType string) ([]string, error) {
	var results []string
	var err error

	switch recordType {
	case "A":
		results, err = r.resolveARecord(domain)
	case "AAAA":
		results, err = r.resolveAAAARecord(domain)
	case "MX":
		results, err = r.resolveMXRecord(domain)
	case "TXT":
		results, err = r.resolveTXTRecord(domain)
	case "NS":
		results, err = r.resolveNSRecord(domain)
	case "PTR":
		results, err = r.resolvePTRRecord(domain) // Reverse DNS lookup
	default:
		err = fmt.Errorf("unsupported record type: %s", recordType)
	}

	return results, err
}

func (r *DNSResolver) resolveARecord(domain string) ([]string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			results = append(results, ipv4.String())
		}
	}
	return results, nil
}

func (r *DNSResolver) resolveAAAARecord(domain string) ([]string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, ip := range ips {
		if ipv6 := ip.To16(); ipv6 != nil && ip.To4() == nil {
			results = append(results, ipv6.String())
		}
	}
	return results, nil
}

func (r *DNSResolver) resolveMXRecord(domain string) ([]string, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, mx := range mxRecords {
		results = append(results, fmt.Sprintf("%s (Priority: %d)", mx.Host, mx.Pref))
	}
	return results, nil
}

func (r *DNSResolver) resolveTXTRecord(domain string) ([]string, error) {
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return nil, err
	}
	return txtRecords, nil
}

func (r *DNSResolver) resolveNSRecord(domain string) ([]string, error) {
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, ns := range nsRecords {
		results = append(results, ns.Host)
	}
	return results, nil
}

func (r *DNSResolver) resolvePTRRecord(ip string) ([]string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (r *DNSResolver) buildCLI() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "dnsresolver",
		Short: "Advanced DNS Resolution Tool",
		Long:  "A comprehensive DNS resolution and investigation tool",
	}

	var resolveCmd = &cobra.Command{
		Use:   "resolve [domain] [record-type]",
		Short: "Resolve DNS records for a domain",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			recordType := strings.ToUpper(args[1])

			results, err := r.ResolveDomain(domain, recordType)
			if err != nil {
				color.Red("Error: %v", err)
				return
			}

			color.Green("Results for %s (%s):", domain, recordType)
			for _, result := range results {
				fmt.Println(result)
			}
		},
	}

	rootCmd.AddCommand(resolveCmd)

	var customDNSCmd = &cobra.Command{
		Use:   "setdns [dns1] [dns2] [..]",
		Short: "Set custom DNS servers",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			r.customServer = args
			color.Green("Custom DNS servers set: %v", r.customServer)
		},
	}

	rootCmd.AddCommand(customDNSCmd)

	return rootCmd
}

func main() {
	homeDir, err := getUserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
	} else {
		fmt.Printf("User home directory: %s\n", homeDir)
	}

	resolver := NewDNSResolver()

	// Load cache from file
	resolver.loadCache()

	// Command-line mode
	if err := resolver.buildCLI().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Save cache before exiting
	resolver.saveCache()
}

// Cache handling
func (r *DNSResolver) checkCache(domain, recordType string) []string {
	r.cache.mu.RLock()
	defer r.cache.mu.RUnlock()
	if record, found := r.cache.cache[domain]; found {
		if results, ok := record.RecordTypes[recordType]; ok {
			return results
		}
	}
	return nil
}

func (r *DNSResolver) cacheResults(domain, recordType string, results []string) {
	r.cache.mu.Lock()
	defer r.cache.mu.Unlock()

	if _, found := r.cache.cache[domain]; !found {
		r.cache.cache[domain] = DNSRecord{
			Domain:      domain,
			RecordTypes: make(map[string][]string),
			Timestamp:   time.Now(),
		}
	}

	r.cache.cache[domain].RecordTypes[recordType] = results
}

func (r *DNSResolver) loadCache() {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return // Ignore if the file doesn't exist
		}
		fmt.Printf("Error loading cache: %v\n", err)
		return
	}

	if err := json.Unmarshal(data, &r.cache.cache); err != nil {
		fmt.Printf("Error parsing cache: %v\n", err)
	}
}

func (r *DNSResolver) saveCache() {
	r.cache.mu.RLock()
	defer r.cache.mu.RUnlock()

	data, err := json.Marshal(r.cache.cache)
	if err != nil {
		fmt.Printf("Error saving cache: %v\n", err)
		return
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		fmt.Printf("Error writing cache to file: %v\n", err)
	}
}

// Additional utility functions for advanced DNS operations

func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

// Additional utility functions for advanced DNS operations
// Removed unused function validateDomain
