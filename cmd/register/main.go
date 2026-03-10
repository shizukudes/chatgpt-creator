package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/verssache/chatgpt-creator/internal/config"
	"github.com/verssache/chatgpt-creator/internal/register"
)

func main() {
	printBanner()

	// Load config
	cfg, err := config.Load("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	// 1. Proxy prompt
	proxy := cfg.Proxy
	if cfg.Proxy == "" {
		fmt.Printf("Proxy (enter to skip): ")
		proxyInput, _ := reader.ReadString('\n')
		proxy = strings.TrimSpace(proxyInput)
	}
	// 2. Total accounts prompt (required)
	fmt.Printf("Total accounts to register: ")
	totalInput, _ := reader.ReadString('\n')
	totalInput = strings.TrimSpace(totalInput)

	if totalInput == "" {
		fmt.Println("Error: total accounts is required.")
		os.Exit(1)
	}
	totalAccounts, err := strconv.Atoi(totalInput)
	if err != nil {
		fmt.Printf("Error: invalid number '%s'.\n", totalInput)
		os.Exit(1)
	}

	// 3. Max workers prompt
	defaultWorkers := 3
	fmt.Printf("Max concurrent workers (default: %d): ", defaultWorkers)
	workersInput, _ := reader.ReadString('\n')
	workersInput = strings.TrimSpace(workersInput)

	maxWorkers := defaultWorkers
	if workersInput != "" {
		if val, err := strconv.Atoi(workersInput); err == nil {
			maxWorkers = val
		}
	}

	// 4. Default password prompt
	defaultPassword := cfg.DefaultPassword
	if cfg.DefaultPassword == "" {
		fmt.Printf("Default password (current: (random), press Enter to use, or enter new): ")
		pwInput, _ := reader.ReadString('\n')
		pwInput = strings.TrimSpace(pwInput)

		if pwInput != "" {
			if len(pwInput) < 12 {
				fmt.Println("Error: password must be at least 12 characters.")
				os.Exit(1)
			}
			defaultPassword = pwInput
		}
	}
	// 5. Default domain prompt
	defaultDomain := cfg.DefaultDomain
	if cfg.DefaultDomain == "" {
		fmt.Printf("Default domain (current: (random from generator.email), press Enter to use, or enter new): ")
		domainInput, _ := reader.ReadString('\n')
		domainInput = strings.TrimSpace(domainInput)

		if domainInput != "" {
			defaultDomain = domainInput
		}
	}
	fmt.Println()
	fmt.Println("-------------------------------------------")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Proxy:          %s\n", proxy)
	fmt.Printf("  Total Accounts: %d\n", totalAccounts)
	fmt.Printf("  Max Workers:    %d\n", maxWorkers)
	if defaultPassword != "" {
		fmt.Printf("  Password:       %s\n", defaultPassword)
	} else {
		fmt.Printf("  Password:       (random)\n")
	}
	if defaultDomain != "" {
		fmt.Printf("  Domain:         %s\n", defaultDomain)
	} else {
		fmt.Printf("  Domain:         (random)\n")
	}
	fmt.Printf("  Output File:    %s\n", cfg.OutputFile)
	fmt.Println("-------------------------------------------")
	fmt.Println()

	// Run the batch
	register.RunBatch(totalAccounts, cfg.OutputFile, maxWorkers, proxy, defaultPassword, defaultDomain)
}

func printBanner() {
	banner := `
   _____ _           _    _____ _____ _______
  / ____| |         | |  / ____|  __ \__   __|
 | |    | |__   __ _| |_| |  __| |__) | | |
 | |    | '_ \ / _` + "`" + ` | __| | |_ |  ___/  | |
 | |____| | | | (_| | |_| |__| | |      | |
  \_____|_| |_|\__,_|\__|\_____|_|      |_|

      ChatGPT Account Registration Bot
               by @verssache
`
	fmt.Println(banner)
}
