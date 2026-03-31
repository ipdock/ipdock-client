package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var version = "dev"

type updateResponse struct {
	Status string `json:"status"`
	IP     string `json:"ip"`
}

var ipServices = []string{
	"https://api.ipify.org",
	"https://icanhazip.com",
	"https://checkip.amazonaws.com",
	"https://ifconfig.me/ip",
}

func getPublicIP() (string, error) {
	// Shuffle to spread load
	shuffled := make([]string, len(ipServices))
	copy(shuffled, ipServices)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	var lastErr error
	client := &http.Client{Timeout: 10 * time.Second}
	for _, svc := range shuffled {
		resp, err := client.Get(svc)
		if err != nil {
			lastErr = err
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		ip := strings.TrimSpace(string(body))
		if ip != "" {
			return ip, nil
		}
	}
	return "", fmt.Errorf("all IP services failed: %v", lastErr)
}

func sendUpdate(apiURL, token, ip string) (*updateResponse, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	url := fmt.Sprintf("%s/api/update?ip=%s", apiURL, ip)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "ipdock-client/"+version)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var result updateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %s", string(body))
	}
	return &result, nil
}

func main() {
	token := os.Getenv("IPDOCK_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "[ipdock] IPDOCK_TOKEN is required")
		os.Exit(1)
	}

	apiURL := os.Getenv("IPDOCK_API_URL")
	if apiURL == "" {
		apiURL = "https://api.ipdock.io"
	}
	apiURL = strings.TrimRight(apiURL, "/")

	intervalSec := 60
	if v := os.Getenv("IPDOCK_INTERVAL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 10 {
			intervalSec = n
		}
	}

	fmt.Printf("[ipdock] Starting client (version=%s, interval=%ds, api=%s)\n", version, intervalSec, apiURL)

	var lastIP string

	// Calculate next forced heartbeat time (24h + 0-6h random jitter)
	nextHeartbeat := func() time.Time {
		jitter := time.Duration(rand.Intn(6*60)) * time.Minute // 0-6 hours
		return time.Now().Add(24*time.Hour + jitter)
	}
	heartbeatDeadline := nextHeartbeat()

	for {
		ip, err := getPublicIP()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ipdock] Failed to detect public IP: %v\n", err)
			time.Sleep(time.Duration(intervalSec) * time.Second)
			continue
		}

		// Decide whether to send an update
		forceHeartbeat := time.Now().After(heartbeatDeadline)

		if ip != lastIP {
			// IP changed — always update
			result, err := sendUpdate(apiURL, token, ip)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ipdock] Update failed: %v\n", err)
			} else {
				fmt.Printf("[ipdock] IP changed: %s → %s (status=%s)\n", lastIP, ip, result.Status)
				lastIP = ip
				heartbeatDeadline = nextHeartbeat()
			}
		} else if forceHeartbeat {
			// IP unchanged but heartbeat interval elapsed — send check-in
			result, err := sendUpdate(apiURL, token, ip)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ipdock] Heartbeat failed: %v\n", err)
			} else {
				fmt.Printf("[ipdock] Heartbeat check-in (%s, status=%s)\n", ip, result.Status)
				heartbeatDeadline = nextHeartbeat()
			}
		} else {
			fmt.Printf("[ipdock] IP unchanged (%s) — skipping update\n", ip)
		}

		time.Sleep(time.Duration(intervalSec) * time.Second)
	}
}
