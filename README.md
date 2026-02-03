# Facebook Ads API Dumper

A lightweight Go program that fetches and dumps raw JSON responses from the Facebook Marketing API. Perfect for debugging, data analysis, or building integrations.

## Features

- 🔍 Fetches ad account details, campaigns, ad sets, ads, and insights
- 📄 Outputs formatted JSON to console and/or files
- 🚀 No external dependencies (uses Go standard library only)
- ⚙️ Simple command-line interface
- 💾 Optional file output with timestamps

## Prerequisites

1. **Facebook Developer Account**: Create an app at [developers.facebook.com](https://developers.facebook.com)
2. **Marketing API Access**: Add the Marketing API product to your app
3. **Access Token**: Generate a token with `ads_read` permission via [Graph API Explorer](https://developers.facebook.com/tools/explorer/)
4. **Ad Account ID**: Find your account ID in [Facebook Ads Manager](https://business.facebook.com/adsmanager) (format: `act_XXXXXXXXXX`)

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
go run main.go -token YOUR_TOKEN -account YOUR_ACCOUNT_ID
```

## Usage

### Basic Usage (Console Output)

```bash
./fb-ads-dump -token YOUR_ACCESS_TOKEN -account 1234567890
```

### Save to Files

```bash
./fb-ads-dump -token YOUR_ACCESS_TOKEN -account 1234567890 -output ./dumps
```

### Command-Line Flags

- `-token` (required): Your Facebook access token
- `-account` (required): Ad account ID without the `act_` prefix
- `-output` (optional): Directory to save JSON files with timestamps

## Example Output

```json
=== ad_account ===
{
  "id": "act_1234567890",
  "name": "My Ad Account",
  "account_id": "1234567890",
  "currency": "USD",
  "timezone_name": "America/Los_Angeles"
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
  ],
  "paging": {...}
}
```

## What Data is Retrieved

- **Ad Account**: Basic account information (name, currency, timezone)
- **Campaigns**: All campaigns with status, objective, and timestamps
- **Ad Sets**: All ad sets with budget information and campaign associations
- **Ads**: All ads with creative details and status
- **Insights**: Account-level performance metrics (impressions, clicks, spend, CTR, CPC)

## API Version

Currently uses Facebook Graph API **v19.0**. Update the `apiVersion` constant in `main.go` to use a different version.

## Security Notes

⚠️ **Never commit your access token to version control!**

- Use command-line flags for tokens (not environment variables in scripts)
- Access tokens grant broad permissions - store them securely
- Long-lived tokens expire after 60 days - implement refresh logic for production
- Consider using a configuration file (add to `.gitignore`) for repeated use

## Limitations

- **Pagination**: Currently fetches only the first page of results (typically 25 items)
- **Rate Limits**: No automatic retry or backoff for rate limit errors
- **Token Refresh**: Manual token renewal required every 60 days

## Future Enhancements

- [ ] Add cursor-based pagination for complete data retrieval
- [ ] Implement exponential backoff for rate limit handling
- [ ] Add token refresh mechanism
- [ ] Support configuration files (YAML/JSON)
- [ ] Add filtering options (date ranges, status filters)
- [ ] Parallel fetching for multiple ad accounts

## Error Handling

The program continues execution even if individual requests fail, logging errors for each endpoint. Common errors:

- **OAuth errors**: Invalid or expired access token
- **Permission errors**: Token lacks required permissions (`ads_read`)
- **Rate limit errors**: Too many requests - wait and retry
- **Invalid account**: Ad account ID doesn't exist or you don't have access

## Contributing

Contributions welcome! Please feel free to submit issues or pull requests.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Resources

- [Facebook Marketing API Documentation](https://developers.facebook.com/docs/marketing-api)
- [Graph API Explorer](https://developers.facebook.com/tools/explorer/)
- [Marketing API Rate Limits](https://developers.facebook.com/docs/marketing-api/overview/rate-limiting)

## Author

Created by [sstreichan](https://github.com/sstreichan)