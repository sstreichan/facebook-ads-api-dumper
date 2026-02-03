package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	baseURL    = "https://graph.facebook.com/v19.0"
	apiVersion = "v19.0"
)

type Config struct {
	AccessToken string
	AdAccountID string
	OutputDir   string
}

type APIClient struct {
	config     Config
	httpClient *http.Client
}

func NewAPIClient(config Config) *APIClient {
	return &APIClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) makeRequest(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?access_token=%s", baseURL, endpoint, c.config.AccessToken)
	
	log.Printf("Requesting: %s", endpoint)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}
	
	return body, nil
}

func (c *APIClient) dumpResponse(name string, data []byte) error {
	// Pretty print to console
	var prettyJSON interface{}
	if err := json.Unmarshal(data, &prettyJSON); err != nil {
		log.Printf("Warning: Invalid JSON from %s", name)
		fmt.Printf("\n=== %s (RAW) ===\n%s\n\n", name, string(data))
		return nil
	}
	
	formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
	fmt.Printf("\n=== %s ===\n%s\n\n", name, string(formatted))
	
	// Save to file if output directory specified
	if c.config.OutputDir != "" {
		filename := fmt.Sprintf("%s/%s_%d.json", c.config.OutputDir, name, time.Now().Unix())
		if err := os.WriteFile(filename, formatted, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		log.Printf("Saved to: %s", filename)
	}
	
	return nil
}

func (c *APIClient) fetchAdAccount() error {
	endpoint := fmt.Sprintf("act_%s?fields=id,name,account_id,currency,timezone_name,business", c.config.AdAccountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("ad_account", data)
}

func (c *APIClient) fetchCampaigns() error {
	endpoint := fmt.Sprintf("act_%s/campaigns?fields=id,name,status,objective,created_time,updated_time", c.config.AdAccountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("campaigns", data)
}

func (c *APIClient) fetchAdSets() error {
	endpoint := fmt.Sprintf("act_%s/adsets?fields=id,name,status,campaign_id,daily_budget,lifetime_budget,created_time", c.config.AdAccountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("adsets", data)
}

func (c *APIClient) fetchAds() error {
	endpoint := fmt.Sprintf("act_%s/ads?fields=id,name,status,adset_id,creative,created_time", c.config.AdAccountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("ads", data)
}

func (c *APIClient) fetchInsights() error {
	endpoint := fmt.Sprintf("act_%s/insights?fields=impressions,clicks,spend,ctr,cpc,date_start,date_stop&level=account&time_range={'since':'2026-01-01','until':'2026-02-03'}", c.config.AdAccountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("insights", data)
}

func main() {
	accessToken := flag.String("token", "", "Facebook access token (required)")
	adAccountID := flag.String("account", "", "Ad account ID without 'act_' prefix (required)")
	outputDir := flag.String("output", "", "Output directory for JSON files (optional)")
	flag.Parse()
	
	if *accessToken == "" || *adAccountID == "" {
		flag.Usage()
		log.Fatal("Both -token and -account flags are required")
	}
	
	// Create output directory if specified
	if *outputDir != "" {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}
	
	config := Config{
		AccessToken: *accessToken,
		AdAccountID: *adAccountID,
		OutputDir:   *outputDir,
	}
	
	client := NewAPIClient(config)
	
	log.Println("Starting Facebook Ads API data dump...")
	
	// Fetch all resources
	if err := client.fetchAdAccount(); err != nil {
		log.Printf("Error fetching ad account: %v", err)
	}
	
	if err := client.fetchCampaigns(); err != nil {
		log.Printf("Error fetching campaigns: %v", err)
	}
	
	if err := client.fetchAdSets(); err != nil {
		log.Printf("Error fetching ad sets: %v", err)
	}
	
	if err := client.fetchAds(); err != nil {
		log.Printf("Error fetching ads: %v", err)
	}
	
	if err := client.fetchInsights(); err != nil {
		log.Printf("Error fetching insights: %v", err)
	}
	
	log.Println("Data dump complete!")
}
