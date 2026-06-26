# Agent Hallucination Guardrail Proxy (AGP) 🛡️
**Track:** Trading Infra

## The Problem
As trading bots transition from rigid scripts to autonomous LLMs (via Claude Code, OpenClaw, etc.), the risk of model hallucination becomes a direct financial liability. An LLM misinterpreting a sentiment signal or glitching in its reasoning loop can easily lead to max-leverage liquidations or account-draining wash trading. 

**Allowing agents to trade safely is just as important as allowing them to trade smartly.**

## The Solution
The **Agent Hallucination Guardrail Proxy (AGP)** acts as a deterministic "kill switch" layer between the AI agent and the exchange. 

AGP intercepts the payload and runs strict, non-AI mathematical checks against user-defined risk parameters. If an AI attempts a trade outside these absolute bounds, the proxy blocks the execution, triggers a Discord alert, and returns an error message to the agent's context window, forcing it to recalibrate.

### 🏗️ Platform Deployment Architecture
While this repository runs the proxy locally for demonstration, the intended architecture is a **Platform Solution**. 
Bitget (or a centralized protocol) hosts the blazing-fast Go Validation Engine on their cloud infrastructure right next to their matching engines. The end user simply points their Node.js MCP Server to Bitget's AGP URL (e.g., `https://agp-guard.bitget.com`). 
This provides a **zero-maintenance, zero-latency** security layer for users deploying autonomous agents.

## Key Features

🚀 **Sub-millisecond Validation Engine (Go)**  
Built in Go to ensure near-zero latency so execution prices aren't ruined during fast market movements. Validates allowed trading pairs, hard-capped leverage limits, and maximum position sizing.

📉 **Dynamic Risk Limits (Market Volatility Scanner)**  
A background goroutine continuously polls the public Bitget Ticker API. If it detects extreme market volatility (e.g., a 24-hour price change of ±10%), the proxy automatically locks state and slashes the AI's allowed leverage and position sizing by 50% to protect the account.

🔌 **Native Bitget V2 Integration**  
Fully integrated with the Bitget Agent Hub. The proxy natively signs approved payloads using HMAC-SHA256 cryptography and routes them directly to the Bitget Futures API.

🤖 **Node.js MCP Wrapper**  
A dedicated Model Context Protocol (MCP) server that wraps the execution tool. The AI agent uses standard tools, completely unaware that its executions are being aggressively monitored and filtered by the remote Go engine.

🚨 **Real-time Discord Alerting**  
Immediate Discord webhook alerts whenever an agent's trade is blocked, an exchange error occurs, or the Dynamic Risk Limits are triggered due to market turbulence.

---

## Quick Start (Local Prototype)

### 1. Prerequisites
- Go 1.22+
- Node.js v24+

### 2. Configuration
Navigate to the engine directory and configure your risk limits and API keys:
```bash
cd agp-engine
# Open risk_config.json and add your Bitget API keys and Discord Webhook URL
```

### 3. Run the Go Validation Engine
```bash
cd agp-engine
go run cmd/server/main.go
```
*The engine will start on port 8080 and the volatility scanner will begin polling the market.*

### 4. Start the MCP Server
In a new terminal window:
```bash
cd agp-mcp
npm install
npm start
```
*Your AI agent can now connect to this MCP server to place strictly guarded trades! Note: In production, you would launch this with `AGP_API_URL=https://agp-guard.bitget.com`.*
