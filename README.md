# ChatGPT Account Registration Bot

Automated bulk ChatGPT account registration bot built with Go. Features concurrent workers, TLS fingerprint spoofing, automatic email generation, OTP verification, cookie export, and retry-until-success logic.

## Features

- **Concurrent Registration** — Configurable worker pool for parallel account creation
- **TLS Fingerprinting** — Randomized Chrome TLS profiles to avoid detection
- **Auto Email Generation** — Generates temporary emails via [generator.email](https://generator.email) or custom domains
- **OTP Verification** — Automatic email OTP retrieval and validation
- **Cookie Export** — Saves full session cookies (JSON format) for browser extension import
- **Browser-Friendly Cookies** — Separate export without cross-domain cookies (avoids import errors)
- **Retry Loop** — Automatically retries failed registrations until target count is reached
- **Proxy Support** — Optional HTTP/SOCKS proxy for all requests
- **Configurable** — JSON config file with interactive prompt overrides

## Requirements

- Go 1.21+

## Installation

```bash
git clone https://github.com/shizukudes/chatgpt-creator.git
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
[22:43:09] [W1] [1/5] Signin | 200
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

## Cookie Export

After successful registration, cookies are automatically saved to the `cookies/` directory:

```
cookies/
├── user@example.com.json           # Full cookies (all domains)
└── user@example.com-browser.json   # Browser-friendly (chatgpt.com only)
```

### Cookie Files

| File | Description | Import Method |
|------|-------------|---------------|
| `email@domain.com.json` | Full session cookies including auth tokens | EditThisCookie, Cookie-Editor (desktop) |
| `email@domain.com-browser.json` | ChatGPT-only cookies, no cross-domain | Cookie-Editor (mobile/Android), browser DevTools |

### Importing Cookies

**EditThisCookie (Chrome)**
1. Install [EditThisCookie](https://chrome.google.com/webstore/detail/editthiscookie/fngmhnnpilhplaeedifhccceomclgfbg)
2. Go to chatgpt.com
3. Click EditThisCookie icon → Import → select `email@domain.com.json`

**Cookie-Editor (Firefox/Chrome)**
1. Install [Cookie-Editor](https://cookie-editor.cgagnier.ca/)
2. Navigate to chatgpt.com
3. Click Cookie-Editor → Import → paste cookie JSON or select file

**Browser DevTools**
1. Open DevTools (F12) → Application → Cookies
2. Select chatgpt.com
3. Manually add each cookie from the JSON file

### Cookie Structure

```json
[
  {
    "name": "__Secure-next-auth.session-token",
    "value": "eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNI...",
    "domain": "chatgpt.com",
    "path": "/",
    "expires": "2025-01-15T12:00:00Z",
    "secure": true,
    "httpOnly": true
  }
]
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
│   │   ├── cookies.go       # Cookie export to JSON (full + browser-friendly)
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
├── cookies/                  # Auto-generated cookie JSON files
├── config.json               # Configuration file
├── go.mod
└── go.sum
```

## Disclaimer

This tool is provided for educational and research purposes only. Use of this tool to create accounts in violation of OpenAI's Terms of Service is solely at your own risk. The author assumes no responsibility for any misuse or consequences arising from the use of this software.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
