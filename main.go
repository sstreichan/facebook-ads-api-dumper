package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	baseURL    = "https://graph.facebook.com/v19.0"
	apiVersion = "v19.0"
)

type Config struct {
	AccessToken string
	OutputDir   string
}

type AdAccount struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
	Currency  string `json:"currency"`
}

type AdAccountsResponse struct {
	Data   []AdAccount `json:"data"`
	Paging struct {
		Cursors struct {
			Before string `json:"before"`
			After  string `json:"after"`
		} `json:"cursors"`
	} `json:"paging"`
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

func (c *APIClient) dumpResponse(name string, data []byte, accountDir string) error {
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
	if c.config.OutputDir != "" && accountDir != "" {
		filename := fmt.Sprintf("%s/%s_%d.json", accountDir, name, time.Now().Unix())
		if err := os.WriteFile(filename, formatted, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		log.Printf("Saved to: %s", filename)
	}
	
	return nil
}

func (c *APIClient) fetchAdAccounts() ([]AdAccount, error) {
	endpoint := "me/adaccounts?fields=id,name,account_id,currency,timezone_name,account_status"
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return nil, err
	}
	
	var response AdAccountsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("parsing ad accounts response: %w", err)
	}
	
	// Also dump the raw response
	c.dumpResponse("all_ad_accounts", data, c.config.OutputDir)
	
	return response.Data, nil
}

func (c *APIClient) fetchAdAccount(accountID string, accountDir string) error {
	endpoint := fmt.Sprintf("%s?fields=id,name,account_id,currency,timezone_name,business,account_status", accountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("ad_account", data, accountDir)
}

func (c *APIClient) fetchCampaigns(accountID string, accountDir string) error {
	endpoint := fmt.Sprintf("%s/campaigns?fields=id,name,status,objective,created_time,updated_time", accountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("campaigns", data, accountDir)
}

func (c *APIClient) fetchAdSets(accountID string, accountDir string) error {
	endpoint := fmt.Sprintf("%s/adsets?fields=id,name,status,campaign_id,daily_budget,lifetime_budget,created_time", accountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("adsets", data, accountDir)
}

func (c *APIClient) fetchAds(accountID string, accountDir string) error {
	endpoint := fmt.Sprintf("%s/ads?fields=id,name,status,adset_id,creative,created_time", accountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("ads", data, accountDir)
}

func (c *APIClient) fetchInsights(accountID string, accountDir string) error {
	endpoint := fmt.Sprintf("%s/insights?fields=impressions,clicks,spend,ctr,cpc,date_start,date_stop&level=account&time_range={'since':'2026-01-01','until':'2026-02-03'}", accountID)
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return err
	}
	return c.dumpResponse("insights", data, accountDir)
}

func (c *APIClient) processAccount(account AdAccount) error {
	log.Printf("\n========================================")
	log.Printf("Processing Account: %s (%s)", account.Name, account.AccountID)
	log.Printf("========================================\n")
	
	// Create account-specific directory if output is enabled
	var accountDir string
	if c.config.OutputDir != "" {
		// Sanitize account name for directory
		safeName := strings.Map(func(r rune) rune {
			if r == '/' || r == '\\' || r == ':' {
				return '_'
			}
			return r
		}, account.Name)
		accountDir = filepath.Join(c.config.OutputDir, fmt.Sprintf("%s_%s", account.AccountID, safeName))
		if err := os.MkdirAll(accountDir, 0755); err != nil {
			return fmt.Errorf("creating account directory: %w", err)
		}
	}
	
	// Fetch all resources for this account
	if err := c.fetchAdAccount(account.ID, accountDir); err != nil {
		log.Printf("Error fetching ad account details: %v", err)
	}
	
	if err := c.fetchCampaigns(account.ID, accountDir); err != nil {
		log.Printf("Error fetching campaigns: %v", err)
	}
	
	if err := c.fetchAdSets(account.ID, accountDir); err != nil {
		log.Printf("Error fetching ad sets: %v", err)
	}
	
	if err := c.fetchAds(account.ID, accountDir); err != nil {
		log.Printf("Error fetching ads: %v", err)
	}
	
	if err := c.fetchInsights(account.ID, accountDir); err != nil {
		log.Printf("Error fetching insights: %v", err)
	}
	
	return nil
}

func main() {
	accessToken := flag.String("token", "", "Facebook access token (required)")
	outputDir := flag.String("output", "", "Output directory for JSON files (optional)")
	flag.Parse()
	
	if *accessToken == "" {
		flag.Usage()
		log.Fatal("The -token flag is required")
	}
	
	// Create output directory if specified
	if *outputDir != "" {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}
	
	config := Config{
		AccessToken: *accessToken,
		OutputDir:   *outputDir,
	}
	
	client := NewAPIClient(config)
	
	log.Println("Starting Facebook Ads API data dump...")
	log.Println("Discovering accessible ad accounts...")
	
	// Fetch all accessible ad accounts
	accounts, err := client.fetchAdAccounts()
	if err != nil {
		log.Fatalf("Failed to fetch ad accounts: %v", err)
	}
	
	if len(accounts) == 0 {
		log.Println("No ad accounts found for this access token.")
		log.Println("Make sure your token has 'ads_read' permission and you have access to at least one ad account.")
		return
	}
	
	log.Printf("Found %d accessible ad account(s)\n", len(accounts))
	
	// Process each account
	successCount := 0
	for i, account := range accounts {
		log.Printf("\nProcessing %d/%d: %s", i+1, len(accounts), account.Name)
		if err := client.processAccount(account); err != nil {
			log.Printf("Error processing account %s: %v", account.Name, err)
		} else {
			successCount++
		}
	}
	
	log.Printf("\n========================================")
	log.Printf("Data dump complete!")
	log.Printf("Successfully processed %d/%d accounts", successCount, len(accounts))
	log.Printf("========================================\n")
}
