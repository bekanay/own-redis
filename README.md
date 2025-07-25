# own-redis

A simple in-memory UDP-based key–value database inspired by Redis, written in Go.

## 🚀 Features

* **SET**: store a key–value pair, with optional expiry (PX in milliseconds)
* **GET**: retrieve a value by key
* **PING**: health check (replies with `PONG`)
* Commands and arguments are case-insensitive
* Concurrent-safe with `sync.RWMutex`
* Automatic cleanup of expired keys

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
