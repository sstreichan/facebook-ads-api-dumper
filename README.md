# Facebook Ads API Dumper

A lightweight Go program that automatically discovers and dumps raw JSON responses from all Facebook ad accounts you have access to. Perfect for debugging, data analysis, or building integrations.

## Features

- 🔍 **Auto-discovers** all ad accounts accessible to your token
- 📊 Fetches ad account details, campaigns, ad sets, ads, and insights
- 📄 Outputs formatted JSON to console and/or files
- 🚀 No external dependencies (uses Go standard library only)
- ⚙️ Simple command-line interface
- 💾 Optional file output with account-specific directories
- 📈 Summary statistics showing accounts processed

## Prerequisites

1. **Facebook Developer Account**: Create an app at [developers.facebook.com](https://developers.facebook.com)
2. **Marketing API Access**: Add the Marketing API product to your app
3. **Access Token**: Generate a token with `ads_read` permission via [Graph API Explorer](https://developers.facebook.com/tools/explorer/)
4. **Ad Account Access**: Ensure you have access to at least one ad account (admin, advertiser, or analyst role)

## Installation

### Option 1: Clone and Build

```bash
git clone https://github.com/sstreichan/facebook-ads-api-dumper.git
cd facebook-ads-api-dumper
go build -o fb-ads-dump main.go
```

### Option 2: Install Directly

```bash
go install github.com/sstreichan/facebook-ads-api-dumper@latest
```

### Option 3: Run Without Installing

```bash
go run main.go -token YOUR_ACCESS_TOKEN
```

## Usage

### Basic Usage (Console Output Only)

```bash
./fb-ads-dump -token YOUR_ACCESS_TOKEN
```

The program will automatically:
1. Discover all ad accounts you have access to
2. Fetch data from each account (campaigns, ad sets, ads, insights)
3. Display formatted JSON to console

### Save to Files

```bash
./fb-ads-dump -token YOUR_ACCESS_TOKEN -output ./dumps
```

This creates a directory structure like:
```
dumps/
├── all_ad_accounts_1738594025.json
├── 1234567890_My_Ad_Account/
│   ├── ad_account_1738594026.json
│   ├── campaigns_1738594027.json
│   ├── adsets_1738594028.json
│   ├── ads_1738594029.json
│   └── insights_1738594030.json
└── 9876543210_Another_Account/
    └── ...
```

### Command-Line Flags

- `-token` (required): Your Facebook access token with `ads_read` permission
- `-output` (optional): Directory to save JSON files organized by account

## Example Output

```
Starting Facebook Ads API data dump...
Discovering accessible ad accounts...
Found 2 accessible ad account(s)

Processing 1/2: My Ad Account

========================================
Processing Account: My Ad Account (1234567890)
========================================

=== ad_account ===
{
  "id": "act_1234567890",
  "name": "My Ad Account",
  "account_id": "1234567890",
  "currency": "USD",
  "timezone_name": "America/Los_Angeles",
  "account_status": 1
}

=== campaigns ===
{
  "data": [
    {
      "id": "120212345678901234",
      "name": "Summer Campaign 2026",
      "status": "ACTIVE",
      "objective": "OUTCOME_TRAFFIC"
    }
  ]
}

...

========================================
Data dump complete!
Successfully processed 2/2 accounts
========================================
```

## What Data is Retrieved

For **each accessible ad account**:

- **Ad Account**: Basic account information (name, currency, timezone, status)
- **Campaigns**: All campaigns with status, objective, and timestamps
- **Ad Sets**: All ad sets with budget information and campaign associations
- **Ads**: All ads with creative details and status
- **Insights**: Account-level performance metrics for recent period (impressions, clicks, spend, CTR, CPC)

## API Version

Currently uses Facebook Graph API **v19.0**. Update the `apiVersion` constant in `main.go` to use a different version.

## Security Notes

⚠️ **Never commit your access token to version control!**

- Use command-line flags for tokens (not hardcoded values)
- Access tokens grant broad permissions - store them securely
- The `ads_read` permission allows reading all ad account data you have access to
- Long-lived tokens expire after 60 days - implement refresh logic for production
- Consider using environment variables or a configuration file (add to `.gitignore`)

## Limitations

- **Pagination**: Currently fetches only the first page of results per endpoint (typically 25 items)
- **Rate Limits**: No automatic retry or backoff for rate limit errors
- **Token Refresh**: Manual token renewal required every 60 days
- **Date Range**: Insights are fetched for a hardcoded date range (edit `fetchInsights()` to customize)

## Future Enhancements

- [ ] Add cursor-based pagination for complete data retrieval
- [ ] Implement exponential backoff for rate limit handling
- [ ] Add token refresh mechanism
- [ ] Support configuration files (YAML/JSON)
- [ ] Add filtering options (date ranges, status filters, specific accounts)
- [ ] Parallel fetching for faster multi-account processing
- [ ] Add progress bars for long-running operations

## Error Handling

The program continues execution even if individual requests fail, logging errors for each endpoint. Common errors:

- **OAuth errors**: Invalid or expired access token
- **Permission errors**: Token lacks `ads_read` permission
- **Rate limit errors**: Too many requests - wait and retry
- **Empty accounts list**: No accessible ad accounts or missing permissions

## Troubleshooting

### "No ad accounts found"

- Verify your token has `ads_read` permission in Graph API Explorer
- Ensure you have at least admin, advertiser, or analyst access to an ad account
- Check that your token hasn't expired (long-lived tokens last 60 days)

### "API error (status 190)"

- Your access token is invalid or expired
- Generate a new token from Graph API Explorer

### "API error (status 17)"

- Rate limit reached
- Wait a few minutes before retrying

## Contributing

Contributions welcome! Please feel free to submit issues or pull requests.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Resources

- [Facebook Marketing API Documentation](https://developers.facebook.com/docs/marketing-api)
- [Graph API Explorer](https://developers.facebook.com/tools/explorer/)
- [Marketing API Rate Limits](https://developers.facebook.com/docs/marketing-api/overview/rate-limiting)
- [Ad Account Permissions](https://www.facebook.com/business/help/442345745885606)

## Author

Created by [sstreichan](https://github.com/sstreichan)
