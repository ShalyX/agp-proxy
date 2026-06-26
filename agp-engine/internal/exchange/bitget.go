package exchange

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"agp-engine/internal/config"
	"agp-engine/internal/validator"
)

// GenerateSignature generates the Bitget API signature
func GenerateSignature(timestamp, method, requestPath, body, secretKey string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// PlaceOrder sends the validated order to Bitget
func PlaceOrder(req validator.OrderRequest, cfg *config.RiskConfig) ([]byte, error) {
	if cfg.BitgetAPIKey == "" {
		return nil, fmt.Errorf("bitget_api_key is not configured in risk_config.json")
	}

	requestPath := "/api/v2/mix/order/place-order"
	url := cfg.BitgetAPIURL + requestPath

	// Map AGP fields to Bitget V2 Futures API fields
	payload := map[string]interface{}{
		"symbol":      req.Symbol,
		"productType": "USDT-FUTURES",
		"marginMode":  "isolated", // or crossed, depending on strategy
		"marginCoin":  "USDT",
		"size":        strconv.FormatFloat(req.SizeUSDT, 'f', 2, 64),
		"side":        req.Side,
		"orderType":   "market", // simple market order for now
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	bodyStr := string(bodyBytes)

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := GenerateSignature(timestamp, "POST", requestPath, bodyStr, cfg.BitgetAPISecret)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("ACCESS-KEY", cfg.BitgetAPIKey)
	httpReq.Header.Set("ACCESS-SIGN", signature)
	httpReq.Header.Set("ACCESS-TIMESTAMP", timestamp)
	httpReq.Header.Set("ACCESS-PASSPHRASE", cfg.BitgetAPIPassphrase)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("bitget request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("bitget API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
