package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"net"
	"time"
  _ "net/http/pprof"
)

func fmtStatStr(stat string, tags map[string]string) string {
	parts := []string{}
	for k, v := range tags {
		if v != "" {
			parts = append(parts, fmt.Sprintf("%s:%s", k, v))
		}
	}

	return fmt.Sprintf("%s|%s", stat, strings.Join(parts, ","))
}

func main() {

	stats, err := NewSimpleClient("localhost:6060")
	if err != nil {
		log.Fatal("could not start stas client: ", err)
	}

	// add handlers to default mux
	http.HandleFunc("/ping", pingHandler(stats))

	s := &http.Server{
		Addr:    ":8080",
	}

	log.Fatal(s.ListenAndServe())
}

// MockClient implements a in-memory statsd client which never flushes to a server.
type SimpleClient struct {
	buffer []string

	c net.PacketConn

	ra *net.UDPAddr
}

func NewSimpleClient(addr string) (*SimpleClient, error){

	c, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return nil, err
	}

	ra, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		c.Close()
		return nil, err
	}

	return &SimpleClient{
		c:  c,
		ra: ra,
	}, nil
}

func (sc *SimpleClient) send(s string) error {

	sc.buffer = append(sc.buffer, s)
	if len(sc.buffer) > 100 {

		b := strings.Join(sc.buffer, ",")
		_, err := sc.c.(*net.UDPConn).WriteToUDP([]byte(b), sc.ra)
		if err != nil {
			return err
		}

		sc.buffer = nil
	}

	return  nil
}

func (sc *SimpleClient) Timing(s string, d time.Duration, sampleRate float64,
	tags map[string]string) error {
	return sc.send(fmtStatStr(
		fmt.Sprintf("%s:%d|ms",s, d/ time.Millisecond), tags),
	)
}

func pingHandler(s *SimpleClient) http.HandlerFunc{

	return func(w http.ResponseWriter, r *http.Request) {
		st := time.Now()
		defer func() {_ = s.Timing("http.ping", time.Since(st), 1.0, nil)}()

		w.WriteHeader(200)

	}
}
