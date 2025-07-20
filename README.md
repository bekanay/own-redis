# own-redis

A simple in-memory UDP-based key–value database inspired by Redis, written in Go.

## 🚀 Features

* **SET**: store a key–value pair, with optional expiry (PX in milliseconds)
* **GET**: retrieve a value by key
* **PING**: health check (replies with `PONG`)
* Commands and arguments are case-insensitive
* Concurrent-safe with `sync.RWMutex`
* Automatic cleanup of expired keys

## 🛠️ Prerequisites

* Go 1.18+
* [`gofumpt`](https://github.com/mvdan/gofumpt) for formatting
* `netcat` (nc) for manual testing
* `bc` for stress-test timing

## 📦 Installation & Build

1. Clone the repo and enter directory:

   ```bash
   git clone https://github.com/yourusername/own-redis.git
   cd own-redis
   ```
2. Install `gofumpt` if you haven’t:

   ```bash
   go install mvdan.cc/gofumpt@latest
   ```
3. Format and build:

   ```bash
   gofumpt -w .
   go build -o own-redis .
   ```

## ⚙️ Usage

```bash
own-redis [--port <N>]
own-redis --help
```

* **--port N**  : UDP port to listen on (default: 8080)
* **--help**    : Show usage information

### Example session with netcat

```bash
# Start server on 8080
./own-redis --port 8080

# In another terminal:
nc -u localhost 8080

> PING
PONG

> SET foo bar
OK

> GET foo
bar

> SET temp 123 PX 1000
OK
> GET temp
123
# wait ≥1s
> GET temp
(nil)
```

## 🔍 Testing

### Automated Go tests

A basic integration test lives in `store_test.go`. It expects the server running on port 8080.

```bash
go test -timeout 5s
```

### Stress testing with Bash

A script `stress_test.sh` performs concurrent SET/GET operations. Make it executable and run:

```bash
chmod +x stress_test.sh
./stress_test.sh
```

**Tip:** for higher throughput replace `nc` spawns with a persistent UDP socket in Bash, or use the provided Go benchmark in `bench.go`:

```bash
go run bench.go
```

## 🧪 Race Detection

Ensure no data races:

```bash
go run -race .
```