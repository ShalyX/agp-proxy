Write-Host "--- AGP Proxy Integration Tests ---" -ForegroundColor Cyan

$baseUrl = "http://localhost:8080/api/v1/order"
$headers = @{"Content-Type"="application/json"}

Write-Host "1. Testing Valid Order..."
$validPayload = '{"symbol": "BTCUSDT", "side": "buy", "size_usdt": 1000, "leverage": 5}'
try {
    $response = Invoke-RestMethod -Uri $baseUrl -Method Post -Headers $headers -Body $validPayload
    Write-Host "SUCCESS: Valid order passed (or tried to hit Bitget)!" -ForegroundColor Green
} catch {
    Write-Host "FAIL: Valid order was blocked. Status: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Red
}

Write-Host "`n2. Testing Hard Leverage Guardrail (20x)..."
$highLevPayload = '{"symbol": "BTCUSDT", "side": "buy", "size_usdt": 1000, "leverage": 20}'
try {
    $response = Invoke-RestMethod -Uri $baseUrl -Method Post -Headers $headers -Body $highLevPayload
    Write-Host "FAIL: High leverage order was allowed through!" -ForegroundColor Red
} catch {
    Write-Host "SUCCESS: High leverage order was blocked as expected." -ForegroundColor Green
}

Write-Host "`n3. Testing Unauthorized Pair (DOGEUSDT)..."
$badPairPayload = '{"symbol": "DOGEUSDT", "side": "buy", "size_usdt": 1000, "leverage": 5}'
try {
    $response = Invoke-RestMethod -Uri $baseUrl -Method Post -Headers $headers -Body $badPairPayload
    Write-Host "FAIL: Unauthorized pair order was allowed through!" -ForegroundColor Red
} catch {
    Write-Host "SUCCESS: Unauthorized pair order was blocked as expected." -ForegroundColor Green
}
