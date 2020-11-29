package main

import (
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var serverCounter uint64
var serverCounterTotal uint64

func serverHandler(w http.ResponseWriter, r *http.Request) {
  buf := make([]byte, 100000)
  io.ReadFull(r.Body, buf)
  r.Body.Close()
  w.WriteHeader(200)
  atomic.AddUint64(&serverCounter, 1)
  atomic.AddUint64(&serverCounterTotal, 1)
}

func sink() {

  promauto.NewCounterFunc(prometheus.CounterOpts{
    Name: "benchmark_sink_total",
  },
  func() float64 {
    return float64(serverCounterTotal)
  },
  )
  http.Handle("/metrics", promhttp.Handler())
  go http.ListenAndServe(":9999", nil)


  h2s := &http2.Server {
  }
  handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    serverHandler(w, r)
  })
  h1s := &http.Server {
    Addr:    ":80",
    Handler: h2c.NewHandler(handler, h2s),
  }
  go func() {
    for {
      time.Sleep(1*time.Second)
      log.Println(serverCounter)
      serverCounter = 0
    }
  }()
  h1s.ListenAndServe()
}