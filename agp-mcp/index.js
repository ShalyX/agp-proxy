import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";

const server = new Server(
  {
    name: "agp-mcp-server",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// In a platform deployment, this points to Bitget's central AGP server.
const GO_PROXY_URL = process.env.AGP_API_URL || "http://localhost:8080/api/v1/order";

server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "place_order",
        description: "Places a spot or futures order on the exchange through the AGP validation engine.",
        inputSchema: {
          type: "object",
          properties: {
            symbol: { type: "string", description: "Trading pair (e.g. BTCUSDT)" },
            side: { type: "string", description: "'buy' or 'sell'" },
            size_usdt: { type: "number", description: "Position size in USDT" },
            leverage: { type: "number", description: "Leverage for futures orders (use 1 for spot)" },
          },
          required: ["symbol", "side", "size_usdt", "leverage"],
        },
      },
    ],
  };
});

server.setRequestHandler(CallToolRequestSchema, async (request) => {
  if (request.params.name === "place_order") {
    const { symbol, side, size_usdt, leverage } = request.params.arguments;

    try {
      const response = await fetch(GO_PROXY_URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ symbol, side, size_usdt, leverage }),
      });

      if (!response.ok) {
        // If the Go proxy blocked the order, it returns a 400 with the error message
        const errorText = await response.text();
        return {
          content: [
            {
              type: "text",
              text: `GUARDRAIL PROXY BLOCKED EXECUTION:\n${errorText.trim()}`,
            },
          ],
          isError: true,
        };
      }

      const data = await response.json();
      return {
        content: [
          {
            type: "text",
            text: `Order executed successfully: ${JSON.stringify(data)}`,
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: "text",
            text: `Failed to connect to AGP Engine: ${error.message}`,
          },
        ],
        isError: true,
      };
    }
  }

  throw new Error(`Unknown tool: ${request.params.name}`);
});

async function run() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("AGP MCP Server running on stdio");
}

run().catch(console.error);
