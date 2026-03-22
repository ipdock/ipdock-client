package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const version = "0.1.0"

var (
	apiURL   string
	token    string
	interval time.Duration
	lastIP   string
)

func init() {
	apiURL = getEnv("IPDOCK_API_URL", "https://api.ipdock.io")
	token = os.Getenv("IPDOCK_TOKEN")
	if token == "" {
		log.Fatal("[ipdock] IPDOCK_TOKEN is required")
	}
	secs, err := strconv.Atoi(getEnv("IPDOCK_INTERVAL", "60"))
	if err != nil || secs < 10 {
		secs = 60
	}
	interval = time.Duration(secs) * time.Second
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getPublicIP tries multiple IP detection services in order
func getPublicIP() (string, error) {
	services := []string{
		"https://api4.ipify.org",
		"https://ipv4.icanhazip.com",
		"https://checkip.amazonaws.com",
	}

	client := &http.Client{Timeout: 10 * time.Second}
	for _, svc := range services {
		resp, err := client.Get(svc)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		ip := strings.TrimSpace(string(body))
		if ip != "" {
			return ip, nil
		}
	}
	return "", fmt.Errorf("could not detect public IP from any service")
}

// sendUpdate posts the IP to the ipdock API
func sendUpdate(ip string) error {
	payload := map[string]string{"ip": ip}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL+"/api/v1/update", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ipdock-client/"+version)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("API returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func run() {
	ip, err := getPublicIP()
	if err != nil {
		log.Printf("[ipdock] ERROR detecting IP: %v", err)
		return
	}

	if ip == lastIP {
		log.Printf("[ipdock] IP unchanged (%s) — skipping update", ip)
		return
	}

	log.Printf("[ipdock] IP changed: %s → %s — sending update", lastIP, ip)
	if err := sendUpdate(ip); err != nil {
		log.Printf("[ipdock] ERROR sending update: %v", err)
		return
	}

	log.Printf("[ipdock] Update sent successfully: %s", ip)
	lastIP = ip
}

func main() {
	log.Printf("[ipdock] client v%s starting", version)
	log.Printf("[ipdock] API: %s | interval: %s", apiURL, interval)

	// Run immediately on start
	run()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		run()
	}
}
