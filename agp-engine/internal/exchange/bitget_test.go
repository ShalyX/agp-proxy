package exchange

import (
	"testing"
)

func TestGenerateSignature(t *testing.T) {
	timestamp := "1645000000000"
	method := "POST"
	requestPath := "/api/v2/mix/order/place-order"
	body := `{"symbol":"BTCUSDT","size":"1000.00"}`
	secretKey := "mysecretkey"

	sig := GenerateSignature(timestamp, method, requestPath, body, secretKey)
	
	if len(sig) == 0 {
		t.Errorf("GenerateSignature() returned empty string")
	}
}
