package validator

import (
	"agp-engine/internal/config"
	"testing"
)

func TestValidateOrder(t *testing.T) {
	cfg := &config.RiskConfig{
		MaxLeverage:         10,
		MaxPositionSizeUSDT: 5000,
		AllowedPairs:        []string{"BTCUSDT"},
	}

	tests := []struct {
		name    string
		req     OrderRequest
		wantErr bool
	}{
		{
			name: "valid order",
			req: OrderRequest{Symbol: "BTCUSDT", Side: "buy", SizeUSDT: 1000, Leverage: 5},
			wantErr: false,
		},
		{
			name: "unauthorized pair",
			req: OrderRequest{Symbol: "DOGEUSDT", Side: "buy", SizeUSDT: 1000, Leverage: 5},
			wantErr: true,
		},
		{
			name: "leverage too high",
			req: OrderRequest{Symbol: "BTCUSDT", Side: "buy", SizeUSDT: 1000, Leverage: 20},
			wantErr: true,
		},
		{
			name: "position size too large",
			req: OrderRequest{Symbol: "BTCUSDT", Side: "buy", SizeUSDT: 10000, Leverage: 5},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrder(tt.req, cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
