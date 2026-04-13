# Mini Exchange

Real-time mini stock exchange system with Go & Fiber — featuring order matching engine, WebSocket live updates, and market simulation.

## Features

| Feature | Description |
|---------|-------------|
| Order Matching Engine | Price-time priority matching |
| WebSocket Updates | Live order book, trade & ticker |
| Market Simulation | Auto-generate orders for testing |
| Clean Architecture | Handler → Service → Repository |
| Structured Logging | Zerolog with configurable format |

---

## Quick Start

```bash
# Clone & setup
git clone <repository-url> && cd mini-exchange
cp .env.example .env

# Run
go mod download && go run main.go
```

Server: `http://localhost:3000`

---

## Architecture

```
HTTP Request ──▶ Handler ──▶ Service ──▶ Repository ──▶ In-Memory Storage
                    │
                    ▼
              MatchingService ──▶ WebSocket Hub ──▶ Clients
```

| Layer | Responsibility |
|-------|----------------|
| **Handlers** | Parse request, validate, return response |
| **Services** | Business logic, order matching |
| **Repositories** | Data storage with mutex locks |
| **WebSocket Hub** | Pub/sub broadcast to clients |
| **Tasks** | Background: matching engine, simulation |

### Matching Algorithm

**Price-Time Priority:**
- BUY matches SELL where `sell.price ≤ buy.price`
- SELL matches BUY where `buy.price ≥ sell.price`
- Partial fills supported, remaining goes to order book

---

## API Reference

**Base URL:** `http://localhost:3000`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/greeting` | Greeting message |
| POST | `/api/order` | Create order |
| GET | `/api/order` | List all orders |
| GET | `/api/trade` | Trade history |
| GET | `/api/market/snapshot` | Market data snapshot |

### Create Order

```bash
curl -X POST http://localhost:3000/api/order \
  -H "Content-Type: application/json" \
  -d '{"stock_code": "IDR", "side": "BUY", "price": 65000, "quantity": 10}'
```

**Request Body:**
```json
{
  "stock_code": "IDR",      // Required: asset code
  "side": "BUY",            // Required: BUY or SELL  
  "price": 65000,           // Required: price > 0
  "quantity": 10            // Required: quantity > 0
}
```

### Market Snapshot

```bash
curl "http://localhost:3000/api/market/snapshot?stock_code=IDR&depth=10"
```

**Query Params:** `stock_code` (required), `depth` (default: 20), `trade_limit` (default: 50)

---

## WebSocket

**Connect:** `ws://localhost:3000/ws`

### Subscribe/Unsubscribe

```json
{"type": "subscribe", "channel": "market.ticker:IDR"}
{"type": "unsubscribe", "channel": "market.ticker:IDR"}
```

### Available Channels

| Channel | Description | Data |
|---------|-------------|------|
| `market.ticker:{stock}` | Price updates | `last_price`, `change`, `volume` |
| `market.orderbook:{stock}` | Order book changes | `buys[]`, `sells[]` |
| `market.trade:{stock}` | Trade executions | `price`, `quantity`, `buy_order_id`, `sell_order_id` |
| `order.update` | Order status changes | `order_id`, `status`, `filled_quantity` |

> Channel names are case-insensitive

### JavaScript Example

```javascript
const ws = new WebSocket('ws://localhost:3000/ws');

ws.onopen = () => {
  ws.send(JSON.stringify({type: 'subscribe', channel: 'market.ticker:IDR'}));
  ws.send(JSON.stringify({type: 'subscribe', channel: 'market.trade:IDR'}));
};

ws.onmessage = (e) => {
  const {channel, data} = JSON.parse(e.data);
  console.log(channel, data);
};
```

---

## Configuration

```env
APP_PORT=:3000
APP_ENV=development
LOG_LEVEL=info
LOG_FORMAT=console
RATE_LIMIT_MAX=100
RATE_LIMIT_EXPIRATION=60
SIMULATION_ENABLED=true
SIMULATION_INTERVAL=2s
```

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | `:3000` | Server port |
| `LOG_LEVEL` | `info` | debug/info/warn/error |
| `SIMULATION_ENABLED` | `true` | Enable market simulation |
| `SIMULATION_INTERVAL` | `2s` | Order generation interval |

---

## Technical Notes

### Race Condition Handling

| Area | Strategy |
|------|----------|
| Order Book | `sync.RWMutex` for read/write protection |
| Matching Engine | Single goroutine worker via channel queue |
| Trade Store | Append-only with mutex lock |
| WebSocket | Buffered channel per client + drop strategy |

### Non-Blocking WebSocket Broadcast

- Per-client buffered channel (`256` buffer)
- Dedicated write goroutine per client
- Drop messages if client buffer full (backpressure)
- Fan-out pattern via central Hub

### Known Bottlenecks

- Single matching engine goroutine (throughput limit)
- In-memory storage (no persistence, memory growth)
- WebSocket fan-out cost at high client count

---

## License

MIT