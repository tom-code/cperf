package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var latenc prometheus.Histogram
type Client struct {
  bucket *tokenBucket
  hcli *http.Client
  url string
  mutex sync.Mutex
  size int
}

func clientLoop(cli *Client) {
  for {
    if cli.bucket.get() {
      go func () {
        data := make([]byte, cli.size)
        start := time.Now()
        _, err := cli.hcli.Post(cli.url, "text/plain", bytes.NewReader(data))
        //_, err := cli.hcli.Get(cli.url)
        if err != nil {
          log.Println(err)
        }
        lat := time.Now().Sub(start)
        latenc.Observe(float64(lat.Milliseconds()))
        //log.Printf("%d\n", lat.Microseconds())
      }()
    } else {
      time.Sleep(1*time.Millisecond)
    }
  }
}

func client() {
  latenc = promauto.NewHistogram(prometheus.HistogramOpts{
    Name:    "bechmark_latency",
    Help:    "http latency",
    Buckets: prometheus.ExponentialBuckets(1, 2, 20),
  })

  http.Handle("/metrics", promhttp.Handler())
  go http.ListenAndServe(":9998", nil)


  hcli := createClient()
  cli := Client{
    bucket: newBucket(0),
    hcli : hcli,
  }
  go clientLoop(&cli)
  for {
    time.Sleep(1*time.Second)
    resp, err := hcli.Get(urlCtrl+"/client")
    if err != nil {
      log.Println(err)
      continue
    }
    if resp.StatusCode != 200 {
      log.Println("ctrl status "+resp.Status)
      continue;
    }
    defer resp.Body.Close()
    var msg ProtoClientResp
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      log.Println(err)
      continue
    }
    log.Println(string(body))
    err = json.Unmarshal(body, &msg)
    if err != nil {
      log.Println(err)
      continue
    }
    cli.mutex.Lock()
    cli.url = msg.Url
    cli.bucket.speed = float64(msg.Speed)
    cli.size = msg.Size
    cli.mutex.Unlock()
    log.Printf("config changed to %f %s\n", cli.bucket.speed, cli.url)
  }
}

