package config

import (
	"os"
	"testing"
)

func TestConfigManagerScaling(t *testing.T) {
	// Create a temporary mock config file
	mockConfig := `{"max_leverage": 10, "max_position_size_usdt": 5000}`
	tmpFile := "mock_risk_config.json"
	os.WriteFile(tmpFile, []byte(mockConfig), 0644)
	defer os.Remove(tmpFile)

	cm, err := LoadConfigManager(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load config manager: %v", err)
	}

	// 1. Check Base Config
	cfg := cm.GetConfig()
	if cfg.MaxLeverage != 10 {
		t.Errorf("Expected base MaxLeverage to be 10, got %d", cfg.MaxLeverage)
	}

	// 2. Trigger Volatility (Penalty should apply)
	changed := cm.SetVolatileMode(true)
	if !changed {
		t.Errorf("Expected SetVolatileMode(true) to return true")
	}

	cfg = cm.GetConfig()
	if cfg.MaxLeverage != 5 {
		t.Errorf("Expected scaled MaxLeverage to be 5, got %d", cfg.MaxLeverage)
	}
	if cfg.MaxPositionSizeUSDT != 2500 {
		t.Errorf("Expected scaled MaxPositionSizeUSDT to be 2500, got %f", cfg.MaxPositionSizeUSDT)
	}

	// 3. Reset Volatility (Base config should be restored)
	changed = cm.SetVolatileMode(false)
	if !changed {
		t.Errorf("Expected SetVolatileMode(false) to return true")
	}

	cfg = cm.GetConfig()
	if cfg.MaxLeverage != 10 {
		t.Errorf("Expected restored MaxLeverage to be 10, got %d", cfg.MaxLeverage)
	}
}
