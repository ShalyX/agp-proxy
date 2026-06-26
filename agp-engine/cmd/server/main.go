package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"agp-engine/internal/config"
	"agp-engine/internal/exchange"
	"agp-engine/internal/market"
	"agp-engine/internal/notifier"
	"agp-engine/internal/validator"
)

var cfgManager *config.ConfigManager

func orderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req validator.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve a safe snapshot of the current active configuration
	cfg := cfgManager.GetConfig()

	// 1. Guardrail Validation
	err := validator.ValidateOrder(req, cfg)
	if err != nil {
		alertMsg := fmt.Sprintf("⚠️ **GUARDRAIL TRIGGERED**\nBlocked Trade: %+v\nReason: %s", req, err.Error())
		notifier.SendAlert(cfg.DiscordWebhookURL, alertMsg)
		
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Execute on Exchange
	log.Printf("✅ Order passed validation (MaxLeverage: %d). Forwarding to Bitget...\n", cfg.MaxLeverage)
	respBody, err := exchange.PlaceOrder(req, cfg)
	if err != nil {
		log.Printf("❌ Exchange execution failed: %v", err)
		notifier.SendAlert(cfg.DiscordWebhookURL, fmt.Sprintf("❌ **EXCHANGE ERROR**\nTrade failed to execute: %v", err))
		
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func main() {
	var err error
	cfgManager, err = config.LoadConfigManager("risk_config.json")
	if err != nil {
		log.Fatalf("Failed to load config manager: %v", err)
	}

	// Start the background Goroutine for fetching market data
	market.StartMonitor(cfgManager)

	http.HandleFunc("/api/v1/order", orderHandler)

	port := "8080"
	log.Printf("Starting AGP Engine with Dynamic Risk Limits on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
