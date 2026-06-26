package config

import (
	"encoding/json"
	"os"
	"sync"
)

type RiskConfig struct {
	MaxLeverage         int      `json:"max_leverage"`
	MaxPositionSizeUSDT float64  `json:"max_position_size_usdt"`
	AllowedPairs        []string `json:"allowed_pairs"`
	DiscordWebhookURL   string   `json:"discord_webhook_url"`
	BitgetAPIKey        string   `json:"bitget_api_key"`
	BitgetAPISecret     string   `json:"bitget_api_secret"`
	BitgetAPIPassphrase string   `json:"bitget_api_passphrase"`
	BitgetAPIURL        string   `json:"bitget_api_url"`
}

type ConfigManager struct {
	mu           sync.RWMutex
	baseConfig   *RiskConfig
	activeConfig *RiskConfig
	isVolatile   bool
}

func LoadConfigManager(path string) (*ConfigManager, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg RiskConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	
	active := cfg // copy by value

	return &ConfigManager{
		baseConfig:   &cfg,
		activeConfig: &active,
		isVolatile:   false,
	}, nil
}

func (cm *ConfigManager) GetConfig() *RiskConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	// Return a copy so the caller can't mutate the underlying struct
	cp := *cm.activeConfig
	return &cp
}

func (cm *ConfigManager) SetVolatileMode(volatile bool) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// If no change in state, return false
	if cm.isVolatile == volatile {
		return false
	}

	cm.isVolatile = volatile
	active := *cm.baseConfig // start with fresh base copy

	if volatile {
		active.MaxLeverage = active.MaxLeverage / 2
		if active.MaxLeverage < 1 {
			active.MaxLeverage = 1
		}
		active.MaxPositionSizeUSDT = active.MaxPositionSizeUSDT / 2
	}

	cm.activeConfig = &active
	return true
}

func (cm *ConfigManager) IsVolatile() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.isVolatile
}
