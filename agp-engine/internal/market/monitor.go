package market

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"agp-engine/internal/config"
	"agp-engine/internal/notifier"
)

const (
	VolatilityThreshold = 0.10 // 10%
	PollInterval        = 30 * time.Second
	TickerURL           = "https://api.bitget.com/api/v2/mix/market/ticker?symbol=BTCUSDT&productType=USDT-FUTURES"
)

type TickerResponse struct {
	Code string `json:"code"`
	Data []struct {
		PriceChangePercent string `json:"priceChangePercent"`
	} `json:"data"`
}

func fetchVolatility() (float64, error) {
	resp, err := http.Get(TickerURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data TickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	if len(data.Data) == 0 {
		return 0, fmt.Errorf("no ticker data found")
	}

	percentStr := data.Data[0].PriceChangePercent
	percent, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		return 0, err
	}

	return percent, nil
}

func StartMonitor(cm *config.ConfigManager) {
	log.Println("📈 Starting Volatility Monitor for BTCUSDT...")
	ticker := time.NewTicker(PollInterval)
	go func() {
		for {
			<-ticker.C
			percent, err := fetchVolatility()
			if err != nil {
				log.Printf("⚠️ Monitor failed to fetch ticker: %v", err)
				continue
			}

			// We treat the absolute price change as volatility
			isVolatile := percent >= VolatilityThreshold || percent <= -VolatilityThreshold

			// If state changed, notify Discord
			changed := cm.SetVolatileMode(isVolatile)
			if changed {
				cfg := cm.GetConfig()
				if isVolatile {
					msg := fmt.Sprintf("🚨 **HIGH VOLATILITY DETECTED (%.2f%%)**\nRisk limits have been halved!\nNew Max Leverage: %dx\nNew Max Position: $%.2f", percent*100, cfg.MaxLeverage, cfg.MaxPositionSizeUSDT)
					notifier.SendAlert(cfg.DiscordWebhookURL, msg)
					log.Println(msg)
				} else {
					msg := fmt.Sprintf("✅ **MARKET STABILIZED (%.2f%%)**\nRisk limits restored to normal.\nMax Leverage: %dx\nMax Position: $%.2f", percent*100, cfg.MaxLeverage, cfg.MaxPositionSizeUSDT)
					notifier.SendAlert(cfg.DiscordWebhookURL, msg)
					log.Println(msg)
				}
			}
		}
	}()
}
