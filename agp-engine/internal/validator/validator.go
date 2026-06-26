package validator

import (
	"fmt"
	"agp-engine/internal/config"
)

type OrderRequest struct {
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"` // "buy" or "sell"
	SizeUSDT float64 `json:"size_usdt"`
	Leverage int     `json:"leverage"`
}

func ValidateOrder(req OrderRequest, cfg *config.RiskConfig) error {
	// 1. Check Allowed Pairs
	pairAllowed := false
	for _, p := range cfg.AllowedPairs {
		if req.Symbol == p {
			pairAllowed = true
			break
		}
	}
	if !pairAllowed {
		return fmt.Errorf("hallucination detected: unauthorized pair %s", req.Symbol)
	}

	// 2. Check Leverage
	if req.Leverage > cfg.MaxLeverage {
		return fmt.Errorf("hallucination detected: leverage %d exceeds max %d", req.Leverage, cfg.MaxLeverage)
	}

	// 3. Check Position Size
	if req.SizeUSDT > cfg.MaxPositionSizeUSDT {
		return fmt.Errorf("hallucination detected: position size %f exceeds max %f", req.SizeUSDT, cfg.MaxPositionSizeUSDT)
	}

	return nil
}
