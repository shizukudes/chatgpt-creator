package register

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/verssache/chatgpt-creator/internal/email"
	"github.com/verssache/chatgpt-creator/internal/util"
)

// registerOne handles a single account registration.
func registerOne(workerID int, tag string, proxy, outputFile, defaultPassword, defaultDomain string, printMu, fileMu *sync.Mutex) (bool, string, string) {
	client, err := NewClient(proxy, tag, workerID, printMu, fileMu)
	if err != nil {
		return false, "", fmt.Sprintf("failed to create client: %v", err)
	}

	emailAddr, err := email.CreateTempEmail(defaultDomain)
	if err != nil {
		return false, "", fmt.Sprintf("failed to create temp email: %v", err)
	}

	password := defaultPassword
	if password == "" {
		password = util.GeneratePassword(14)
	}

	firstName, lastName := util.RandomName()
	birthdate := util.RandomBirthdate()

	client.print(fmt.Sprintf("Starting registration for %s", emailAddr))

	err = client.RunRegister(emailAddr, password, firstName+" "+lastName, birthdate)
	if err != nil {
		return false, emailAddr, err.Error()
	}

	// Append to file
	fileMu.Lock()
	defer fileMu.Unlock()

	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, emailAddr, fmt.Sprintf("failed to open output file: %v", err)
	}
	defer f.Close()

	line := fmt.Sprintf("%s|%s\n", emailAddr, password)
	if _, err := f.WriteString(line); err != nil {
		return false, emailAddr, fmt.Sprintf("failed to write to output file: %v", err)
	}

	return true, emailAddr, ""
}

// RunBatch runs concurrent registration tasks with retry until target success count is reached.
func RunBatch(totalAccounts int, outputFile string, maxWorkers int, proxy, defaultPassword, defaultDomain string) {
	var printMu sync.Mutex
	var fileMu sync.Mutex

	var remaining int64 = int64(totalAccounts)
	var successCount int64
	var failureCount int64
	var attemptNum int64

	startTime := time.Now()

	var wg sync.WaitGroup

	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				// Claim a slot before starting work
				if atomic.AddInt64(&remaining, -1) < 0 {
					// No more slots needed, put it back and exit
					atomic.AddInt64(&remaining, 1)
					return
				}

				attempt := atomic.AddInt64(&attemptNum, 1)
				tag := fmt.Sprintf("%d/%d", attempt, totalAccounts)

				success, emailAddr, errStr := registerOne(workerID, tag, proxy, outputFile, defaultPassword, defaultDomain, &printMu, &fileMu)
				if success {
					atomic.AddInt64(&successCount, 1)
					ts := time.Now().Format("15:04:05")
					printMu.Lock()
					fmt.Printf("[%s] [W%d] ✓ SUCCESS: %s\n", ts, workerID, emailAddr)
					printMu.Unlock()
				} else {
					atomic.AddInt64(&failureCount, 1)
					// Failed — return the slot so it gets retried
					atomic.AddInt64(&remaining, 1)
					ts := time.Now().Format("15:04:05")
					printMu.Lock()
					fmt.Printf("[%s] [W%d] ✗ FAILURE: %s | %s\n", ts, workerID, emailAddr, errStr)
					printMu.Unlock()
				}
			}
		}(w)
	}

	wg.Wait()

	elapsed := time.Since(startTime)
	elapsedStr := formatDuration(elapsed)

	fmt.Printf("\n--- Batch Registration Summary ---\n")
	fmt.Printf("Target:    %d\n", totalAccounts)
	fmt.Printf("Success:   %d\n", successCount)
	fmt.Printf("Attempts:  %d\n", attemptNum)
	fmt.Printf("Failures:  %d\n", failureCount)
	fmt.Printf("Elapsed:   %s\n", elapsedStr)
	fmt.Printf("----------------------------------\n")
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
