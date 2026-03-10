# ChatGPT Account Registration Bot

Automated bulk ChatGPT account registration bot built with Go. Features concurrent workers, TLS fingerprint spoofing, automatic email generation, OTP verification, and retry-until-success logic.

## Features

- **Concurrent Registration** — Configurable worker pool for parallel account creation
- **TLS Fingerprinting** — Randomized Chrome TLS profiles to avoid detection
- **Auto Email Generation** — Generates temporary emails via [generator.email](https://generator.email) or custom domains
- **OTP Verification** — Automatic email OTP retrieval and validation
- **Retry Loop** — Automatically retries failed registrations until target count is reached
- **Proxy Support** — Optional HTTP/SOCKS proxy for all requests
- **Configurable** — JSON config file with interactive prompt overrides

## Requirements

- Go 1.21+

## Installation

```bash
git clone https://github.com/verssache/chatgpt-creator.git
cd chatgpt-creator
go mod download
```

## Usage

```bash
go run cmd/register/main.go
```

### Interactive Prompts

```
Proxy (enter to skip):
Total accounts to register: 5
Max concurrent workers (default: 3): 2
Default password (current: (random), press Enter to use, or enter new):
Default domain (current: (random from generator.email), press Enter to use, or enter new):
```

### Example Output

```
[22:43:08] [W1] [1/5] Starting registration flow...
[22:43:09] [W1] [1/5] Visit Homepage (Try 1) | 200
[22:43:09] [W1] [1/5] Get CSRF | 200
[22:43:10] [W1] [1/5] Signin | 200
[22:43:12] [W1] [1/5] Authorize | 200
[22:43:15] [W1] [1/5] Register | 200
[22:43:17] [W1] [1/5] Send OTP | 200
[22:43:19] [W1] [1/5] Validate OTP [483291] | 200
[22:43:24] [W1] [1/5] Create Account | 200
[22:43:33] [W1] [1/5] Callback | 200
[22:43:33] [W1] SUCCESS: johndoe8x2kq@smartmail.de

--- Batch Registration Summary ---
Target:    5
Success:   5
Attempts:  6
Failures:  1
Elapsed:   1m 45s
----------------------------------
```

## Configuration

Create a `config.json` in the project root (optional):

```json
{
  "proxy": "",
  "output_file": "results.txt",
  "default_password": "",
  "default_domain": ""
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `proxy` | string | `""` | HTTP/SOCKS proxy URL. Leave empty for direct connection |
| `output_file` | string | `results.txt` | File path for saving registered accounts |
| `default_password` | string | `""` | Password for all accounts. Must be 12+ chars. Empty = random |
| `default_domain` | string | `""` | Email domain to use. Empty = random from generator.email |

Environment variable `PROXY` overrides the config file proxy value.

## Output Format

Registered accounts are saved to the output file in the format:

```
email|password
```

## Project Structure

```
.
├── cmd/
│   └── register/
│       └── main.go          # Entry point, interactive prompts
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration loading & validation
│   ├── register/
│   │   ├── batch.go         # Batch orchestration, worker pool, retry logic
│   │   ├── client.go        # HTTP client with TLS fingerprinting
│   │   └── flow.go          # Registration flow (CSRF → signup → OTP → callback)
│   ├── email/
│   │   └── generator.go     # Temporary email generation
│   ├── chrome/
│   │   └── profiles.go      # Chrome TLS profile randomization
│   └── util/
│       ├── helpers.go       # Utility functions
│       ├── names.go         # Random name generation (gofakeit)
│       ├── password.go      # Random password generation
│       └── trace.go         # Datadog trace headers
├── config.json               # Configuration file
├── go.mod
└── go.sum
```

## Disclaimer

This tool is provided for educational and research purposes only. Use of this tool to create accounts in violation of OpenAI's Terms of Service is solely at your own risk. The author assumes no responsibility for any misuse or consequences arising from the use of this software.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
