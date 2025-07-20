package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Store struct {
	mu      sync.RWMutex
	data    map[string]string
	expires map[string]time.Time
}

func NewStore() *Store {
	s := &Store{
		data:    make(map[string]string),
		expires: make(map[string]time.Time),
	}
	go s.cleanupLoop()
	return s
}

func (s *Store) cleanupLoop() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for k, exp := range s.expires {
			if now.After(exp) {
				delete(s.data, k)
				delete(s.expires, k)
			}
		}
		s.mu.Unlock()
	}
}

func (s *Store) Set(key, val string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
	if ttl > 0 {
		s.expires[key] = time.Now().Add(ttl)
	} else {
		delete(s.expires, key)
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	exp, hasExp := s.expires[key]
	if hasExp && time.Now().After(exp) {
		return "", false
	}
	val, ok := s.data[key]
	return val, ok
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "")
	flag.Usage = func() {
		fmt.Println(`Own Redis

Usage:
  own-redis [--port <N>]
  own-redis --help

Options:
  --help       Show this screen.
  --port N     Port number.
`)
	}
	flag.Parse()

	store := NewStore()
	addr := net.UDPAddr{Port: port, IP: net.IPv4zero}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Listening on port: " + strconv.Itoa(addr.Port))
	buf := make([]byte, 4096)
	for {
		n, client, _ := conn.ReadFromUDP(buf)
		go func(msg string) {
			reply := handleLine(store, msg)
			conn.WriteToUDP([]byte(reply+"\n"), client)
		}(string(buf[:n]))
	}
}

func handleLine(s *Store, line string) string {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return ""
	}
	cmd := strings.ToUpper(parts[0])
	args := parts[1:]

	switch cmd {
	case "PING":
		return "PONG"
	case "SET":
		if len(args) < 2 {
			return "(error) ERR wrong number of arguments for 'SET' command"
		}
		key := args[0]
		// join the rest, detect PX
		// ... parse ttl if present ...
		value := strings.Join(args[1:], " ")
		s.Set(key, value, 0)
		return "OK"
	case "GET":
		if len(args) != 1 {
			return "(error) ERR wrong number of arguments for 'GET' command"
		}
		if v, ok := s.Get(args[0]); ok {
			return v
		}
		return "(nil)"
	default:
		return "(error) ERR unknown command"
	}
}
